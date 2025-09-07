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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelSvcs interface {
	SearchHotels(ctx context.Context, query map[string]string) (any, error)
	TestStreaming(ctx context.Context, w http.ResponseWriter, query map[string]string) error
	TestStreaming2(ctx context.Context, w io.Writer, query map[string]string) error
}

type hotelssvcs struct {
	repo       repo.HotelRepo
	liteApiSdk *liteApiSdk.LiteApiSdk
}

func NewHotelSvcs(i *do.Injector) (HotelSvcs, error) {
	return &hotelssvcs{
		repo:       do.MustInvoke[repo.HotelRepo](i),
		liteApiSdk: do.MustInvoke[*liteApiSdk.LiteApiSdk](i),
	}, nil
}

func (h *hotelssvcs) SearchHotels(ctx context.Context, query map[string]string) (any, error) {

	resp, err := h.liteApiSdk.GetHotels(query, "en", 3, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	type aux struct {
		Data     []models.Hotel `json:"data"`
		HotelIds []string       `json:"hotelIds"`
		Total    int            `json:"total"`
	}

	var result aux
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil

}

func (h *hotelssvcs) SaveHotels(ctx context.Context, hotels []models.Hotel) error {
	for _, hotel := range hotels {
		filter := bson.M{"id": hotel.Id}
		update := bson.M{"$set": hotel}

		upsert := true
		after := options.After
		opts := &options.FindOneAndUpdateOptions{
			Upsert:         &upsert,
			ReturnDocument: &after,
		}

		if _, err := h.repo.Patch(ctx, filter, update, opts); err != nil {
			return err
		}
	}

	return nil

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
