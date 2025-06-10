package filter

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelReqFilters struct {
	CustomerName *string             `json:"customerName"`
	Status       *string             `json:"status"`
	CustomerId   *primitive.ObjectID `json:"customerId"`
}

// func NewTravelReqFilters(query string) (*TravelReqFilters, error) {
// 	f := &TravelReqFilters{}

// 	if query != "" {
// 		if err := json.Unmarshal([]byte(query), &f); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return f, nil
// }

func (f TravelReqFilters) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A
	if f.CustomerName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CustomerName)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$customerName",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, nameFilter)
	}

	if f.Status != nil {
		m["status"] = *f.Status
	}

	if f.CustomerId != nil {
		m["customerId"] = *f.CustomerId
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
