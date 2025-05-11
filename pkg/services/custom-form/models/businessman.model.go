package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Businessman struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
}
