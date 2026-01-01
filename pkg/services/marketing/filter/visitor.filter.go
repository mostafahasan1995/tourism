package filter

import (
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VisitorFilter struct {
	Page           int                `bson:"page" json:"page"`
	Size           int                `bson:"size" json:"size"`
	HotelId        primitive.ObjectID `bson:"hotelId,omitempty" json:"hotelId,omitempty"`
	FullName       string             `bson:"fullName" json:"fullName"`
	Nationality    string             `bson:"nationality" json:"nationality"`
	Email          string             `bson:"email" json:"email"`
	Phone          string             `bson:"phone" json:"phone"`
	Interests      []string           `bson:"interests" json:"interests"`
	Tags           []string           `bson:"tags" json:"tags"`
	Source         string             `bson:"source" json:"source"`
	IsActive       *bool              `bson:"isActive,omitempty" json:"isActive,omitempty"`
	IsVIP          *bool              `bson:"isVIP,omitempty" json:"isVIP,omitempty"`
	VisitDateFrom  *time.Time         `bson:"visitDateFrom,omitempty" json:"visitDateFrom,omitempty"`
	VisitDateTo    *time.Time         `bson:"visitDateTo,omitempty" json:"visitDateTo,omitempty"`
	MinVisitCount  *int               `bson:"minVisitCount,omitempty" json:"minVisitCount,omitempty"`
	MaxVisitCount  *int               `bson:"maxVisitCount,omitempty" json:"maxVisitCount,omitempty"`
	MinTotalSpent  *float64           `bson:"minTotalSpent,omitempty" json:"minTotalSpent,omitempty"`
	MaxTotalSpent  *float64           `bson:"maxTotalSpent,omitempty" json:"maxTotalSpent,omitempty"`
	SearchText     string             `bson:"searchText" json:"searchText"`
	Name           *string            `json:"name,omitempty"`
	Country        *string            `json:"country,omitempty"`
	RegistrationId *string            `json:"registrationId,omitempty"`
	Status         *string            `json:"status,omitempty"`
	DateFrom       *time.Time         `json:"dateFrom,omitempty"`
	DateTo         *time.Time         `json:"dateTo,omitempty"`
}

func (f VisitorFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A

	if f.Name != nil && *f.Name != "" {
		ands = append(ands, bson.M{"name": bson.M{"$regex": *f.Name, "$options": "i"}})
	}

	if f.Email != "" {
		ands = append(ands, bson.M{"email": bson.M{"$regex": f.Email, "$options": "i"}})
	}

	if f.Phone != "" {
		ands = append(ands, bson.M{"phone": bson.M{"$regex": f.Phone, "$options": "i"}})
	}

	if f.Country != nil && *f.Country != "" {
		ands = append(ands, bson.M{"country": bson.M{"$regex": *f.Country, "$options": "i"}})
	}

	if f.RegistrationId != nil && *f.RegistrationId != "" {
		ands = append(ands, bson.M{"registrationId": bson.M{"$regex": *f.RegistrationId, "$options": "i"}})
	}

	if f.Status != nil && *f.Status != "" {
		ands = append(ands, bson.M{"status": *f.Status})
	}

	dateFilter := bson.M{}
	if f.DateFrom != nil {
		dateFilter["$gte"] = f.DateFrom
	}

	if f.DateTo != nil {
		dateFilter["$lte"] = f.DateTo
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
func (f *VisitorFilter) ParseQueryParams(q url.Values) {
	// Name filter
	if name := q.Get("name"); name != "" {
		f.Name = &name
	}

	// Email filter
	if email := q.Get("email"); email != "" {
		f.Email = email
	}

	// Phone filter
	if phone := q.Get("phone"); phone != "" {
		f.Phone = phone
	}

	// Country filter
	if country := q.Get("country"); country != "" {
		f.Country = &country
	}

	// Registration ID filter
	if regId := q.Get("registrationId"); regId != "" {
		f.RegistrationId = &regId
	}

	// Status filter
	if status := q.Get("status"); status != "" {
		f.Status = &status
	}

	// Date filters
	if from := q.Get("dateFrom"); from != "" {
		date, err := time.Parse("2006-01-02", from)
		if err == nil {
			f.DateFrom = &date
		}
	}

	if to := q.Get("dateTo"); to != "" {
		date, err := time.Parse("2006-01-02", to)
		if err == nil {
			// Set to end of day
			date = date.Add(24*time.Hour - time.Second)
			f.DateTo = &date
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

func (f VisitorFilter) ToBsonFilter() bson.M {
	// Keep it super simple - just exclude trashed items
	filter := bson.M{"trash": false}
	return filter
}
