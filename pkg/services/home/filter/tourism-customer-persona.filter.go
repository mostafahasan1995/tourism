package filter

import (
	"fmt"
	"regexp"

	"github.com/goccy/go-json"

	"go.mongodb.org/mongo-driver/bson"
)

// CustomerPersonaFilter defines filter criteria for customer personas
type CustomerPersonaFilter struct {
	Title      *string `json:"title"`
	IsActive   *bool   `json:"isActive"`
	Preference *string `json:"preference"` // Search across different preference fields
	ProgramId  *string `json:"programId"`

	Page *int `json:"page"`
	Size *int `json:"size"`
}

// NewCustomerPersonaFilter creates a new filter instance from a query string
func NewCustomerPersonaFilter(query string) (*CustomerPersonaFilter, error) {
	f := &CustomerPersonaFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// BuildPipeline creates a MongoDB aggregation pipeline based on the filter
func (f *CustomerPersonaFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Add trash filter by default
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A
	if f.Title != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Title)
		re, _ := regexp.Compile(pattern)
		titleFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$title",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, titleFilter)
	}

	if f.IsActive != nil {
		ands = append(ands, bson.M{"isActive": *f.IsActive})
	}

	if f.Preference != nil {
		// Search across all preference fields
		preferencePattern := fmt.Sprintf(".*%s.*", *f.Preference)
		re, _ := regexp.Compile(preferencePattern)
		preferenceFilter := bson.M{"$or": []bson.M{
			{"explorationPreferences": bson.M{"$elemMatch": bson.M{"$regex": re.String(), "$options": "i"}}},
			{"relaxationPreferences": bson.M{"$elemMatch": bson.M{"$regex": re.String(), "$options": "i"}}},
			{"travelMustHaves": bson.M{"$elemMatch": bson.M{"$regex": re.String(), "$options": "i"}}},
		}}
		ands = append(ands, preferenceFilter)
	}

	if f.ProgramId != nil {
		ands = append(ands, bson.M{"programId": *f.ProgramId})
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
		{"$sort": bson.M{"displayOrder": 1}},
	}

	return pipeline
}

// ToBsonFilter converts the filter struct to a MongoDB filter (kept for backward compatibility)
func (f *CustomerPersonaFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	if f.Title != nil {
		filterConditions = append(filterConditions, bson.M{"title": bson.M{"$regex": *f.Title, "$options": "i"}})
	}

	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	if f.Preference != nil {
		// Search across all preference fields
		preferenceFilter := bson.M{"$or": []bson.M{
			{"explorationPreferences": bson.M{"$elemMatch": bson.M{"$regex": *f.Preference, "$options": "i"}}},
			{"relaxationPreferences": bson.M{"$elemMatch": bson.M{"$regex": *f.Preference, "$options": "i"}}},
			{"travelMustHaves": bson.M{"$elemMatch": bson.M{"$regex": *f.Preference, "$options": "i"}}},
		}}
		filterConditions = append(filterConditions, preferenceFilter)
	}

	if f.ProgramId != nil {
		filterConditions = append(filterConditions, bson.M{"programId": *f.ProgramId})
	}

	return bson.M{"$and": filterConditions}
}
