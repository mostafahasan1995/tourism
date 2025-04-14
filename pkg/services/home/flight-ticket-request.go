package home

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlightTicketRequestSvcs interface {
	GetOne(ctx context.Context, id string) (*models.FlightTicketRequest, error)
	GetAll(ctx context.Context, filter filter.FlightTicketRequestFilter) (models.FlightTicketRequestPagination, error)
	Add(ctx context.Context, data *models.FlightTicketRequestDto) error
	AddMany(ctx context.Context, data []models.FlightTicketRequestDto) error
	Update(ctx context.Context, id string, data *models.FlightTicketRequestDto) error
	Delete(ctx context.Context, id string) error
}

type flightTicketRequestsvcs struct {
	repo repo.FlightTicketRequestRepo
}

func NewFlightTicketRequestSvcs(i *do.Injector) (FlightTicketRequestSvcs, error) {
	return &flightTicketRequestsvcs{
		repo: do.MustInvoke[repo.FlightTicketRequestRepo](i),
	}, nil
}

func (l *flightTicketRequestsvcs) GetOne(ctx context.Context, id string) (*models.FlightTicketRequest, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *flightTicketRequestsvcs) GetAll(ctx context.Context, filter filter.FlightTicketRequestFilter) (models.FlightTicketRequestPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.FlightTicketRequestPagination{}, err
	}

	return data, nil
}

func (l *flightTicketRequestsvcs) Add(ctx context.Context, data *models.FlightTicketRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	flightTicketRequest := &models.FlightTicketRequest{
		FlightTicketRequestDto: models.FlightTicketRequestDto{
			DestinationFrom:        data.DestinationFrom,
			DestinationTo:          data.DestinationTo,
			TripType:               data.TripType,
			TravelClass:            data.TravelClass,
			NumberOfAdults:         data.NumberOfAdults,
			NumberOfChildren:       data.NumberOfChildren,
			NumberOfInfants:        data.NumberOfInfants,
			BestDepartureTime:      data.BestDepartureTime,
			StopoverPreferences:    data.StopoverPreferences,
			PreferredContactMethod: data.PreferredContactMethod,
			PreferredAirlines:      data.PreferredAirlines,
			ExtraLuggage:           data.ExtraLuggage,
			SpecialMeals:           data.SpecialMeals,
		},
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	if err := l.repo.Add(ctx, flightTicketRequest); err != nil {
		return err
	}

	return nil

}
func (l *flightTicketRequestsvcs) AddMany(ctx context.Context, data []models.FlightTicketRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	var writeOps []mongo.WriteModel

	for _, flr := range data {
		flightTicketRequest := &models.FlightTicketRequest{
			FlightTicketRequestDto: models.FlightTicketRequestDto{
				DestinationFrom:        flr.DestinationFrom,
				DestinationTo:          flr.DestinationTo,
				TripType:               flr.TripType,
				TravelClass:            flr.TravelClass,
				NumberOfAdults:         flr.NumberOfAdults,
				NumberOfChildren:       flr.NumberOfChildren,
				NumberOfInfants:        flr.NumberOfInfants,
				BestDepartureTime:      flr.BestDepartureTime,
				StopoverPreferences:    flr.StopoverPreferences,
				PreferredContactMethod: flr.PreferredContactMethod,
				PreferredAirlines:      flr.PreferredAirlines,
				ExtraLuggage:           flr.ExtraLuggage,
				SpecialMeals:           flr.SpecialMeals,
			},
			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}
		writeOp := mongo.NewInsertOneModel()
		writeOp.SetDocument(flightTicketRequest)
		writeOps = append(writeOps, writeOp)
		if len(writeOps) == 0 {
			return errors.New("empty write ops")
		}
	}
	_, errInsrt := l.repo.BulkWrite(ctx, writeOps)
	if errInsrt != nil {
		return err
	}

	return nil
}
func (a *flightTicketRequestsvcs) Update(ctx context.Context, id string, data *models.FlightTicketRequestDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *flightTicketRequestsvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
