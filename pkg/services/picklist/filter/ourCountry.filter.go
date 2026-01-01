package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type OurCountryFilter struct {
	SearchWord string `json:"searchWord"`
	Page       int64  `json:"page"`
	Size       int64  `json:"size"`
}

func (f OurCountryFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	if f.SearchWord != "" {
		m["name"] = bson.M{
			"$regex":   f.SearchWord,
			"$options": "i",
		}
	}

	return []bson.M{
		{
			"$match": m,
		},
	}
}
