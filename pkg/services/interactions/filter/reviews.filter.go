package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsFilter struct {
	Status string `bson:"status" json:"status"`

	// Entity filters
	Type string             `bson:"type" json:"type"` // hotel, program, destination, etc.
	Ref  primitive.ObjectID `bson:"ref" json:"ref"`   // Reference ID
	Refs []string           `bson:"refs" json:"refs"` // Multiple reference IDs

	// Customer filters
	Customer string `bson:"customer" json:"customer"` // customer or agent
	Username string `bson:"username" json:"username"`
	UserId   string `bson:"userId" json:"userId"`

	// Program filters (legacy support)
	ProgramId primitive.ObjectID `bson:"programId" json:"programId"`

	// Independent destination and countries
	Destination string `bson:"destination" json:"destination"`
	Country     string `bson:"country" json:"country"`

	// Rating filters
	MinRating float64 `bson:"minRating" json:"minRating"`
	MaxRating float64 `bson:"maxRating" json:"maxRating"`

	// Content filters
	Search    string `bson:"search" json:"search"` // Search in description and advice
	HasImages *bool  `bson:"hasImages" json:"hasImages"`
	HasAdvice *bool  `bson:"hasAdvice" json:"hasAdvice"`
}

func (f *ReviewsFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add status filter if provided
	if f.Status != "" {
		filterConditions = append(filterConditions, bson.M{"status": f.Status})
	}

	// Entity type filter
	if f.Type != "" {
		filterConditions = append(filterConditions, bson.M{"type": f.Type})
	}

	// Single reference filter
	if !f.Ref.IsZero() {
		filterConditions = append(filterConditions, bson.M{"ref": f.Ref})
	}

	// Multiple references filter
	if len(f.Refs) > 0 {
		var refIds []primitive.ObjectID
		for _, refStr := range f.Refs {
			if refId, err := primitive.ObjectIDFromHex(refStr); err == nil {
				refIds = append(refIds, refId)
			}
		}
		if len(refIds) > 0 {
			filterConditions = append(filterConditions, bson.M{"ref": bson.M{"$in": refIds}})
		}
	}

	// Customer type filter
	if f.Customer != "" {
		filterConditions = append(filterConditions, bson.M{"customer": f.Customer})
	}

	// Username filter
	if f.Username != "" {
		filterConditions = append(filterConditions, bson.M{"username": bson.M{"$regex": f.Username, "$options": "i"}})
	}

	// UserId filter (legacy support)
	if f.UserId != "" {
		filterConditions = append(filterConditions, bson.M{"userId": f.UserId})
	}

	// Program filter (legacy support)
	if !f.ProgramId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"programId": f.ProgramId})
	}

	// Independent destination filter
	if f.Destination != "" {
		filterConditions = append(filterConditions, bson.M{"destination": bson.M{"$regex": f.Destination, "$options": "i"}})
	}

	// Independent country filter
	if f.Country != "" {
		filterConditions = append(filterConditions, bson.M{"countries": bson.M{"$in": []string{f.Country}}})
	}

	// Rating range filters
	if f.MinRating > 0 {
		filterConditions = append(filterConditions, bson.M{"value": bson.M{"$gte": f.MinRating}})
	}
	if f.MaxRating > 0 {
		filterConditions = append(filterConditions, bson.M{"value": bson.M{"$lte": f.MaxRating}})
	}

	// Content search filter
	if f.Search != "" {
		filterConditions = append(filterConditions, bson.M{
			"$or": []bson.M{
				{"description": bson.M{"$regex": f.Search, "$options": "i"}},
				{"adviceForTravelers": bson.M{"$regex": f.Search, "$options": "i"}},
				{"text": bson.M{"$regex": f.Search, "$options": "i"}}, // Legacy field
			},
		})
	}

	// Images filter
	if f.HasImages != nil {
		if *f.HasImages {
			filterConditions = append(filterConditions, bson.M{"images": bson.M{"$exists": true, "$ne": []interface{}{}}})
		} else {
			filterConditions = append(filterConditions, bson.M{
				"$or": []bson.M{
					{"images": bson.M{"$exists": false}},
					{"images": []interface{}{}},
				},
			})
		}
	}

	// Advice filter
	if f.HasAdvice != nil {
		if *f.HasAdvice {
			filterConditions = append(filterConditions, bson.M{"adviceForTravelers": bson.M{"$exists": true, "$ne": ""}})
		} else {
			filterConditions = append(filterConditions, bson.M{
				"$or": []bson.M{
					{"adviceForTravelers": bson.M{"$exists": false}},
					{"adviceForTravelers": ""},
				},
			})
		}
	}

	return bson.M{"$and": filterConditions}
}
