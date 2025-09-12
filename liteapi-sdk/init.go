package liteApiSdk

import (
	"context"
	"fmt"
	"larsa-tourism-microservices/pkg/util"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
)

type LiteApiInitFunc func(ctx context.Context) (*LiteApiSdk, error)

func LiteApiSdkInit() LiteApiInitFunc {
	return func(ctx context.Context) (*LiteApiSdk, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get request cfg: %w", err)
		}

		serviceToken, err := common.GetServiceToken("tourism", cfg.Hp)
		if err != nil {
			return nil, fmt.Errorf("failed to get service token: %w", err)
		}
		cfg.Hp.ServiceToken = serviceToken

		optionValue, err := common.GetOptionValue("LITEAPI_API_KEY", cfg.Hp)
		if err != nil {
			return nil, fmt.Errorf("failed to get option value: %w", err)
		}

		if key, ok := optionValue.(string); ok {
			if key != "" {
				return NewLiteApiSdk(key), nil
			}
		}

		return nil, fmt.Errorf("invalid option value")

	}
}

var m = map[string]string{
	"/hotels/rates":                      "https://api.liteapi.travel/v3.0",
	"/hotels/min-rates":                  "https://api.liteapi.travel/v3.0",
	"/rates/prebook":                     "https://book.liteapi.travel/v3.0",
	"/rates/book":                        "https://book.liteapi.travel/v3.0",
	"/rates/cancel":                      "https://book.liteapi.travel/v3.0",
	"/rates/get":                         "https://book.liteapi.travel/v3.0",
	"/rates/get-list":                    "https://book.liteapi.travel/v3.0",
	"/rates/get-details":                 "https://book.liteapi.travel/v3.0",
	"/rates/get-reviews":                 "https://book.liteapi.travel/v3.0",
	"/rates/get-data-reviews":            "https://book.liteapi.travel/v3.0",
	"/rates/get-guests-ids":              "https://book.liteapi.travel/v3.0",
	"/rates/get-guests-bookings":         "https://book.liteapi.travel/v3.0",
	"/rates/get-voucher-by-id":           "https://book.liteapi.travel/v3.0",
	"/rates/get-vouchers":                "https://book.liteapi.travel/v3.0",
	"/rates/create-voucher":              "https://book.liteapi.travel/v3.0",
	"/rates/update-voucher":              "https://book.liteapi.travel/v3.0",
	"/rates/update-voucher-status":       "https://book.liteapi.travel/v3.0",
	"/rates/get-loyalty":                 "https://book.liteapi.travel/v3.0",
	"/rates/enable-loyalty":              "https://book.liteapi.travel/v3.0",
	"/rates/update-loyalty":              "https://book.liteapi.travel/v3.0",
	"/rates/retrieve-weekly-analytics":   "https://book.liteapi.travel/v3.0",
	"/rates/retrieve-analytics-report":   "https://book.liteapi.travel/v3.0",
	"/rates/retrieve-market-analytics":   "https://book.liteapi.travel/v3.0",
	"/rates/retrieve-most-booked-hotels": "https://book.liteapi.travel/v3.0",
	"/rates/request":                     "https://book.liteapi.travel/v3.0",
	"/rates/get-loyalties":               "https://book.liteapi.travel/v3.0",
	"/rates/enable-loyalties":            "https://book.liteapi.travel/v3.0",
	"/rates/update-loyalties":            "https://book.liteapi.travel/v3.0",
}
