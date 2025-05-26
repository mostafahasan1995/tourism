package filter

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type GeneralProgramFilter struct {
}

func NewGeneralProgramFilter(query string) (*GeneralProgramFilter, error) {

	f := &GeneralProgramFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *GeneralProgramFilter) BuildPipeline(m bson.M) []bson.M {

	return []bson.M{
		{
			"$match": m,
		},
	}

}
