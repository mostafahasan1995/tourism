package picklist

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/picklist/filter"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/services/picklist/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Destination, error)
	GetAll(ctx context.Context, query any) ([]models.Destination, error)
	Add(ctx context.Context, data *models.DestinationDto) (*models.Destination, error)
	Update(ctx context.Context, id string, data *models.DestinationDto) (*models.Destination, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	//v2
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Destination, error)
}

type destinationSvcs struct {
	repo repo.DestinationRepo
}

func NewDestinationSvcs(i *do.Injector) (DestinationSvcs, error) {
	return &destinationSvcs{
		repo: do.MustInvoke[repo.DestinationRepo](i),
	}, nil
}

func (d *destinationSvcs) GetOne(ctx context.Context, id string) (*models.Destination, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return d.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (d *destinationSvcs) GetAll(ctx context.Context, query any) ([]models.Destination, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.DestinationFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	var result []models.Destination
	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (d *destinationSvcs) Add(ctx context.Context, data *models.DestinationDto) (*models.Destination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	destination := &models.Destination{
		Id:             primitive.NewObjectID(),
		DestinationDto: *data,
		CreatedAt:      time.Now(),
		CreatedBy:      cfg.User.Id,
	}

	if err := d.repo.Add(ctx, destination); err != nil {
		return nil, err
	}

	return destination, nil
}

func (d *destinationSvcs) Update(ctx context.Context, id string, data *models.DestinationDto) (*models.Destination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	destination := &models.Destination{
		Id:             _id,
		DestinationDto: *data,
		UpdatedAt:      time.Now(),
		UpdatedBy:      cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": destination}

	updatedDestination, err := d.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedDestination, nil
}

func (d *destinationSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = d.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *destinationSvcs) Count(ctx context.Context, filter any) (int64, error) {
	return d.repo.Count(ctx, filter)
}

// v2
func (d *destinationSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Destination, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	var result []models.Destination
	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}
