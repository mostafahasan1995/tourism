package filter

import (
	"go.mongodb.org/mongo-driver/bson"
)

type DestinationFilter struct {
	SearchWord *string `json:"searchWord"`
	Name       *string `json:"name"`
	// Page       int64  `json:"page"`
	// Size       int64  `json:"size"`
}

// func NewDestinationFilter(query string) (*DestinationFilter, error) {
// 	f := &DestinationFilter{}
// 	if query != "" {
// 		if err := json.Unmarshal([]byte(query), f); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return f, nil
// }

func (f DestinationFilter) BuildPipeline(m bson.M) []bson.M {

	if f.SearchWord != nil {
		m["$or"] = []bson.M{
			{
				"name": bson.M{
					"$regex":   *f.SearchWord,
					"$options": "i",
				},
			},
			{
				"country": bson.M{
					"$regex":   *f.SearchWord,
					"$options": "i",
				},
			},
		}
	}

	if f.Name != nil {
		m["name"] = *f.Name
	}

	return []bson.M{
		{
			"$match": m,
		},
	}
}
