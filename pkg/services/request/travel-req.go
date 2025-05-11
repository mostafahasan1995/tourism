package request

import (
	"context"
	"encoding/json"
	"errors"
	"larsa-tourism-microservices/pkg/services/request/models"
	"larsa-tourism-microservices/pkg/services/request/repo"

	"github.com/samber/do"
)

type TravelReqSvcs interface {
	Get(ctx context.Context, skip int64, limit int64) (*models.TravelReqWithPagination, error)
	Add(ctx context.Context, serviceType string, data json.RawMessage) (any, error)
}

type travelreqsvcs struct {
	repo       repo.TravelReqRepo
	vipCarRepo repo.VipCarRequestRepo
}

func NewTravelReqSvcs(i *do.Injector) (TravelReqSvcs, error) {
	return &travelreqsvcs{
		repo:       do.MustInvoke[repo.TravelReqRepo](i),
		vipCarRepo: do.MustInvoke[repo.VipCarRequestRepo](i),
	}, nil
}

func (t *travelreqsvcs) Get(ctx context.Context, skip int64, limit int64) (*models.TravelReqWithPagination, error) {
	return nil, nil
}

func (t *travelreqsvcs) Add(ctx context.Context, serviceType string, data json.RawMessage) (any, error) {

	switch serviceType {
	case "vipcar":
		{
			var vipCar models.VipCarRequest
			err := json.Unmarshal(data, &vipCar)
			if err != nil {
				return nil, err
			}
			if err := t.vipCarRepo.Add(ctx, &vipCar); err != nil {
				return nil, err
			}

			// travelReq := models.TravelReq{
			// 	Id:          primitive.NewObjectID(),
			// 	ReqId:       "req-00001",
			// 	ServiceType: serviceType,
			// 	//Package:     vipCar.Id,
			// 	//Program:     vipCar.Id,
			// 	Customer:    vipCar.Id,
			// 	Status:      "pending",
			// 	Ref:         vipCar.Id,
			// }

			return &vipCar, nil
		}
	}

	return nil, errors.New("service type not found")
}
