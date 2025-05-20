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
	var req models.CustomPlanDto
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	customPlanReq := &models.CustomPlan{
		Id:            primitive.NewObjectID(),
		CustomPlanDto: req,
	}

	if err := c.repo.Add(ctx, customPlanReq); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id:            customPlanReq.Id,
		CustomerName:  customPlanReq.ClientName,
		CustomerPhone: customPlanReq.Phone,
		CustomerEmail: customPlanReq.Email,
		Nationality:   customPlanReq.Nationality,
	}, nil
}

func (c *customPlanSvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.CustomPlanDto
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	customPlanReq := &models.CustomPlan{
		Id:            _id,
		CustomPlanDto: req,
	}

	_, err = c.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": customPlanReq})
	return err
}
