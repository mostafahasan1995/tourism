package liteapi

import (
	"context"
	"encoding/json"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
)

type RatesSvcs interface {
	GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error)
}

type ratessvcs struct {
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
}

func NewRatesSvcs(i *do.Injector) (RatesSvcs, error) {
	return &ratessvcs{
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
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
