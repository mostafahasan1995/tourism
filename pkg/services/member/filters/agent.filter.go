package filters

import (
	"fmt"
	"regexp"

	"github.com/goccy/go-json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentFilter struct {
	Name           *string              `json:"name"`
	Status         *string              `json:"status"`
	DestinationIds []primitive.ObjectID `json:"destinationIds"`
}

func NewAgentFilter(query string) (*AgentFilter, error) {
	filter := &AgentFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), filter); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

var destinationLookup = []bson.M{
	{
		"$lookup": bson.M{
			"from": "tourismDestinations",
			"let":  bson.M{"countries": "$countries"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$in": bson.A{"$name", "$$countries"}},
							},
						},
					},
				},
			},
			"as": "destinations",
		},
	},
}

func (f *AgentFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A
	if f.Name != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Name)
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
