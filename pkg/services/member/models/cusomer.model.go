package models

import (
	"larsa-tourism-microservices/pkg/types"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerDto struct {
	Name                string            `bson:"name" json:"name"`
	Nationality         string            `bson:"nationality" json:"nationality"`
	Company             string            `bson:"company" json:"company"`
	TripCoordinatorName string            `bson:"tripCoordinatorName" json:"tripCoordinatorName"`
	About               string            `bson:"about" json:"about"`
	Image               []types.FileField `bson:"image" json:"image"`
	ClientContact       MemberContact     `bson:"clientContact" json:"clientContact"`
	CoordinatorContact  MemberContact     `bson:"coordinatorContact" json:"coordinatorContact"`
	SocialMedia         []SocialMedia     `bson:"socialMedia" json:"socialMedia"`
	Security            MemberSecurity    `bson:"security" json:"security"`
}

type SocialMedia struct {
	Platform string `bson:"platform" json:"platform"` //e.g. facebook, instagram, twitter, linkedin, etc.
	Link     string `bson:"link" json:"link"`
}

type Customer struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` //same as user id
	CustomerId  string             `bson:"customerId,omitempty" json:"customerId,omitempty"`
	CustomerDto `bson:",inline"`
	Status      string             `bson:"status" json:"status"`
	Trash       bool               `bson:"trash" json:"trash"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy   primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

type CustomerWithPagination struct {
	Customers  []Customer       `bson:"customers" json:"customers"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

//

type CustomerRegisterData struct {
	ClientName  string `bson:"clientName" json:"clientName"`
	ClientPhone string `bson:"clientPhone" json:"clientPhone"`
	ClientEmail string `bson:"clientEmail" json:"clientEmail"`
	Nationality string `bson:"nationality" json:"nationality"`
	Password    string `bson:"password" json:"password"`
}
