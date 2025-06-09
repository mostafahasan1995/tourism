package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExhibitorProfileFilter struct {
	Page         int                `bson:"page" json:"page"`
	Size         int                `bson:"size" json:"size"`
	HotelId      primitive.ObjectID `bson:"hotelId,omitempty" json:"hotelId,omitempty"`
	HotelName    string             `bson:"hotelName" json:"hotelName"`
	IsActive     *bool              `bson:"isActive,omitempty" json:"isActive,omitempty"`
	IsPublished  *bool              `bson:"isPublished,omitempty" json:"isPublished,omitempty"`
	Rating       *float64           `bson:"rating,omitempty" json:"rating,omitempty"`
	MinRating    *float64           `bson:"minRating,omitempty" json:"minRating,omitempty"`
	MaxRating    *float64           `bson:"maxRating,omitempty" json:"maxRating,omitempty"`
	PropertyType string             `bson:"propertyType" json:"propertyType"`
}

func (f *ExhibitorProfileFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add hotel ID filter if provided
	if !f.HotelId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"hotelId": f.HotelId})
	}

	// Add hotel name search if provided
	if f.HotelName != "" {
		filterConditions = append(filterConditions, bson.M{
			"heroSection.hotelName": bson.M{"$regex": f.HotelName, "$options": "i"},
		})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	// Add published status filter if provided
	if f.IsPublished != nil {
		filterConditions = append(filterConditions, bson.M{"isPublished": *f.IsPublished})
	}

	// Add exact rating filter if provided
	if f.Rating != nil {
		filterConditions = append(filterConditions, bson.M{"heroSection.rating": *f.Rating})
	}

	// Add rating range filter if provided
	if f.MinRating != nil && f.MaxRating != nil {
		filterConditions = append(filterConditions, bson.M{
			"heroSection.rating": bson.M{
				"$gte": *f.MinRating,
				"$lte": *f.MaxRating,
			},
		})
	} else if f.MinRating != nil {
		filterConditions = append(filterConditions, bson.M{
			"heroSection.rating": bson.M{"$gte": *f.MinRating},
		})
	} else if f.MaxRating != nil {
		filterConditions = append(filterConditions, bson.M{
			"heroSection.rating": bson.M{"$lte": *f.MaxRating},
		})
	}

	// Add property type filter if provided
	if f.PropertyType != "" {
		filterConditions = append(filterConditions, bson.M{
			"heroSection.propertyType": bson.M{"$regex": f.PropertyType, "$options": "i"},
		})
	}

	return bson.M{"$and": filterConditions}
}
