package helpers

import (
	"net/http"
	"strconv"

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

// ExtractPaginationParams extracts page and size parameters from HTTP request query string
func ExtractPaginationParams(r *http.Request) (page, size int) {
	// Extract page parameter (default to 1)
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			page = 1
		}
	} else {
		page = 1
	}

	// Extract size parameter (default to 10)
	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		} else {
			size = 10
		}
	} else {
		size = 10
	}

	return page, size
}
