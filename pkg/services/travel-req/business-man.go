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

type BusinessManSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type businessmansvcs struct {
	repo repo.BusinessManRepo
}

func NewBusinessManSvcs(i *do.Injector) (BusinessManSvcs, error) {
	return &businessmansvcs{
		repo: do.MustInvoke[repo.BusinessManRepo](i),
	}, nil
}

func (b *businessmansvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return b.repo.GetByFilter(ctx, filter)
}

func (b *businessmansvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.BusinessMan
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	req.Id = primitive.NewObjectID()

	if err := b.repo.Add(ctx, &req); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id:            req.Id,
		CustomerName:  req.ClientName,
		CustomerPhone: req.ClientPhone,
		CustomerEmail: req.ClientEmail,
		Nationality:   req.Nationality,
	}, nil
}

func (b *businessmansvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.BusinessMan
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	req.Id = _id

	_, err = b.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": req})
	return err
}
