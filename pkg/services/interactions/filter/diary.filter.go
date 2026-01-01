package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiaryFilter struct {
	IDs []primitive.ObjectID `json:"ids"`
}

func (f DiaryFilter) BuildPipeline(m bson.M) []bson.M {
	if m == nil {
		m = bson.M{}
	}

	// Default filter for non-trashed items
	m["trash"] = bson.M{"$ne": true}

	// Filter by IDs if provided
	if len(f.IDs) > 0 {
		m["_id"] = bson.M{"$in": f.IDs}
	}

	return []bson.M{
		{"$match": m},
	}
}

// ToBsonFilter converts the filter to a MongoDB filter
func (f DiaryFilter) ToBsonFilter() bson.M {
	filter := bson.M{"trash": bson.M{"$ne": true}}

	// Filter by IDs if provided
	if len(f.IDs) > 0 {
		filter["_id"] = bson.M{"$in": f.IDs}
	}

	return filter
}
