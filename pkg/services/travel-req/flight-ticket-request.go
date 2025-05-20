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

type FlightTicketRequestSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (any, error)
	Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error)
	Update(ctx context.Context, id string, data json.RawMessage) error
}

type flightTicketRequestSvcs struct {
	repo repo.FlightTicketRequestRepo
}

func NewFlightTicketRequestSvcs(i *do.Injector) (FlightTicketRequestSvcs, error) {
	return &flightTicketRequestSvcs{
		repo: do.MustInvoke[repo.FlightTicketRequestRepo](i),
	}, nil
}

func (f *flightTicketRequestSvcs) GetByFilter(ctx context.Context, filter bson.M) (any, error) {
	return f.repo.GetByFilter(ctx, filter)
}

func (f *flightTicketRequestSvcs) Add(ctx context.Context, data json.RawMessage) (*models.ReqAddData, error) {
	var req models.FlightTicketRequestDto
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	flightTicketReq := &models.FlightTicketRequest{
		Id:                     primitive.NewObjectID(),
		FlightTicketRequestDto: req,
	}

	if err := f.repo.Add(ctx, flightTicketReq); err != nil {
		return nil, err
	}

	return &models.ReqAddData{
		Id: flightTicketReq.Id,
	}, nil
}

func (f *flightTicketRequestSvcs) Update(ctx context.Context, id string, data json.RawMessage) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var req models.FlightTicketRequestDto
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	flightTicketReq := &models.FlightTicketRequest{
		Id:                     _id,
		FlightTicketRequestDto: req,
	}

	_, err = f.repo.Patch(ctx, bson.M{"_id": _id}, bson.M{"$set": flightTicketReq})
	return err
}
