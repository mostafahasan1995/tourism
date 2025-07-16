package models

import (
	picklistmodels "larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelerStoryDto struct {
	Name         transl.Localizable[string] `bson:"name" json:"name"`
	Destinations []primitive.ObjectID       `bson:"destinations" json:"destinations"`
	DiaryTitle   transl.Localizable[string] `bson:"diaryTitle" json:"diaryTitle"`
	TripType     string                     `bson:"tripType" json:"tripType"`
	Bio          transl.Localizable[string] `bson:"bio" json:"bio"`
	CoverImage   types.FileField            `bson:"coverImage" json:"coverImage"`
	Status       string                     `bson:"status" json:"status"`
	ZoneName     string                     `bson:"zoneName" json:"zoneName"`
	Description  transl.Localizable[string] `bson:"description" json:"description"`
	TravelImages []types.FileField          `bson:"travelImages" json:"travelImages"`
	TravelVideos []types.FileField          `bson:"travelVideos" json:"travelVideos"`
	VideoUrl     string                     `bson:"videoUrl" json:"videoUrl"`
	Activities   []primitive.ObjectID       `bson:"activities" json:"activities"`
	Tips         transl.Localizable[string] `bson:"tips" json:"tips"`
	TipsImage    []types.FileField          `bson:"tipsImage" json:"tipsImage"`
}

type TravelerStory struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TravelerStoryDto `bson:",inline"`
	RejectReason     string             `bson:"rejectReason" json:"rejectReason"`
	Feedback         string             `bson:"feedback" json:"feedback"`
	Trash            bool               `bson:"trash" json:"trash"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedBy        primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt        time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	UpdatedBy        primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

type TravelerStoryWithPagination struct {
	TravelerStories []TravelerStoryRes `bson:"travelerStories" json:"travelerStories"`
	Pagination      types.Pagination   `bson:"pagination" json:"pagination"`
}

type TravelerStoryRes struct {
	TravelerStory    `bson:",inline"`
	IsFav            bool `bson:"isFav" json:"isFav"`
	DestinationsData []picklistmodels.Destination
	ActivitiesData   []picklistmodels.Activities
}

//

type TravelerStoryStatusDto struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=pending approved rejected"` //pending - approved - rejected
	Reason string `bson:"reason" json:"reason" validate:"required_if=Status rejected"`
}

type TravelerStoryFeedback struct {
	Feedback string `bson:"feedback" json:"feedback"`
}
