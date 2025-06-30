package picklist

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
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

type CarsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Cars, error)
	GetAll(ctx context.Context, query any) ([]models.Cars, error)
	GetPaginated(ctx context.Context, query any) (*models.CarsPagination, error)
	Add(ctx context.Context, data *models.CarsDto) (*models.Cars, error)
	Update(ctx context.Context, id string, data *models.CarsDto) (*models.Cars, error)
	Delete(ctx context.Context, id string) error
}

type carsSvcs struct {
	repo repo.CarsRepo
}

func NewCarsSvcs(i *do.Injector) (CarsSvcs, error) {
	return &carsSvcs{
		repo: do.MustInvoke[repo.CarsRepo](i),
	}, nil
}

func (c *carsSvcs) GetOne(ctx context.Context, id string) (*models.Cars, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return c.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (c *carsSvcs) GetAll(ctx context.Context, query any) ([]models.Cars, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.CarsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	var result []models.Cars
	errAg := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (c *carsSvcs) GetPaginated(ctx context.Context, query any) (*models.CarsPagination, error) {
	filters, err := helpers.ParseFilters[filter.CarsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	return c.repo.GetAll(ctx, *filters)
}

func (c *carsSvcs) Add(ctx context.Context, data *models.CarsDto) (*models.Cars, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	cars := &models.Cars{
		CarsDto:   *data,
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	if err := c.repo.Add(ctx, cars); err != nil {
		return nil, err
	}

	return cars, nil
}

func (c *carsSvcs) Update(ctx context.Context, id string, data *models.CarsDto) (*models.Cars, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	cars := &models.Cars{
		Id:        _id,
		CarsDto:   *data,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": cars}

	updatedCars, err := c.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedCars, nil
}

func (c *carsSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = c.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
