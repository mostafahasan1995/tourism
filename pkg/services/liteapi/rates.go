package liteapi

import (
	"context"
	"encoding/json"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"

	"github.com/samber/do"
)

type RatesSvcs interface {
	GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error)
	GetMinRates(ctx context.Context, data any) (any, error)
	PreBook(ctx context.Context, data map[string]any) (any, error)
	Book(ctx context.Context, data map[string]any) (any, error)
}

type ratessvcs struct {
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
	preBookRepo     repo.PreBookRepo
	bookingRepo     repo.BookingRepo
}

func NewRatesSvcs(i *do.Injector) (RatesSvcs, error) {
	return &ratessvcs{
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		preBookRepo:     do.MustInvoke[repo.PreBookRepo](i),
		bookingRepo:     do.MustInvoke[repo.BookingRepo](i),
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

	var result models.PreBook
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go func(ctx context.Context, data models.PreBook) {
		if err := r.preBookRepo.Add(ctx, &data); err != nil {
			fmt.Printf("error adding prebook: %v", err)
		}

	}(context.WithoutCancel(ctx), result)

	return result, nil
}

func (r *ratessvcs) Book(ctx context.Context, data map[string]any) (any, error) {
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

	var result models.Booking
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go func(ctx context.Context, data models.Booking) {
		if err := r.bookingRepo.Add(ctx, &data); err != nil {
			fmt.Printf("error adding booking: %v", err)
		}
	}(context.WithoutCancel(ctx), result)

	return result, nil
}
