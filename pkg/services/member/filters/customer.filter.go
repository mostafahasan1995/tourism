package filters

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CustomerFilter struct {
	CustomerName *string    `json:"customerName"`
	Date         *time.Time `json:"date"`
	Status       *string    `json:"status"`
	Email        *string    `json:"email"`
}

func NewCustomerFilter(query string) (*CustomerFilter, error) {
	filter := &CustomerFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), filter); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (f CustomerFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	var ands bson.A
	if f.CustomerName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CustomerName)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$name",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, nameFilter)
	}

	if f.Date != nil {
		// Filter by date (exact match on the date part)
		startOfDay := time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 0, 0, 0, 0, f.Date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		dateFilter := bson.M{
			"createdAt": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
		}
		ands = append(ands, dateFilter)
	}

	if f.Status != nil {
		statusFilter := bson.M{
			"status": *f.Status,
		}
		ands = append(ands, statusFilter)
	}

	if f.Email != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Email)
		re, _ := regexp.Compile(pattern)
		emailFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$clientContact.email",
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
	}

	return pipeline
}
