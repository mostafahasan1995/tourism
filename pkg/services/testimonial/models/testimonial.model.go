package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestimonialDto struct {
	Message string            `bson:"message" json:"message"`
	Album   []types.FileField `bson:"album" json:"album"`
}

type Testimonial struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId         primitive.ObjectID `bson:"userId" json:"userId"`
	TestimonialDto `bson:",inline"`
	CreatedAt      time.Time `bson:"createdAt" json:"createdAt"`
}
