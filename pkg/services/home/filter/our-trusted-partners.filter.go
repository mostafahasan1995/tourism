package filter

import (
	"encoding/json"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

type TrustedPartnersFilter struct {
	Title    *string `json:"title"`
	IsActive *bool   `json:"isActive"`

	Page *int `json:"page"`
	Size *int `json:"size"`
}

func NewTrustedPartnersFilter(query string) (*TrustedPartnersFilter, error) {
	f := &TrustedPartnersFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *TrustedPartnersFilter) BuildPipeline(m bson.M) []bson.M {
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

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}

// Keep the old method for backward compatibility
func (f *TrustedPartnersFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	if f.Title != nil {
		filterConditions = append(filterConditions, bson.M{"title": bson.M{"$regex": *f.Title, "$options": "i"}})
	}

	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	return bson.M{"$and": filterConditions}
}
