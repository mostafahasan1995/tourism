package travelreq

import (
	"context"
	"encoding/json"

	"larsa-tourism-microservices/pkg/services/travel-req/enums"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
)

type ReqSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) // todo: return uniform response
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type ReqType struct {
	Svcs     ReqSvcs
	CollName string
}

type ReqTypes map[enums.ServiceType]ReqType

func NewReqTypes(i *do.Injector) (ReqTypes, error) {
	return map[enums.ServiceType]ReqType{
		enums.ServiceTypeBusinessMan: {
			Svcs:     do.MustInvoke[BusinessManSvcs](i),
			CollName: "tourismBusinessmenRequests",
		},
		enums.ServiceTypeCustomPlan: {
			Svcs:     do.MustInvoke[CustomPlanSvcs](i),
			CollName: "tourismCustomPlanRequests",
		},
		enums.ServiceTypeDelegation: {
			Svcs:     do.MustInvoke[DelegationSvcs](i),
			CollName: "tourismDelegationRequests",
		},
		enums.ServiceTypeFlightTicket: {
			Svcs:     do.MustInvoke[FlightTicketRequestSvcs](i),
			CollName: "tourismFlightTicketRequests",
		},
		enums.ServiceTypePartner: {
			Svcs:     do.MustInvoke[PartnerRequestSvcs](i),
			CollName: "tourismPartnerRequests",
		},
		enums.ServiceTypeVipCar: {
			Svcs:     do.MustInvoke[VipCarRequestSvcs](i),
			CollName: "tourismVipCarRequest",
		},
	}, nil
}
