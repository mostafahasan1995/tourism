package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerDto struct {
	Name                transl.Localizable[string] `bson:"name" json:"name"`
	Nationality         string                     `bson:"nationality" json:"nationality"`
	Company             transl.Localizable[string] `bson:"company" json:"company"`
	TripCoordinatorName string                     `bson:"tripCoordinatorName" json:"tripCoordinatorName"`
	About               transl.Localizable[string] `bson:"about" json:"about"`
	Image               []types.FileField          `bson:"image" json:"image"`
	ClientContact       MemberContact              `bson:"clientContact" json:"clientContact"`
	CoordinatorContact  MemberContact              `bson:"coordinatorContact" json:"coordinatorContact"`
	SocialMedia         []SocialMedia              `bson:"socialMedia" json:"socialMedia"`
	Security            MemberSecurity             `bson:"security" json:"security"`
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
type CustomerWithNewsletter struct {
	Customer               `bson:",inline" json:",inline"`
	SubscribedToNewsletter bool `bson:"subscribedToNewsletter" json:"subscribedToNewsletter"`
}

type CustomerWithNewsletterAndPagination struct {
	Customers  []CustomerWithNewsletter `bson:"customers" json:"customers"`
	Pagination types.Pagination         `bson:"pagination" json:"pagination"`
}
type CustomerWithPagination struct {
	Customers  []Customer       `bson:"customers" json:"customers"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

//

type CustomerRegisterData struct {
	ClientName  transl.Localizable[string] `bson:"clientName" json:"clientName"`
	ClientPhone types.PhoneNumber          `bson:"clientPhone" json:"clientPhone"`
	ClientEmail string                     `bson:"clientEmail" json:"clientEmail"`
	Nationality string                     `bson:"nationality" json:"nationality"`
	Password    string                     `bson:"password" json:"password"`
}

func BuildNewsletterPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from": "tourismNewsletters",
			"let":  bson.M{"email": "$security.email"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$and": []bson.M{
							{"$eq": []any{"$email", "$$email"}},
							{"$eq": []any{"$status", "active"}},
						},
					},
				}},
			},
			"as": "newsletter",
		}},
		{"$addFields": bson.M{
			"subscribedToNewsletter": bson.M{"$gt": []any{bson.M{"$size": "$newsletter"}, 0}},
		}},
		{"$project": bson.M{"newsletter": 0}},
	}
}
