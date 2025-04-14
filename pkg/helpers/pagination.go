package helpers

import (
	"go.mongodb.org/mongo-driver/bson"
)

func BuildPaginationPipeline(pipeline []bson.M, name string, skip, limit int64, sortField string, sortOrder int) []bson.M {
	return append(pipeline, bson.M{
		"$facet": bson.M{
			"pagination": []bson.M{
				{"$count": "totalCount"},
				{"$addFields": bson.M{
					"perPage": limit,
					"totalPages": bson.M{
						"$ceil": bson.M{
							"$divide": []interface{}{"$totalCount", limit},
						},
					},
				}},
			},
			name: []bson.M{
				{"$sort": bson.M{sortField: sortOrder}},
				{"$skip": skip},
				{"$limit": limit},
			},
		},
	}, bson.M{
		"$unwind": bson.M{
			"path":                       "$pagination",
			"preserveNullAndEmptyArrays": true,
		},
	})

}
