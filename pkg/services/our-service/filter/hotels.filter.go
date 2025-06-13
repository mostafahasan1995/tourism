package filter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type HotelsFilter struct {
	SearchWord             string `json:"searchWord"`
	IsDisplayInPerfectStay *bool  `json:"isDisplayInPerfectStay"`
	HotelType              string `json:"hotelType"`
	//
	SortBy          []string `json:"sortBy"`
	ReviewScore     int      `json:"reviewScore"`
	MealOptions     []string `json:"mealOptions"`
	BedType         []string `json:"bedType"`
	HotelDacilities []string `json:"hotelDacilities"`
	RoomType        []string `json:"roomType"`
	PetFriendly     string   `json:"petFriendly"`

	//
	HotelTypes             []string                     `json:"hotelTypes"`
	RoomAmenities          []string                     `json:"roomAmenities"`
	NearbyAttractions      []string                     `json:"nearbyAttractions"`
	Locations              []string                     `json:"locations"`
	CheckInAndCheckOut     CheckInAndCheckOutFilter     `json:"checkInAndCheckOut"`
	PriceRange             PriceRangeFilter             `json:"priceRange"`
	Ratings                float64                      `json:"ratings"`
	DistanceFromCityCenter DistanceFromCityCenterFilter `json:"distanceFromCityCenter"`
}

func (f HotelsFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	if f.SearchWord != "" {
		m["name"] = bson.M{
			"$regex":   f.SearchWord,
			"$options": "i",
		}
	}

	// Handle both single hotelType and array of hotelTypes
	if f.HotelType != "" {
		m["hotelType"] = bson.M{
			"$regex":   f.HotelType,
			"$options": "i",
		}
	} else if len(f.HotelTypes) > 0 {
		m["hotelType"] = bson.M{"$in": f.HotelTypes}
	}

	if len(f.Locations) > 0 {
		m["location"] = bson.M{"$in": f.Locations}
	}

	if f.IsDisplayInPerfectStay != nil {
		m["isDisplayInPerfectStay"] = *f.IsDisplayInPerfectStay
	}

	if len(f.RoomAmenities) > 0 {
		m["roomAmenities"] = bson.M{"$in": f.RoomAmenities}
	}

	if len(f.NearbyAttractions) > 0 {
		m["nearbyAttractions"] = bson.M{"$in": f.NearbyAttractions}
	}

	if priceFilter := f.PriceRange.BuildPriceFilter(); priceFilter != nil {
		m["price"] = priceFilter
	}

	if f.Ratings != 0 {
		m["ratings"] = f.Ratings
	}

	if f.DistanceFromCityCenter.GreaterThan > 0 || f.DistanceFromCityCenter.LessThan > 0 {
		distanceFromCityFilter := bson.M{}
		if f.DistanceFromCityCenter.GreaterThan > 0 {
			distanceFromCityFilter["$gte"] = f.DistanceFromCityCenter.GreaterThan
		}
		if f.DistanceFromCityCenter.LessThan > 0 {
			distanceFromCityFilter["$lte"] = f.DistanceFromCityCenter.LessThan
		}
		m["distanceFromCityCenter"] = distanceFromCityFilter
	}

	return []bson.M{
		{"$match": m},
	}
}

// ToBsonFilter converts the filter to a MongoDB filter
func (f HotelsFilter) ToBsonFilter() bson.M {
	filter := bson.M{"trash": bson.M{"$ne": true}}

	if f.SearchWord != "" {
		filter["name"] = bson.M{
			"$regex":   f.SearchWord,
			"$options": "i",
		}
	}

	// Handle both single hotelType and array of hotelTypes
	if f.HotelType != "" {
		filter["hotelType"] = bson.M{
			"$regex":   f.HotelType,
			"$options": "i",
		}
	} else if len(f.HotelTypes) > 0 {
		filter["hotelType"] = bson.M{"$in": f.HotelTypes}
	}

	if len(f.Locations) > 0 {
		filter["location"] = bson.M{"$in": f.Locations}
	}

	if f.IsDisplayInPerfectStay != nil {
		filter["isDisplayInPerfectStay"] = *f.IsDisplayInPerfectStay
	}

	if len(f.RoomAmenities) > 0 {
		filter["roomAmenities"] = bson.M{"$in": f.RoomAmenities}
	}

	if len(f.NearbyAttractions) > 0 {
		filter["nearbyAttractions"] = bson.M{"$in": f.NearbyAttractions}
	}

	if priceFilter := f.PriceRange.BuildPriceFilter(); priceFilter != nil {
		filter["price"] = priceFilter
	}

	if f.Ratings != 0 {
		filter["ratings"] = f.Ratings
	}

	if f.DistanceFromCityCenter.GreaterThan > 0 || f.DistanceFromCityCenter.LessThan > 0 {
		distanceFromCityFilter := bson.M{}
		if f.DistanceFromCityCenter.GreaterThan > 0 {
			distanceFromCityFilter["$gte"] = f.DistanceFromCityCenter.GreaterThan
		}
		if f.DistanceFromCityCenter.LessThan > 0 {
			distanceFromCityFilter["$lte"] = f.DistanceFromCityCenter.LessThan
		}
		filter["distanceFromCityCenter"] = distanceFromCityFilter
	}

	return filter
}

type CheckInAndCheckOutFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type DistanceFromCityCenterFilter struct {
	LessThan    int `json:"lessThan"`
	GreaterThan int `json:"greaterThan"`
}
