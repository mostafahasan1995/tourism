package filter

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type DestinationFilter struct {
	SearchWord string `json:"searchWord"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
}

func NewDestinationFilter(query string) (*DestinationFilter, error) {
	f := &DestinationFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (f DestinationFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	if f.SearchWord != "" {
		m["$or"] = []bson.M{
			{
				"name": bson.M{
					"$regex":   f.SearchWord,
					"$options": "i",
				},
			},
			{
				"country": bson.M{
					"$regex":   f.SearchWord,
					"$options": "i",
				},
			},
		}
	}

	return []bson.M{
		{
			"$match": m,
		},
	}
}
