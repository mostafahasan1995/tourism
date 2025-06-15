package filter

import (
	"net/url"
	"strconv"
	"time"

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

func (f ExhibitorProfileFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A

	// Add hotel ID filter if provided
	if !f.HotelId.IsZero() {
		ands = append(ands, bson.M{"hotelId": f.HotelId})
	}

	// Add hotel name search if provided
	if f.HotelName != "" {
		ands = append(ands, bson.M{
			"heroSection.hotelName": bson.M{"$regex": f.HotelName, "$options": "i"},
		})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		ands = append(ands, bson.M{"isActive": *f.IsActive})
	}

	// Add published status filter if provided
	if f.IsPublished != nil {
		ands = append(ands, bson.M{"isPublished": *f.IsPublished})
	}

	// Add exact rating filter if provided
	if f.Rating != nil {
		ands = append(ands, bson.M{"heroSection.rating": *f.Rating})
	}

	// Add rating range filter if provided
	if f.MinRating != nil && f.MaxRating != nil {
		ands = append(ands, bson.M{
			"heroSection.rating": bson.M{
				"$gte": *f.MinRating,
				"$lte": *f.MaxRating,
			},
		})
	} else if f.MinRating != nil {
		ands = append(ands, bson.M{
			"heroSection.rating": bson.M{"$gte": *f.MinRating},
		})
	} else if f.MaxRating != nil {
		ands = append(ands, bson.M{
			"heroSection.rating": bson.M{"$lte": *f.MaxRating},
		})
	}

	// Add property type filter if provided
	if f.PropertyType != "" {
		ands = append(ands, bson.M{
			"heroSection.propertyType": bson.M{"$regex": f.PropertyType, "$options": "i"},
		})
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}

func (f ExhibitorProfileFilter) ToBsonFilter() bson.M {
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

// ExhibitorRequestFilter defines the filter criteria for exhibitor requests
type ExhibitorRequestFilter struct {
	HotelName       *string    `json:"hotelName,omitempty"`
	Location        *string    `json:"location,omitempty"`
	Email           *string    `json:"email,omitempty"`
	Status          *string    `json:"status,omitempty"`
	RequestDateFrom *time.Time `json:"requestDateFrom,omitempty"`
	RequestDateTo   *time.Time `json:"requestDateTo,omitempty"`
	Page            int        `json:"page"`
	Size            int        `json:"size"`
}

func (f ExhibitorRequestFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A

	if f.HotelName != nil && *f.HotelName != "" {
		ands = append(ands, bson.M{"heroSection.hotelName": bson.M{"$regex": *f.HotelName, "$options": "i"}})
	}

	if f.Location != nil && *f.Location != "" {
		ands = append(ands, bson.M{"contactInfo.location": bson.M{"$regex": *f.Location, "$options": "i"}})
	}

	if f.Email != nil && *f.Email != "" {
		ands = append(ands, bson.M{"contactInfo.email": bson.M{"$regex": *f.Email, "$options": "i"}})
	}

	if f.Status != nil && *f.Status != "" {
		ands = append(ands, bson.M{"status": *f.Status})
	}

	dateFilter := bson.M{}
	if f.RequestDateFrom != nil {
		dateFilter["$gte"] = f.RequestDateFrom
	}

	if f.RequestDateTo != nil {
		dateFilter["$lte"] = f.RequestDateTo
	}

	if len(dateFilter) > 0 {
		ands = append(ands, bson.M{"createdAt": dateFilter})
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}

// ParseQueryParams parses URL query parameters into filter values
func (f *ExhibitorRequestFilter) ParseQueryParams(q url.Values) {
	// Hotel name filter
	if hotelName := q.Get("hotelName"); hotelName != "" {
		f.HotelName = &hotelName
	}

	// Location filter
	if location := q.Get("location"); location != "" {
		f.Location = &location
	}

	// Email filter
	if email := q.Get("email"); email != "" {
		f.Email = &email
	}

	// Status filter
	if status := q.Get("status"); status != "" {
		f.Status = &status
	}

	// Date filters
	if from := q.Get("requestDateFrom"); from != "" {
		date, err := time.Parse("2006-01-02", from)
		if err == nil {
			f.RequestDateFrom = &date
		}
	}

	if to := q.Get("requestDateTo"); to != "" {
		date, err := time.Parse("2006-01-02", to)
		if err == nil {
			// Set to end of day
			date = date.Add(24*time.Hour - time.Second)
			f.RequestDateTo = &date
		}
	}

	// Set pagination
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	f.Page = page

	size, _ := strconv.Atoi(q.Get("size"))
	if size < 1 {
		size = 10 // Default size
	}
	f.Size = size
}

// ToBsonFilter converts the filter to a BSON filter for MongoDB
func (f *ExhibitorRequestFilter) ToBsonFilter() bson.M {
	filter := bson.M{"trash": false}

	if f.HotelName != nil && *f.HotelName != "" {
		filter["heroSection.hotelName"] = bson.M{"$regex": *f.HotelName, "$options": "i"}
	}

	if f.Location != nil && *f.Location != "" {
		filter["contactInfo.location"] = bson.M{"$regex": *f.Location, "$options": "i"}
	}

	if f.Email != nil && *f.Email != "" {
		filter["contactInfo.email"] = bson.M{"$regex": *f.Email, "$options": "i"}
	}

	if f.Status != nil && *f.Status != "" {
		filter["status"] = *f.Status
	}

	dateFilter := bson.M{}
	if f.RequestDateFrom != nil {
		dateFilter["$gte"] = f.RequestDateFrom
	}

	if f.RequestDateTo != nil {
		dateFilter["$lte"] = f.RequestDateTo
	}

	if len(dateFilter) > 0 {
		filter["createdAt"] = dateFilter
	}

	return filter
}
