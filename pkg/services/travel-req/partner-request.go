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

type PartnerRequestSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type partnerRequestSvcs struct {
	repo repo.PartnerRequestRepo
}

func NewPartnerRequestSvcs(i *do.Injector) (PartnerRequestSvcs, error) {
	return &partnerRequestSvcs{
		repo: do.MustInvoke[repo.PartnerRequestRepo](i),
	}, nil
}

func (p *partnerRequestSvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return p.repo.GetByFilter(ctx, filter)
}

func (p *partnerRequestSvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.PartnerRequestDto
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	partnerReq := &models.PartnerRequest{
		Id:                primitive.NewObjectID(),
		PartnerRequestDto: req,
	}

	if err := p.repo.Add(ctx, partnerReq); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id:            partnerReq.Id,
		CustomerName:  partnerReq.ContactDetail.FullName,
		CustomerPhone: partnerReq.ContactDetail.PhoneNumber,
		CustomerEmail: partnerReq.ContactDetail.Email,
	}, nil
}

func (p *partnerRequestSvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.PartnerRequestDto
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	partnerReq := &models.PartnerRequest{
		Id:                _id,
		PartnerRequestDto: req,
	}

	_, err = p.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": partnerReq})
	return err
}
