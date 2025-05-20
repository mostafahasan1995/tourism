package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CustomPlanDto struct {
	TripType        string             `bson:"tripType" json:"tripType"` //e.g. family - honeymoon - luxury retreat - other
	ClientName      string             `bson:"clientName" json:"clientName"`
	Phone           string             `bson:"phone" json:"phone"`
	Email           string             `bson:"email" json:"email"`
	Nationality     string             `bson:"nationality" json:"nationality"`
	TripDuration    int                `bson:"tripDuration" json:"tripDuration"`
	Destinations    []CFDestination    `bson:"destinations" json:"destinations"`
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}
type CustomPlan struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CustomPlanDto `bson:",inline"`
}
