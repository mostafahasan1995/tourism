package liteapi

import (
	"context"
	"encoding/json"
	"errors"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"larsa-tourism-microservices/pkg/services/member"
	"time"

	"larsa-tourism-microservices/pkg/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatesSvcs interface {
	GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error)
	GetMinRates(ctx context.Context, data any) (any, error)
	PreBook(ctx context.Context, data map[string]any) (any, error)
	Book(ctx context.Context, data map[string]any) (any, error)
	MyPrebooks(ctx context.Context) ([]models.UserPrebook, error)
	MyBookings(ctx context.Context) ([]models.UserBooking, error)
	CancelBooking(ctx context.Context, bookingId string) (any, error)
}

type ratessvcs struct {
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
	preBookRepo     repo.PreBookRepo
	bookingRepo     repo.BookingRepo
	customersvcs    member.CustomerSvcs
}

func NewRatesSvcs(i *do.Injector) (RatesSvcs, error) {
	return &ratessvcs{
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		preBookRepo:     do.MustInvoke[repo.PreBookRepo](i),
		bookingRepo:     do.MustInvoke[repo.BookingRepo](i),
		customersvcs:    do.MustInvoke[member.CustomerSvcs](i),
	}, nil
}

func (r *ratessvcs) GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error) {
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetFullRates(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.RatesList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ratessvcs) GetMinRates(ctx context.Context, data any) (any, error) {
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetMinRates(data)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ratessvcs) PreBook(ctx context.Context, data map[string]any) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.customersvcs.GetOne(ctx, cfg.User.Id.Hex())
	if err != nil {
		return nil, errors.New("error get customer")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.PreBook(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.PreBookData
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	userPrebook := &models.UserPrebook{
		PreBook:    result.Data,
		GusetLevel: result.GuestLevel,
		UserId:     cfg.User.Id,
		Status:     "pending",
		CreatedAt:  time.Now(),
		Trash:      false,
	}

	go util.WithRetry(func() error {
		return r.preBookRepo.Add(context.WithoutCancel(ctx), userPrebook)
	}, 3)

	//in case the liteApi request success we must send the response
	return result, nil
}

func (r *ratessvcs) Book(ctx context.Context, data map[string]any) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.customersvcs.GetOne(ctx, cfg.User.Id.Hex())
	if err != nil {
		return nil, errors.New("error get customer")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.Book(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.BookingData
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	userBooking := &models.UserBooking{
		Booking:    result.Data,
		GusetLevel: result.GuestLevel,
		UserId:     cfg.User.Id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	go util.WithRetry(func() error {
		return r.bookingRepo.Add(context.WithoutCancel(ctx), userBooking)
	}, 3)

	return result, nil
}

func (r *ratessvcs) MyPrebooks(ctx context.Context) ([]models.UserPrebook, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": cfg.User.Id, "trash": false}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var result []models.UserPrebook
	err = r.preBookRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ratessvcs) MyBookings(ctx context.Context) ([]models.UserBooking, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": cfg.User.Id}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var result []models.UserBooking
	err = r.bookingRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (r *ratessvcs) CancelBooking(ctx context.Context, bookingId string) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.bookingRepo.GetByFilter(ctx, bson.M{"userId": cfg.User.Id, "bookingId": bookingId})
	if err != nil {
		return nil, errors.New("booking not found or not belongs to the user")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.CancelBooking(bookingId)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go util.WithRetry(func() error {
		filter := bson.M{"userId": cfg.User.Id, "bookingId": bookingId}
		update := bson.M{"$set": bson.M{"status": "CANCELLED", "updatedAt": time.Now()}}
		_, err := r.bookingRepo.Patch(context.WithoutCancel(ctx), filter, update)
		return err
	}, 3)

	return result, nil
}
