package travelreq

import (
	"context"
	"encoding/json"

	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
)

type ReqSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) // todo: return uniform response
}

type ReqType struct {
	Svcs     ReqSvcs
	CollName string
}

type ReqTypes map[string]ReqType

func NewReqTypes(i *do.Injector) (ReqTypes, error) {
	return map[string]ReqType{
		"vipcar": {
			Svcs:     do.MustInvoke[VipCarRequestSvcs](i),
			CollName: "tourismVipCarRequest",
		},
	}, nil
}
