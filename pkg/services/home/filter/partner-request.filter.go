package filter

import (
	"encoding/json"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

// PartnerRequestFilter defines filter criteria for partner requests
type PartnerRequestFilter struct {
	CompanyName     *string `json:"companyName"`
	BusinessType    *string `json:"businessType"`
	CompanyLocation *string `json:"companyLocation"`
	Status          *string `json:"status"`
	Email           *string `json:"email"`

	Page *int `json:"page"`
	Size *int `json:"size"`
}

// NewPartnerRequestFilter creates a new filter instance from a query string
func NewPartnerRequestFilter(query string) (*PartnerRequestFilter, error) {
	f := &PartnerRequestFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// BuildPipeline creates a MongoDB aggregation pipeline based on the filter
func (f *PartnerRequestFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Add trash filter by default
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A

	if f.CompanyName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CompanyName)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$companyName",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, nameFilter)
	}

	if f.BusinessType != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.BusinessType)
		re, _ := regexp.Compile(pattern)
		typeFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$businessType",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, typeFilter)
	}

	if f.CompanyLocation != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CompanyLocation)
		re, _ := regexp.Compile(pattern)
		locationFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$companyLocation",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, locationFilter)
	}

	if f.Status != nil {
		ands = append(ands, bson.M{"status": *f.Status})
	}

	if f.Email != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Email)
		re, _ := regexp.Compile(pattern)
		emailFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$email",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, emailFilter)
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
		{"$sort": bson.M{"createdAt": -1}}, // Sort by creation date, newest first
	}

	return pipeline
}
