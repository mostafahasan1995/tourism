package filter

import (
	"encoding/json"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsFilter struct {
	Status   *string `json:"status"`
	Type     *string `json:"type"`
	Ref      *string `json:"ref"`
	UserId   *string `json:"userId"`
	Username *string `json:"username"`
	Customer *string `json:"customer"`
	Search   *string `json:"search"`
}

func NewReviewsFilter(query string) (*ReviewsFilter, error) {
	f := &ReviewsFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (f *ReviewsFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A

	ands = append(ands, bson.M{"trash": bson.M{"$ne": true}})

	if f.Status != nil {
		ands = append(ands, bson.M{"status": *f.Status})
	}

	if f.Type != nil {
		ands = append(ands, bson.M{"type": *f.Type})
	}

	if f.Ref != nil {
		if refId, err := primitive.ObjectIDFromHex(*f.Ref); err == nil {
			ands = append(ands, bson.M{"ref": refId})
		}
	}

	if f.UserId != nil {
		ands = append(ands, bson.M{"userId": *f.UserId})
	}

	if f.Username != nil {
		ands = append(ands, bson.M{"username": bson.M{"$regex": *f.Username, "$options": "i"}})
	}

	if f.Customer != nil {
		ands = append(ands, bson.M{"customer": *f.Customer})
	}

	if f.Search != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Search)
		re, _ := regexp.Compile(pattern)
		searchFilter := bson.M{
			"$or": []bson.M{
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$description",
							"regex":   re.String(),
							"options": "i",
						},
					},
				},
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$adviceForTravelers",
							"regex":   re.String(),
							"options": "i",
						},
					},
				},
			},
		}
		ands = append(ands, searchFilter)
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
