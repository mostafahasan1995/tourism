package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BusinessManDto struct {
	Purpose             string `bson:"purpose" json:"purpose"` //meeting - investment - conference - other
	ClientName          string `bson:"clientName" json:"clientName"`
	ClientPhone         string `bson:"clientPhone" json:"clientPhone"`
	ClientEmail         string `bson:"clientEmail" json:"clientEmail"`
	ClientIsCoordinator bool   `bson:"clientIsCoordinator" json:"clientIsCoordinator"`
	CoordinatorName     string `bson:"coordinatorName" json:"coordinatorName"`
	CoordinatorPhone    string `bson:"coordinatorPhone" json:"coordinatorPhone"`
	CoordinatorEmail    string `bson:"coordinatorEmail" json:"coordinatorEmail"`
	Nationality         string `bson:"nationality" json:"nationality"`
	TripDuration        int    `bson:"tripDuration" json:"tripDuration"`

	Destinations    []CFDestination    `bson:"destinations" json:"destinations"`
	TripCoordinator primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
	ContactMethod   []string           `bson:"contactMethod" json:"contactMethod"`
	SpecialReq      string             `bson:"specialReq" json:"specialReq"`
}
type BusinessMan struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	BusinessManDto `bson:",inline"`
}
