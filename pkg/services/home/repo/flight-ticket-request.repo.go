package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FlightTicketRequestRepo interface {
	dbrepo.MainRepo[models.FlightTicketRequest]
	GetOne(ctx context.Context, id string) (*models.FlightTicketRequest, error)
	GetAll(ctx context.Context, filter filter.FlightTicketRequestFilter) (models.FlightTicketRequestPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.FlightTicketRequestDto) error
	Delete(ctx context.Context, id string) error
}

type flightTicketRequestrepo struct {
	dbrepo.MainRepoImpl[models.FlightTicketRequest]

	db       *mongo.Client
	collName string
}

func NewFlightTicketRequestRepo(i *do.Injector) (FlightTicketRequestRepo, error) {
	return &flightTicketRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FlightTicketRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFlightTicketRequest",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismFlightTicketRequest",
	}, nil
}

func (l *flightTicketRequestrepo) GetOne(ctx context.Context, id string) (*models.FlightTicketRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.FlightTicketRequest
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *flightTicketRequestrepo) GetAll(ctx context.Context, filter filter.FlightTicketRequestFilter) (models.FlightTicketRequestPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.FlightTicketRequestPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.FlightTicketRequestPagination{}, err
	}

	// Pagination defaults and limits
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = int(totalCount) // return all if invalid
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Query options with pagination
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.FlightTicketRequestPagination{}, err
	}

	var programs []models.FlightTicketRequest
	if err := cur.All(ctx, &programs); err != nil {
		return models.FlightTicketRequestPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := models.FlightTicketRequestPagination{
		FlightTicketRequest: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *flightTicketRequestrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.FlightTicketRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preFlightTicketRequest, err := l.GetOne(ctx, id.Hex())
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

		Id:        id,
		Trash:     false,
		CreatedAt: preFlightTicketRequest.CreatedAt,
		CreatedBy: preFlightTicketRequest.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": flightTicketRequest}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedFlightTicketRequest models.FlightTicketRequest
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedFlightTicketRequest); err != nil {
		return err
	}

	return nil
}
func (l *flightTicketRequestrepo) Delete(ctx context.Context, id string) error {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": primitive.NilObjectID,
	}}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&opt,
	)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
