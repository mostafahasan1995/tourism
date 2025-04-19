package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type TourismProgramFilter struct {
	Destinations []string `bson:"destinations" json:"destinations"`
	TravelTypes  []string `bson:"travelTypes" json:"travelTypes"`
	Interests    []string `bson:"interests" json:"interests"`

	Durations      []int            `bson:"durations" json:"durations"`
	GroupSize      string           `bson:"groupSize" json:"groupSize"`
	SearchWord     string           `bson:"searchWord" json:"searchWord"`
	PriceRange     PriceRangeFilter `bson:"priceRange" json:"priceRange"`
	ActivitiesName []string         `bson:"activitiesName" json:"activitiesName"`
	Page           int              `bson:"page" json:"page"`
	Size           int              `bson:"size" json:"size"`
}

func (f *TourismProgramFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}
	if f.SearchWord != "" {
		filterConditions = append(filterConditions, bson.M{
			"name": bson.M{
				"$regex":   f.SearchWord,
				"$options": "i",
			},
		})
	}

	if len(f.TravelTypes) > 0 {
		filterConditions = append(filterConditions, bson.M{"travelType": bson.M{"$in": f.TravelTypes}})
	}
	if len(f.Interests) > 0 {
		filterConditions = append(filterConditions, bson.M{"interests": bson.M{"$in": f.Interests}})
	}

	if len(f.Destinations) > 0 {
		filterConditions = append(filterConditions, bson.M{"destination": bson.M{"$in": f.Destinations}})
	}

	if len(f.Durations) > 0 {
		filterConditions = append(filterConditions, bson.M{"duration": bson.M{"$in": f.Durations}})
	}

	if f.PriceRange.From > 0 || f.PriceRange.To > 0 {
		priceFilter := bson.M{}
		if f.PriceRange.From > 0 {
			priceFilter["$gte"] = f.PriceRange.From
		}
		if f.PriceRange.To > 0 {
			priceFilter["$lte"] = f.PriceRange.To
		}
		filterConditions = append(filterConditions, bson.M{"price": priceFilter})
	}

	if len(f.ActivitiesName) > 0 {
		filterConditions = append(filterConditions, bson.M{
			"activities.name": bson.M{"$in": f.ActivitiesName},
		})
	}

	return bson.M{"$and": filterConditions}
}
