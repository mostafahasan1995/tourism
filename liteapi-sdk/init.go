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
