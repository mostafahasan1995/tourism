package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type TrustedPartnersFilter struct {
	Title    string `bson:"title" json:"title"`
	IsActive *bool  `bson:"isActive" json:"isActive"`

	Page int `bson:"page" json:"page"`
	Size int `bson:"size" json:"size"`
}

func (f *TrustedPartnersFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	if f.Title != "" {
		filterConditions = append(filterConditions, bson.M{"title": bson.M{"$regex": f.Title, "$options": "i"}})
	}

	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	return bson.M{"$and": filterConditions}
}
