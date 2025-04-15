package filter

import (

	"go.mongodb.org/mongo-driver/bson"
)

type OurCountryFilter struct {

	Page                   int                          `bson:"page" json:"page"`
	Size                   int                          `bson:"size" json:"size"`
}

func (f *OurCountryFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	return bson.M{"$and": filterConditions}
}
