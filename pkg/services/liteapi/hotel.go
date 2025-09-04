package liteapi

import (
	"context"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"time"

	"github.com/samber/do"
)

type HotelSvcs interface {
	SearchHotels(ctx context.Context, query map[string]string) (any, error)
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

	//step 1: check if we have the hotels in the database

	resp, err := h.liteApiSdk.GetHotels(query, "en", 3, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Data)
	}

	return resp, nil

}
