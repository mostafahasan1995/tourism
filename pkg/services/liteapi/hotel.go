package liteapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"net/http"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelSvcs interface {
	GetHotels(ctx context.Context, query map[string]string) (*models.HotelListRes, error)
	GetHotelDetails(ctx context.Context, id, language, advancedAccessibilityOnly string) (*models.HotelDetailsData, error)
	TestStreaming(ctx context.Context, w http.ResponseWriter, query map[string]string) error
	TestStreaming2(ctx context.Context, w io.Writer, query map[string]string) error
}

type hotelssvcs struct {
	repo             repo.HotelRepo
	hotelDetailsRepo repo.HotelDetailsRepo
	liteApiInitFunc  liteApiSdk.LiteApiInitFunc
}

func NewHotelSvcs(i *do.Injector) (HotelSvcs, error) {
	return &hotelssvcs{
		repo:             do.MustInvoke[repo.HotelRepo](i),
		hotelDetailsRepo: do.MustInvoke[repo.HotelDetailsRepo](i),
		liteApiInitFunc:  do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
	}, nil
}

func (h *hotelssvcs) GetHotels(ctx context.Context, query map[string]string) (*models.HotelListRes, error) {
	lang := "en"
	if query["language"] != "" {
		lang = query["language"]
	}

	liteApiSdk, err := h.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetHotels(query, lang, 3, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.HotelListRes
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go func(ctx context.Context, hotels []models.Hotel) {
		writeOps := []mongo.WriteModel{}
		for _, hotel := range hotels {
			updateOp := mongo.NewUpdateOneModel().SetFilter(bson.M{"id": hotel.Id}).SetUpdate(bson.M{"$set": hotel}).SetUpsert(true)
			writeOps = append(writeOps, updateOp)
		}

		if _, err := h.repo.BulkWrite(ctx, writeOps); err != nil {
			fmt.Printf("error bulk writing hotels: %v", err)
		}

	}(context.WithoutCancel(ctx), result.Data)

	return &result, nil

}

func (h *hotelssvcs) GetHotelDetails(ctx context.Context, id, language, advancedAccessibilityOnly string) (*models.HotelDetailsData, error) {
	filter := bson.M{
		"id":        id,
		"expiresAt": bson.M{"$gt": time.Now()},
	}

	result, err := h.hotelDetailsRepo.GetByFilter(ctx, filter)
	if err != nil {
		liteApiSdk, err := h.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		lang := "en"
		if language != "" {
			lang = language
		}

		resp, err := liteApiSdk.GetHotelDetails(id, lang, advancedAccessibilityOnly)
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.HotelDetailsData
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		result.Data["expiresAt"] = time.Now().Add(ExpireTime)

		go func(ctx context.Context, data models.HotelDetails) {
			writeOps := []mongo.WriteModel{}
			updateOp := mongo.NewUpdateOneModel().SetFilter(bson.M{"id": data["id"]}).SetUpdate(bson.M{"$set": data}).SetUpsert(true)
			writeOps = append(writeOps, updateOp)

			h.hotelDetailsRepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), result.Data)

		return &result, nil
	}

	return &models.HotelDetailsData{
		Data: *result,
	}, nil

}

func (h *hotelssvcs) TestStreaming(ctx context.Context, w http.ResponseWriter, query map[string]string) error {

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		pw.Write([]byte(`[`))

		for i := 0; i < 10000; i++ {
			pw.Write([]byte(fmt.Sprintf(`{"id":%d,"name":"Hotel %d"}`, i, i)))
			pw.Write([]byte(`,`))
			if i > 100 {
				time.Sleep(100 * time.Millisecond)
			}
		}

		pw.Write([]byte(`]`))
	}()

	io.Copy(w, pr)

	return nil

}

func (h *hotelssvcs) TestStreaming2(ctx context.Context, w io.Writer, query map[string]string) error {

	fmt.Fprintf(w, `[`)

	for i := 0; i < 10000; i++ {
		w.Write([]byte(fmt.Sprintf(`{"id":%d,"name":"Hotel %d"}`, i, i)))
		w.Write([]byte(`,`))
		if i > 100 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Fprintf(w, `]`)

	return nil

}
