package helpers

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type Filtering interface {
	BuildPipeline(m bson.M) []bson.M
}

func ParseFilters[T Filtering](query any) (*T, error) {
	t, ok := query.(*T)
	if ok {
		return t, nil
	} else if str, ok := query.(string); ok {
		if str == "" {
			return new(T), nil
		}

		t := new(T)
		if err := json.Unmarshal([]byte(str), t); err != nil {
			return nil, err
		}
		return t, nil
	} else {
		return new(T), nil
	}
}
