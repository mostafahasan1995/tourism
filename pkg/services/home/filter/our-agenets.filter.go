package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type OurAgentsFilter struct {
	LanguagesSpoken   []string `bson:"languagesSpoken" json:"languagesSpoken"`
	CountriesYouServe []string `bson:"countriesYouServe" json:"countriesYouServe"`

	Page int `bson:"page" json:"page"`
	Size int `bson:"size" json:"size"`
}

func (f *OurAgentsFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	if len(f.LanguagesSpoken) > 0 {
		filterConditions = append(filterConditions, bson.M{"languagesSpoken": bson.M{"$in": f.LanguagesSpoken}})
	}
	if len(f.CountriesYouServe) > 0 {
		filterConditions = append(filterConditions, bson.M{"countriesYouServe": bson.M{"$in": f.CountriesYouServe}})
	}

	return bson.M{"$and": filterConditions}
}
