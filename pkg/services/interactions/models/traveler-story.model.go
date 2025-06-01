package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelerStoryDto struct {
	Name         string               `bson:"name" json:"name"`
	Destinations []primitive.ObjectID `bson:"destinations" json:"destinations"`
	DiaryTitle   string               `bson:"diaryTitle" json:"diaryTitle"`
	TripType     string               `bson:"tripType" json:"tripType"`
	Bio          string               `bson:"bio" json:"bio"`
	CoverImage   []types.FileField    `bson:"coverImage" json:"coverImage"`
	Status       string               `bson:"status" json:"status"`
	ZoneName     string               `bson:"zoneName" json:"zoneName"`
	Description  string               `bson:"description" json:"description"`
	TravelImages []types.FileField    `bson:"travelImages" json:"travelImages"`
	TravelVideos []types.FileField    `bson:"travelVideos" json:"travelVideos"`
	VideoUrl     string               `bson:"videoUrl" json:"videoUrl"`
	Activities   []primitive.ObjectID `bson:"activities" json:"activities"`
	Tips         string               `bson:"tips" json:"tips"`
	TipsImage    []types.FileField    `bson:"tipsImage" json:"tipsImage"`
}

type TravelerStory struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TravelerStoryDto `bson:",inline"`
	RejectReason     string             `bson:"rejectReason" json:"rejectReason"`
	Feedback         string             `bson:"feedback" json:"feedback"`
	Trash            bool               `bson:"trash" json:"trash"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy        primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy        primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}

type TravelerStoryWithPagination struct {
	TravelerStories []TravelerStory  `bson:"travelerStories" json:"travelerStories"`
	Pagination      types.Pagination `bson:"pagination" json:"pagination"`
}

//

type TravelerStoryStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"` //pending - approved - rejected
	Reason string `bson:"reason" json:"reason" validate:"required_if=Status rejected"`
}

type TravelerStoryFeedback struct {
	Feedback string `bson:"feedback" json:"feedback"`
}
