package filters

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

type AgentJoinFilter struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

func (f AgentJoinFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A
	if f.Name != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Name)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$fullName",
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

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
