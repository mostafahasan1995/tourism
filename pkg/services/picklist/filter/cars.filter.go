package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type CarsFilter struct {
	SearchWord string `json:"searchWord"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

func (f CarsFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	// Search by car type if searchWord is provided
	if f.SearchWord != "" {
		m["carType"] = bson.M{
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

// ToBsonFilter returns a bson.M representation of the filter
func (f CarsFilter) ToBsonFilter() bson.M {
	filter := bson.M{"trash": bson.M{"$ne": true}}

	if f.SearchWord != "" {
		filter["carType"] = bson.M{
			"$regex":   f.SearchWord,
			"$options": "i",
		}
	}

	return filter
}
