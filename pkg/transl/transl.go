package transl

import (
	"context"
	"encoding/json"
	"fmt"
	"larsa-tourism-microservices/pkg/util"
)

type Localizable[T any] map[string]T

func (l Localizable[T]) MarshalJSON(ctx context.Context) ([]byte, error) {
	fmt.Println("calling custom marshaler")
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	lang := cfg.Lang

	content, exists := l[lang]
	if exists {
		return json.Marshal(content)
	} else {
		return json.Marshal(l["en"])
	}
}

func (l *Localizable[T]) UnmarshalJSON(ctx context.Context, data []byte) error {

	fmt.Println("calling custom unmarshaler")

	var result map[string]T
	err := json.Unmarshal(data, &result)

	if err != nil {
		var t T
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		result = map[string]T{
			"en": t,
		}
	}

	*l = result

	return nil
}
