package filters

import (
	"encoding/json"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

type AgentFilter struct {
	Name *string `json:"name"`
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

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
