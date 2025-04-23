package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaveDto struct {
	Type  string             `bson:"type" json:"type"`
	Ref   primitive.ObjectID `bson:"ref" json:"ref"`
	IsFav bool               `bson:"isFav" json:"isFav"`
}

type Fave struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	FaveDto   `bson:",inline"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
