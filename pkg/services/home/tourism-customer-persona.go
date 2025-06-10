package home

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerPersonaSvcs interface {
	GetOne(ctx context.Context, id string) (*models.CustomerPersona, error)
	GetAll(ctx context.Context) ([]models.CustomerPersona, error)
	Get(ctx context.Context, skip int64, limit int64, queryString string) (*mongo.Cursor, int64, error)
	Add(ctx context.Context, data *models.CustomerPersonaDto) (*models.CustomerPersona, error)
	Update(ctx context.Context, id string, data *models.CustomerPersonaDto) (*models.CustomerPersona, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.CustomerPersona, error)
	Delete(ctx context.Context, id string) error
}

type customerPersonaSvcs struct {
	repo repo.CustomerPersonaRepo
}

func NewCustomerPersonaSvcs(i *do.Injector) (CustomerPersonaSvcs, error) {
	return &customerPersonaSvcs{
		repo: do.MustInvoke[repo.CustomerPersonaRepo](i),
	}, nil
}

func (s *customerPersonaSvcs) GetOne(ctx context.Context, id string) (*models.CustomerPersona, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *customerPersonaSvcs) GetAll(ctx context.Context) ([]models.CustomerPersona, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	var result []models.CustomerPersona
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *customerPersonaSvcs) Get(ctx context.Context, skip int64, limit int64, queryString string) (*mongo.Cursor, int64, error) {
	// Create filter from query string
	filter, err := filter.NewCustomerPersonaFilter(queryString)
	if err != nil {
		return nil, 0, err
	}

	// Get base filter and build pipeline
	baseFilter := bson.M{}
	pipeline := filter.BuildPipeline(baseFilter)

	// Create a copy of the pipeline for counting
	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	// Count total documents using the same pipeline
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination to the main pipeline
	if len(pipeline) > 0 && pipeline[len(pipeline)-1]["$sort"] == nil {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"displayOrder": 1}})
	}
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// Execute query
	var cursor *mongo.Cursor
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		cursor = cur
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return cursor, count, nil
}

func (s *customerPersonaSvcs) Add(ctx context.Context, data *models.CustomerPersonaDto) (*models.CustomerPersona, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	persona := &models.CustomerPersona{
		CustomerPersonaDto: *data,
		Id:                 primitive.NewObjectID(),
		Trash:              false,
		CreatedAt:          time.Now(),
		CreatedBy:          cfg.User.Id,
		UpdatedAt:          time.Now(),
		UpdatedBy:          cfg.User.Id,
	}

	if err := s.repo.Add(ctx, persona); err != nil {
		return nil, err
	}

	return persona, nil
}

func (s *customerPersonaSvcs) Update(ctx context.Context, id string, data *models.CustomerPersonaDto) (*models.CustomerPersona, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get the existing persona first to preserve creation info
	existingPersona, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	persona := &models.CustomerPersona{
		CustomerPersonaDto: *data,
		Id:                 _id,
		Trash:              false,
		CreatedAt:          existingPersona.CreatedAt,
		CreatedBy:          existingPersona.CreatedBy,
		UpdatedAt:          time.Now(),
		UpdatedBy:          cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": persona}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *customerPersonaSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.CustomerPersona, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}

	for key, value := range updates {
		updateDoc[key] = value
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *customerPersonaSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}
