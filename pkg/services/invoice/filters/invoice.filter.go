package filters

import (
	"encoding/json"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

type InvoiceFilter struct {
	CustomerName *string `json:"customerName"`
}

func NewInvoiceFilter(query string) (*InvoiceFilter, error) {
	f := &InvoiceFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), f); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *InvoiceFilter) BuildPipeline(m bson.M) []bson.M {

	var ands bson.A
	if f.CustomerName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CustomerName)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$customer.name",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, nameFilter)
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline

}
