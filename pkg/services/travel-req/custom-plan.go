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

type CustomPlanSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type customPlanSvcs struct {
	repo repo.CustomPlanRepo
}

func NewCustomPlanSvcs(i *do.Injector) (CustomPlanSvcs, error) {
	return &customPlanSvcs{
		repo: do.MustInvoke[repo.CustomPlanRepo](i),
	}, nil
}

func (c *customPlanSvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return c.repo.GetByFilter(ctx, filter)
}

func (c *customPlanSvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.CustomPlan
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	req.Id = primitive.NewObjectID()

	if err := c.repo.Add(ctx, &req); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id:            req.Id,
		CustomerName:  req.ClientName,
		CustomerPhone: req.Phone,
		CustomerEmail: req.Email,
		Nationality:   req.Nationality,
	}, nil
}

func (c *customPlanSvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.CustomPlan
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	req.Id = _id

	_, err = c.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": req})
	return err
}
