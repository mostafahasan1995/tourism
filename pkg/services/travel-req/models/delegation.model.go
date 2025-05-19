package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delegation struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DelegationType      string             `bson:"delegationType" json:"delegationType"`
	OrganizationName    string             `bson:"organizationName" json:"organizationName"`
	ClientName          string             `bson:"clientName" json:"clientName"`
	ClientPhone         string             `bson:"clientPhone" json:"clientPhone"`
	ClientEmail         string             `bson:"clienteEmail" json:"clientEmail"`
	Nationality         string             `bson:"nationality" json:"nationality"`
	TripDuration        int                `bson:"tripDuration" json:"tripDuration"`
	TripCoordinatorName string             `bson:"tripCoordinatorName" json:"tripCoordinatorName"`
	Destinations        []CFDestination    `bson:"destinations" json:"destinations"`
	TripCoordinator     primitive.ObjectID `bson:"tripCoordinator" json:"tripCoordinator"`
	ContactMethod       []string           `bson:"contactMethod" json:"contactMethod"`
	SpecialReq          string             `bson:"specialReq" json:"specialReq"`
}
