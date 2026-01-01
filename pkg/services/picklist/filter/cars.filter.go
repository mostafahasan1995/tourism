package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type CarsFilter struct {
	SearchWord string `json:"searchWord"`
	CarType    string `json:"carType"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

func (f CarsFilter) BuildPipeline(m bson.M) []bson.M {

	// Search by car type if searchWord is provided
	if f.SearchWord != "" {
		m["carType"] = bson.M{
			"$regex":   f.SearchWord,
			"$options": "i",
		}
	}

	if f.CarType != "" {
		m["carType.en"] = f.CarType
	}

	return []bson.M{
		{
			"$match": m,
		},
	}
}
