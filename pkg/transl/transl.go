package transl

import (
	"context"
	"larsa-tourism-microservices/pkg/util"

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
)

type Localizable[T any] map[string]T

func (l *Localizable[T]) MarshalJSON(ctx context.Context) ([]byte, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	lang := cfg.Lang
	if lang == "" {
		return json.Marshal(*l)
	} else {
		content, exists := (*l)[lang]
		if exists {
			return json.Marshal(content)
		} else {
			return json.Marshal((*l)["en"])
		}
	}

}

func (l *Localizable[T]) UnmarshalJSON(ctx context.Context, data []byte) error {
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

func (l *Localizable[T]) GetContentByLang(lang string) T {
	content, exists := (*l)[lang]
	if exists {
		return content
	} else {
		return (*l)["en"]
	}
}

// todo: no need for this function
func (l *Localizable[T]) UnmarshalBSON(data []byte) error {

	var result map[string]T
	if err := bson.Unmarshal(data, &result); err != nil {
		var str any = string(data)

		result = map[string]T{
			"en": str.(T),
		}
	}

	*l = result

	return nil
}
