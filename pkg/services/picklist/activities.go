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

type ActivitiesSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Activities, error)
	GetAll(ctx context.Context, query string) ([]models.Activities, error)
	Add(ctx context.Context, data *models.ActivitiesDto) (*models.Activities, error)
	Update(ctx context.Context, id string, data *models.ActivitiesDto) (*models.Activities, error)
	Delete(ctx context.Context, id string) error
}

type activitiesSvcs struct {
	repo repo.ActivitiesRepo
}

func NewActivitiesSvcs(i *do.Injector) (ActivitiesSvcs, error) {
	return &activitiesSvcs{
		repo: do.MustInvoke[repo.ActivitiesRepo](i),
	}, nil
}

func (a *activitiesSvcs) GetOne(ctx context.Context, id string) (*models.Activities, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return a.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (a *activitiesSvcs) GetAll(ctx context.Context, query string) ([]models.Activities, error) {
	match := bson.M{"trash": false}

	filters, err := filter.NewActivitiesFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	var result []models.Activities
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (a *activitiesSvcs) Add(ctx context.Context, data *models.ActivitiesDto) (*models.Activities, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	activities := &models.Activities{
		Id:            primitive.NewObjectID(),
		ActivitiesDto: *data,
		CreatedAt:     time.Now(),
		CreatedBy:     cfg.User.Id,
	}

	if err := a.repo.Add(ctx, activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func (a *activitiesSvcs) Update(ctx context.Context, id string, data *models.ActivitiesDto) (*models.Activities, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	activities := &models.Activities{
		Id:            _id,
		ActivitiesDto: *data,
		UpdatedAt:     time.Now(),
		UpdatedBy:     cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": activities}

	updatedActivities, err := a.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedActivities, nil
}

func (a *activitiesSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
