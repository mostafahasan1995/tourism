package filter

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaveFilter struct {
	Type   *string `json:"type"`
	UserId *string `json:"userId"`
	IsFav  *bool   `json:"isFav"`
	RefId  *string `json:"refId"`
}

func NewFaveFilter(query string) (*FaveFilter, error) {
	f := &FaveFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (f FaveFilter) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A

	if f.Type != nil {
		ands = append(ands, bson.M{"type": models.FaveType(*f.Type)})
	}

	if f.UserId != nil {
		if userId, err := primitive.ObjectIDFromHex(*f.UserId); err == nil {
			ands = append(ands, bson.M{"userId": userId})
		}
	}

	if f.IsFav != nil {
		ands = append(ands, bson.M{"isFav": *f.IsFav})
	}

	if f.RefId != nil {
		if refId, err := primitive.ObjectIDFromHex(*f.RefId); err == nil {
			ands = append(ands, bson.M{"refId": refId})
		}
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
