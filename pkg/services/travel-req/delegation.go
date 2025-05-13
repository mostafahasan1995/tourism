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

type DelegationSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type delegationSvcs struct {
	repo repo.DelegationRepo
}

func NewDelegationSvcs(i *do.Injector) (DelegationSvcs, error) {
	return &delegationSvcs{
		repo: do.MustInvoke[repo.DelegationRepo](i),
	}, nil
}

func (d *delegationSvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return d.repo.GetByFilter(ctx, filter)
}

func (d *delegationSvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.Delegation
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	req.Id = primitive.NewObjectID()

	if err := d.repo.Add(ctx, &req); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id:            req.Id,
		CustomerName:  req.OrganizationName,
		CustomerPhone: req.Phone,
		CustomerEmail: req.Email,
	}, nil
}

func (d *delegationSvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.Delegation
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	req.Id = _id

	_, err = d.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": req})
	return err
}
