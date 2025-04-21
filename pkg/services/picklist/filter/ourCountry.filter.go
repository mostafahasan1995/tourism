package filter

import (

	"go.mongodb.org/mongo-driver/bson"
)

type OurCountryFilter struct {
	SearchWord             string                       `bson:"searchWord" json:"searchWord"`

	Page                   int                          `bson:"page" json:"page"`
	Size                   int                          `bson:"size" json:"size"`
}

func (f *OurCountryFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}
	if f.SearchWord != "" {
		filterConditions = append(filterConditions, bson.M{
			"name": bson.M{
				"$regex": f.SearchWord,
				"$options": "i",
			},
		})
	}
	return bson.M{"$and": filterConditions}
}
