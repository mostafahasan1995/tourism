package filter

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivitiesFilter struct {
	SearchWord    *string              `json:"searchWord"`
	ActivitiesIds []primitive.ObjectID `json:"activitiesIds"`
	Page          *int64               `json:"page"`
	Size          *int64               `json:"size"`
}

func (f ActivitiesFilter) BuildPipeline(m bson.M) []bson.M {

	var ands []bson.M
	if f.SearchWord != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.SearchWord)
		re, _ := regexp.Compile(pattern)
		wordFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$name",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, wordFilter)
	}

	if len(f.ActivitiesIds) > 0 {
		m["_id"] = bson.M{"$in": f.ActivitiesIds}
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	return []bson.M{
		{
			"$match": m,
		},
	}
}
