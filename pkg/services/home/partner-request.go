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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartnerRequestSvcs interface {
	GetOne(ctx context.Context, id string) (*models.PartnerRequest, error)
	GetAll(ctx context.Context) ([]models.PartnerRequest, error)
	Get(ctx context.Context, skip int64, limit int64, query any) (*mongo.Cursor, int64, error)
	Add(ctx context.Context, data *models.PartnerRequestDto) (*models.PartnerRequest, error)
	Update(ctx context.Context, id string, data *models.PartnerRequestDto) (*models.PartnerRequest, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.PartnerRequest, error)
	UpdateStatus(ctx context.Context, id string, status string) (*models.PartnerRequest, error)
	Delete(ctx context.Context, id string) error
}

type partnerRequestSvcs struct {
	repo repo.PartnerRequestRepo
}

func NewPartnerRequestSvcs(i *do.Injector) (PartnerRequestSvcs, error) {
	return &partnerRequestSvcs{
		repo: do.MustInvoke[repo.PartnerRequestRepo](i),
	}, nil
}

func (s *partnerRequestSvcs) GetOne(ctx context.Context, id string) (*models.PartnerRequest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *partnerRequestSvcs) GetAll(ctx context.Context) ([]models.PartnerRequest, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$sort": bson.M{"createdAt": -1}}, // Sort by creation date, newest first
	}

	var result []models.PartnerRequest
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

func (s *partnerRequestSvcs) Get(ctx context.Context, skip int64, limit int64, query any) (*mongo.Cursor, int64, error) {
	// Create filter from query
	filters, err := helpers.ParseFilters[filter.PartnerRequestFilter](query)
	if err != nil {
		return nil, 0, errors.New("invalid query")
	}

	// Get base filter and build pipeline
	baseFilter := bson.M{}
	pipeline := filters.BuildPipeline(baseFilter)

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
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
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

func (s *partnerRequestSvcs) Add(ctx context.Context, data *models.PartnerRequestDto) (*models.PartnerRequest, error) {
	// Default status is pending for new requests
	if data.Status == "" {
		data.Status = "pending"
	}

	request := &models.PartnerRequest{
		PartnerRequestDto: *data,
		Id:                primitive.NewObjectID(),
		Trash:             false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Try to get user info if authenticated
	cfg, _ := util.GetReqAppCfg(ctx)

	request.CreatedBy = cfg.User.Id
	request.UpdatedBy = cfg.User.Id

	if err := s.repo.Add(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *partnerRequestSvcs) Update(ctx context.Context, id string, data *models.PartnerRequestDto) (*models.PartnerRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get the existing request first to preserve creation info
	existingRequest, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	request := &models.PartnerRequest{
		PartnerRequestDto: *data,
		Id:                _id,
		Trash:             false,
		CreatedAt:         existingRequest.CreatedAt,
		CreatedBy:         existingRequest.CreatedBy,
		UpdatedAt:         time.Now(),
		UpdatedBy:         cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": request}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *partnerRequestSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.PartnerRequest, error) {
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

func (s *partnerRequestSvcs) UpdateStatus(ctx context.Context, id string, status string) (*models.PartnerRequest, error) {
	if status != "pending" && status != "approved" && status != "rejected" {
		return nil, helpers.BadRequest("Invalid status. Must be one of: pending, approved, rejected")
	}

	updates := map[string]interface{}{
		"status": status,
	}

	return s.Patch(ctx, id, updates)
}

func (s *partnerRequestSvcs) Delete(ctx context.Context, id string) error {
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
