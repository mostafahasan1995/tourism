package models

import (
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelReq struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ReqId       string             `bson:"reqId" json:"reqId"`
	Package     primitive.ObjectID `bson:"package" json:"package"`
	Program     primitive.ObjectID `bson:"program" json:"program"`
	ServiceType string             `bson:"serviceType" json:"serviceType"`
	Date        time.Time          `bson:"date" json:"date"`
	Customer    primitive.ObjectID `bson:"customer" json:"customer"`
	Status      string             `bson:"status" json:"status"`
	Ref         primitive.ObjectID `bson:"ref" json:"ref"`
}

type TravelReqWithPagination struct {
	TravelReqs []TravelReq       `bson:"travelReqs" json:"travelReqs"`
	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
