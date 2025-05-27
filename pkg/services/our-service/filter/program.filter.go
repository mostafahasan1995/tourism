package filter

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type ProgramFilter struct {
	CustomerName *string `json:"customerName"`
}

func NewProgramFilter(query string) (*ProgramFilter, error) {
	f := &ProgramFilter{}

	if query != "" {
		if err := json.Unmarshal([]byte(query), f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *ProgramFilter) BuildPipeline(m bson.M) []bson.M {

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline

}
