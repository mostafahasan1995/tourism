package filter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type HotelsFilter struct {
	SearchWord             string `bson:"searchWord" json:"searchWord"`
	IsDisplayInPerfectStay *bool `bson:"isDisplayInPerfectStay" json:"isDisplayInPerfectStay"`

	HotelTypes             []string                     `bson:"hotelTypes" json:"hotelTypes"`
	RoomAmenities          []string                     `bson:"roomAmenities" json:"roomAmenities"`
	NearbyAttractions      []string                     `bson:"nearbyAttractions" json:"nearbyAttractions"`
	Locations              []string                     `bson:"locations" json:"locations"`
	CheckInAndCheckOut     CheckInAndCheckOutFilter     `bson:"checkInAndCheckOut" json:"checkInAndCheckOut"`
	PriceRange             PriceRangeFilter             `bson:"priceRange" json:"priceRange"`
	Ratings                float64                          `bson:"ratings" json:"ratings"`
	Page                   int                          `bson:"page" json:"page"`
	Size                   int                          `bson:"size" json:"size"`
	DistanceFromCityCenter DistanceFromCityCenterFilter `bson:"distanceFromCityCenter" json:"distanceFromCityCenter"`
}

func (f *HotelsFilter) ToBsonFilter() bson.M {
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

	if len(f.HotelTypes) > 0 {
		filterConditions = append(filterConditions, bson.M{"hotelType": bson.M{"$in": f.HotelTypes}})
	}

	if len(f.Locations) > 0 {
		filterConditions = append(filterConditions, bson.M{"location": bson.M{"$in": f.Locations}})
	}
	if f.IsDisplayInPerfectStay != nil {
		filterConditions = append(filterConditions, bson.M{"isDisplayInPerfectStay": *f.IsDisplayInPerfectStay})
	}

	// if f.CheckInAndCheckOut != "" {
	// 	filterConditions = append(filterConditions, bson.M{"checkInAndCheckOut": f.CheckInAndCheckOut})
	// }

	if len(f.RoomAmenities) > 0 {
		filterConditions = append(filterConditions, bson.M{"roomAmenities": bson.M{"$in": f.RoomAmenities}})
	}

	if len(f.NearbyAttractions) > 0 {
		filterConditions = append(filterConditions, bson.M{"nearbyAttractions": bson.M{"$in": f.NearbyAttractions}})
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

	if f.Ratings != 0 {
		filterConditions = append(filterConditions, bson.M{"ratings": f.Ratings})
	}

	if f.DistanceFromCityCenter.GreaterThan > 0 || f.DistanceFromCityCenter.LessThan > 0 {
		distanceFromCityFilter := bson.M{}
		if f.DistanceFromCityCenter.GreaterThan > 0 {
			distanceFromCityFilter["$gte"] = f.DistanceFromCityCenter.GreaterThan
		}
		if f.DistanceFromCityCenter.LessThan > 0 {
			distanceFromCityFilter["$lte"] = f.DistanceFromCityCenter.LessThan
		}
		filterConditions = append(filterConditions, bson.M{"distanceFromCityCenter": distanceFromCityFilter})
	}
	return bson.M{"$and": filterConditions}
}

type CheckInAndCheckOutFilter struct {
	From time.Time `bson:"from" json:"from"`
	To   time.Time `bson:"to" json:"to"`
}

type DistanceFromCityCenterFilter struct {
	LessThan    int `bson:"lessThan" json:"lessThan"`
	GreaterThan int `bson:"greaterThan" json:"greaterThan"`
}
