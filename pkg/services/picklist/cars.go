package picklist

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/picklist/filter"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/services/picklist/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Cars, error)
	GetAll(ctx context.Context, query any) ([]models.Cars, error)
	GetPaginated(ctx context.Context, skip, limit int64, query any) (*models.CarsPagination, error)
	Add(ctx context.Context, data *models.CarsDto) (*models.Cars, error)
	Update(ctx context.Context, id string, data *models.CarsDto) (*models.Cars, error)
	Delete(ctx context.Context, id string) error
	//v2
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Cars, error)
	GetPaginatedV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.CarsPagination, error)
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

func (c *carsSvcs) GetPaginated(ctx context.Context, skip, limit int64, query any) (*models.CarsPagination, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.CarsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := c.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Cars
	errAg := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.CarsPagination{
		Cars:       result,
		Pagination: pagination,
	}, nil
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

// v2
func (c *carsSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Cars, error) {
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

	var result []models.Cars
	errAg := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (c *carsSvcs) GetPaginatedV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.CarsPagination, error) {
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

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := c.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Cars
	errAg := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.CarsPagination{
		Cars:       result,
		Pagination: pagination,
	}, nil
}
