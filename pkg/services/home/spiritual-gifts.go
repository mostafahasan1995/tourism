package home

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SpiritualGiftSvcs defines the interface for spiritual gift services
type SpiritualGiftSvcs interface {
	Get(ctx context.Context, skip, limit int64, query any) (*models.SpiritualGiftWithPagination, error)
	GetById(ctx context.Context, id string) (*models.SpiritualGift, error)
	GetCurrent(ctx context.Context) (*models.SpiritualGift, error)
	Save(ctx context.Context, data *models.SpiritualGiftDto) (*models.SpiritualGift, error)
	Update(ctx context.Context, id string, data *models.SpiritualGiftDto) (*models.SpiritualGift, error)
	Toggle(ctx context.Context, id string, enabled bool) (*models.SpiritualGift, error)
	Delete(ctx context.Context, id string) error
}

// spiritualGiftSvcs implements the SpiritualGiftSvcs interface
type spiritualGiftSvcs struct {
	repo repo.SpiritualGiftRepo
}

// NewSpiritualGiftSvcs creates a new instance of SpiritualGiftSvcs
func NewSpiritualGiftSvcs(i *do.Injector) (SpiritualGiftSvcs, error) {
	return &spiritualGiftSvcs{
		repo: do.MustInvoke[repo.SpiritualGiftRepo](i),
	}, nil
}

// Get retrieves spiritual gifts with pagination and filtering
func (s *spiritualGiftSvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.SpiritualGiftWithPagination, error) {
	match := bson.M{"trash": false}

	f, err := helpers.ParseFilters[filter.SpiritualGiftFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.SpiritualGift
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.SpiritualGiftWithPagination{
		SpiritualGifts: result,
		Pagination:     pagination,
	}, nil
}

// GetById retrieves a spiritual gift by ID
func (s *spiritualGiftSvcs) GetById(ctx context.Context, id string) (*models.SpiritualGift, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "trash": false}
	return s.repo.GetByFilter(ctx, filter)
}

// GetCurrent retrieves the current active spiritual gift configuration
func (s *spiritualGiftSvcs) GetCurrent(ctx context.Context) (*models.SpiritualGift, error) {
	// Find the most recently created or updated spiritual gift
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$sort": bson.M{"updatedAt": -1, "createdAt": -1}},
		{"$limit": 1},
	}

	var results []models.SpiritualGift
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &results)
	})

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		// Return a default disabled configuration if none exists
		return &models.SpiritualGift{
			SpiritualGiftDto: models.SpiritualGiftDto{
				Enabled: false,
			},
		}, nil
	}

	return &results[0], nil
}

// Save creates a new spiritual gift record
func (s *spiritualGiftSvcs) Save(ctx context.Context, data *models.SpiritualGiftDto) (*models.SpiritualGift, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	spiritualGift := &models.SpiritualGift{
		Id:               primitive.NewObjectID(),
		SpiritualGiftDto: *data,
		Trash:            false,
		CreatedAt:        time.Now(),
		CreatedBy:        cfg.User.Id,
		UpdatedAt:        time.Now(),
		UpdatedBy:        cfg.User.Id,
	}

	if err := s.repo.Add(ctx, spiritualGift); err != nil {
		return nil, err
	}

	return spiritualGift, nil
}

// Update updates an existing spiritual gift record
func (s *spiritualGiftSvcs) Update(ctx context.Context, id string, data *models.SpiritualGiftDto) (*models.SpiritualGift, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Check if record exists
	existing, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the record
	spiritualGift := &models.SpiritualGift{
		Id:               objID,
		SpiritualGiftDto: *data,
		Trash:            false,
		CreatedAt:        existing.CreatedAt,
		CreatedBy:        existing.CreatedBy,
		UpdatedAt:        time.Now(),
		UpdatedBy:        cfg.User.Id,
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": spiritualGift}

	if _, err := s.repo.Patch(ctx, filter, update); err != nil {
		return nil, err
	}

	return spiritualGift, nil
}

// Toggle enables or disables a spiritual gift record
func (s *spiritualGiftSvcs) Toggle(ctx context.Context, id string, enabled bool) (*models.SpiritualGift, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"enabled":   enabled,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	if _, err := s.repo.Patch(ctx, filter, update); err != nil {
		return nil, err
	}

	return s.GetById(ctx, id)
}

// Delete marks a spiritual gift record as trash
func (s *spiritualGiftSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}
