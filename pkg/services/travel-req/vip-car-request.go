package travelreq

import (
	"context"
	"encoding/json"
	"larsa-tourism-microservices/pkg/services/travel-req/models"
	"larsa-tourism-microservices/pkg/services/travel-req/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VipCarRequestSvcs interface {
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
}

type vipCarRequestsvcs struct {
	repo repo.VipCarRequestRepo
}

func NewVipCarRequestSvcs(i *do.Injector) (VipCarRequestSvcs, error) {
	return &vipCarRequestsvcs{
		repo: do.MustInvoke[repo.VipCarRequestRepo](i),
	}, nil
}

func (v *vipCarRequestsvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.VipCarRequestDto
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	vipCarReq := &models.VipCarRequest{
		Id:               primitive.NewObjectID(),
		VipCarRequestDto: req,
	}

	if err := v.repo.Add(ctx, vipCarReq); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id: vipCarReq.Id,
	}, nil
}

func (v *vipCarRequestsvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return v.repo.GetByFilter(ctx, filter)
}
