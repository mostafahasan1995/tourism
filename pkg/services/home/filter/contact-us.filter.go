package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type ContactUsFilter struct {
	Page   int    `bson:"page" json:"page"`
	Size   int    `bson:"size" json:"size"`
	Status string `bson:"status" json:"status"`
}

func (f *ContactUsFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add status filter if provided
	if f.Status != "" {
		filterConditions = append(filterConditions, bson.M{"status": f.Status})
	}

	return bson.M{"$and": filterConditions}
}
