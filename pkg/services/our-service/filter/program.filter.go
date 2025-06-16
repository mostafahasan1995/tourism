package filter

import (
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramFilter struct {
	// Basic filters
	CustomerName *string             `json:"customerName"`
	Title        *string             `json:"title"`
	ServiceType  []enums.ServiceType `json:"serviceType"`
	TravelReqId  *primitive.ObjectID `json:"travelReqId"`
	CustomerId   *primitive.ObjectID `json:"customerId"`
	AgentId      *primitive.ObjectID `json:"agentId"`
	Status       *string             `json:"status"`
	Package      *primitive.ObjectID `json:"package"`
	ProgramType  *string             `json:"programType"`
	Source       *string             `json:"source"`
	Company      *string             `json:"company"`
	Coordinator  *string             `json:"coordinator"`
	Purpose      *string             `json:"purpose"`
	StartDate    *time.Time          `json:"startDate"`
	EndDate      *time.Time          `json:"endDate"`
	GroupSize    *enums.GroupSize    `json:"groupSize"`
	CreatedBy    *primitive.ObjectID `json:"createdBy"`
	UpdatedBy    *primitive.ObjectID `json:"updatedBy"`

	// Website-specific filters
	Destinations   []primitive.ObjectID `json:"destinations"`   // Filter by destination IDs
	Activities     []primitive.ObjectID `json:"activities"`     // Filter by activity IDs
	MinPrice       *float64             `json:"minPrice"`       // Minimum price filter
	MaxPrice       *float64             `json:"maxPrice"`       // Maximum price filter
	Accommodation  []string             `json:"accommodation"`  // Filter by accommodation types (luxury, family, adventure, etc.)
	Transportation []string             `json:"transportation"` // Filter by transportation types
	Meals          []string             `json:"meals"`          // Filter by meal types
	Interests      []string             `json:"interests"`      // Filter by interests (hiking, diving, sightseeing, safari, adventure, relaxation, outdoors, food)
	Duration       []int                `json:"duration"`       // Filter by trip duration in days
	Recommended    *bool                `json:"recommended"`    // Filter for recommended programs
	ShowInWebsite  *bool                `json:"showInWebsite"`  // Filter for programs that should be shown on website
}

func (f ProgramFilter) BuildPipeline(m bson.M) []bson.M {
	// Add basic filters to the match stage
	if f.Title != nil {
		m["title"] = bson.M{"$regex": *f.Title, "$options": "i"}
	}

	if len(f.ServiceType) > 0 {
		serviceTypes := make([]string, len(f.ServiceType))
		for i, st := range f.ServiceType {
			serviceTypes[i] = string(st)
		}
		m["serviceType"] = bson.M{"$in": serviceTypes}
	}

	if f.TravelReqId != nil {
		m["travelReqId"] = *f.TravelReqId
	}

	if f.CustomerId != nil {
		m["customerId"] = *f.CustomerId
	}

	if f.AgentId != nil {
		m["agentId"] = *f.AgentId
	}

	if f.Status != nil {
		m["status"] = *f.Status
	}

	if f.Package != nil {
		m["package"] = *f.Package
	}

	if f.ProgramType != nil {
		m["programType"] = *f.ProgramType
	}

	if f.Source != nil {
		m["source"] = bson.M{"$regex": *f.Source, "$options": "i"}
	}

	if f.Company != nil {
		m["company"] = bson.M{"$regex": *f.Company, "$options": "i"}
	}

	if f.Coordinator != nil {
		m["coordinator"] = bson.M{"$regex": *f.Coordinator, "$options": "i"}
	}

	if f.Purpose != nil {
		m["purpose"] = bson.M{"$regex": *f.Purpose, "$options": "i"}
	}

	if f.StartDate != nil {
		m["startDate"] = bson.M{"$gte": *f.StartDate}
	}

	if f.EndDate != nil {
		m["endDate"] = bson.M{"$lte": *f.EndDate}
	}

	if f.GroupSize != nil {
		m["groupSize"] = *f.GroupSize
	}

	if f.CreatedBy != nil {
		m["createdBy"] = *f.CreatedBy
	}

	if f.UpdatedBy != nil {
		m["updatedBy"] = *f.UpdatedBy
	}

	// Website-specific filters for general programs
	if len(f.Destinations) > 0 {
		m["generalType.destinations.from"] = bson.M{"$in": f.Destinations}
	}

	if len(f.Activities) > 0 {
		m["generalType.activities"] = bson.M{"$in": f.Activities}
	}

	if len(f.Accommodation) > 0 {
		m["generalType.includes.accommodation"] = bson.M{"$in": f.Accommodation}
	}

	if len(f.Transportation) > 0 {
		m["generalType.includes.transportation"] = bson.M{"$in": f.Transportation}
	}

	if len(f.Meals) > 0 {
		m["generalType.includes.meals"] = bson.M{"$in": f.Meals}
	}

	// Price range filters
	priceFilter := bson.M{}
	if f.MinPrice != nil {
		priceFilter["$gte"] = *f.MinPrice
	}
	if f.MaxPrice != nil {
		priceFilter["$lte"] = *f.MaxPrice
	}
	if len(priceFilter) > 0 {
		m["generalType.pricing.person.price"] = priceFilter
	}

	// Show in website filter
	if f.ShowInWebsite != nil {
		m["generalType.pricing.person.showInWebsite"] = *f.ShowInWebsite
	}

	// Duration filter (calculate from start and end dates)
	if len(f.Duration) > 0 {
		m["$expr"] = bson.M{
			"$in": []interface{}{
				bson.M{
					"$divide": []interface{}{
						bson.M{
							"$subtract": []interface{}{"$endDate", "$startDate"},
						},
						1000 * 60 * 60 * 24, // Convert milliseconds to days
					},
				},
				f.Duration,
			},
		}
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	// Handle customer name lookup
	if f.CustomerName != nil {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "customers",
					"localField":   "customerId",
					"foreignField": "_id",
					"as":           "customer",
				},
			},
			bson.M{
				"$match": bson.M{
					"customer.name": bson.M{"$regex": *f.CustomerName, "$options": "i"},
				},
			},
			bson.M{
				"$project": bson.M{
					"customer": 0,
				},
			},
		)
	}

	// Handle interests filter by looking up activities
	if len(f.Interests) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "activities",
					"localField":   "generalType.activities",
					"foreignField": "_id",
					"as":           "activityDetails",
				},
			},
			bson.M{
				"$match": bson.M{
					"$or": []bson.M{
						{"activityDetails.category": bson.M{"$in": f.Interests}},
						{"activityDetails.tags": bson.M{"$in": f.Interests}},
						{"activityDetails.name": bson.M{"$regex": "(?i)(" + joinStrings(f.Interests, "|") + ")"}},
					},
				},
			},
			bson.M{
				"$project": bson.M{
					"activityDetails": 0,
				},
			},
		)
	}

	return pipeline
}

// Helper function to join strings for regex
func joinStrings(strs []string, separator string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}
	return result
}
