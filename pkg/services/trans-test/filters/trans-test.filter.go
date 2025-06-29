package filters

import "go.mongodb.org/mongo-driver/bson"

type TransTestFilters struct{}

func (f TransTestFilters) BuildPipeline(m bson.M) []bson.M {

	return []bson.M{
		{"$match": m},
	}
}
