package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

// CustomerPersonaFilter defines filter criteria for customer personas
type CustomerPersonaFilter struct {
	Title      string `bson:"title" json:"title"`
	IsActive   *bool  `bson:"isActive" json:"isActive"`
	Preference string `bson:"preference" json:"preference"` // Search across different preference fields
	ProgramId  string `bson:"programId" json:"programId"`
}

// ToBsonFilter converts the filter struct to a MongoDB filter
func (f *CustomerPersonaFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	if f.Title != "" {
		filterConditions = append(filterConditions, bson.M{"title": bson.M{"$regex": f.Title, "$options": "i"}})
	}

	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	if f.Preference != "" {
		// Search across all preference fields
		preferenceFilter := bson.M{"$or": []bson.M{
			{"explorationPreferences": bson.M{"$elemMatch": bson.M{"$regex": f.Preference, "$options": "i"}}},
			{"relaxationPreferences": bson.M{"$elemMatch": bson.M{"$regex": f.Preference, "$options": "i"}}},
			{"travelMustHaves": bson.M{"$elemMatch": bson.M{"$regex": f.Preference, "$options": "i"}}},
		}}
		filterConditions = append(filterConditions, preferenceFilter)
	}

	if f.ProgramId != "" {
		filterConditions = append(filterConditions, bson.M{"programId": f.ProgramId})
	}

	return bson.M{"$and": filterConditions}
}
