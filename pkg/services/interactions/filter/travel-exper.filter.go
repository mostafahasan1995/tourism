package filter

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelerStoryFilter struct {
	Status      *string             `json:"status"`
	Username    *string             `json:"username"`
	Title       *string             `json:"title"`
	Destination *primitive.ObjectID `json:"destination"`
	Activity    *primitive.ObjectID `json:"activity"`
	Search      *string             `json:"search"`
}

func (f TravelerStoryFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	if f.Status != nil {
		ands = append(ands, bson.M{"status": *f.Status})
	}

	if f.Username != nil {
		ands = append(ands, bson.M{"username": bson.M{"$regex": *f.Username, "$options": "i"}})
	}

	if f.Title != nil {
		ands = append(ands, bson.M{"title": bson.M{"$regex": *f.Title, "$options": "i"}})
	}

	if f.Destination != nil {
		ands = append(ands, bson.M{"destinations": f.Destination})
	}

	if f.Activity != nil {
		ands = append(ands, bson.M{"activities": f.Activity})
	}

	if f.Search != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Search)
		re, _ := regexp.Compile(pattern)
		searchFilter := bson.M{
			"$or": []bson.M{
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$title",
							"regex":   re.String(),
							"options": "i",
						},
					},
				},
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$description",
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

type ClientStoryFilter struct {
	Title  *string `json:"title"`
	Search *string `json:"search"`
}

func (f ClientStoryFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	if f.Title != nil {
		ands = append(ands, bson.M{"title": bson.M{"$regex": *f.Title, "$options": "i"}})
	}

	if f.Search != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Search)
		re, _ := regexp.Compile(pattern)
		searchFilter := bson.M{
			"$or": []bson.M{
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$title",
							"regex":   re.String(),
							"options": "i",
						},
					},
				},
				{
					"$expr": bson.M{
						"$regexMatch": bson.M{
							"input":   "$description",
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
