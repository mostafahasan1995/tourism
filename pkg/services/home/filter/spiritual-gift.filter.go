package filter

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiritualGiftFilter struct {
	Enabled   *bool               `json:"enabled,omitempty"`
	PackageId *primitive.ObjectID `json:"packageId,omitempty"`
	ProgramId *string             `json:"programId,omitempty"`
}

func (f SpiritualGiftFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A

	if f.Enabled != nil {
		enabledFilter := bson.M{"enabled": *f.Enabled}
		ands = append(ands, enabledFilter)
	}

	if f.PackageId != nil && !f.PackageId.IsZero() {
		packageFilter := bson.M{"packageId": *f.PackageId}
		ands = append(ands, packageFilter)
	}

	if f.ProgramId != nil && *f.ProgramId != "" {
		pattern := fmt.Sprintf(".*%s.*", *f.ProgramId)
		re, _ := regexp.Compile(pattern)
		programFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$programId",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, programFilter)
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
