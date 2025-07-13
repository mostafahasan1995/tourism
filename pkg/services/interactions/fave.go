package interactions

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/interactions/filter"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FaveSvcs interface {
	Add(ctx context.Context, data *models.FaveDto) (*models.Fave, error)
	Patch(ctx context.Context, id string, data *models.FaveDto) (*models.Fave, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, query any) ([]models.Fave, error)
	Get(ctx context.Context, skip, limit int64, query any) (*models.FavePagination, error)
	Fav(ctx context.Context, data *models.FaveDto) (string, error)
	GetAllByType(ctx context.Context, faveType models.FaveType) ([]models.FaveItem, error)
}

type faveSvcs struct {
	repo repo.FaveRepo
}

func NewFaveSvcs(i *do.Injector) (FaveSvcs, error) {
	return &faveSvcs{
		repo: do.MustInvoke[repo.FaveRepo](i),
	}, nil
}

func (s *faveSvcs) Fav(ctx context.Context, data *models.FaveDto) (string, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return "", err
	}
	existingFave, err := s.repo.GetByFilter(ctx, bson.M{"userId": cfg.User.Id, "refId": data.RefId, "type": data.Type})
	if err != nil && err.Error() != "mongo: no documents in result" {
		return "", err
	}

	if existingFave != nil {
		err := s.repo.DeleteMain(ctx, bson.M{"_id": existingFave.Id})
		if err != nil {
			return "", err
		}
		return "Removed From Favorites", nil
	}

	data.IsFav = true
	data.Type = models.FaveType(data.Type)
	fave := &models.Fave{
		Id:        primitive.NewObjectID(),
		UserId:    cfg.User.Id,
		FaveDto:   *data,
		CreatedAt: time.Now(),
	}

	err = s.repo.Add(ctx, fave)

	if err != nil {
		return "", err
	}

	return "Added To Favorites", nil
}

func (s *faveSvcs) Add(ctx context.Context, data *models.FaveDto) (*models.Fave, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	}

	fave := &models.Fave{
		Id:        primitive.NewObjectID(),
		UserId:    userId,
		FaveDto:   *data,
		CreatedAt: time.Now(),
	}

	err = s.repo.Add(ctx, fave)
	if err != nil {
		return nil, err
	}

	return fave, nil
}

func (s *faveSvcs) Patch(ctx context.Context, id string, data *models.FaveDto) (*models.Fave, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objId}

	// Build update document with only provided fields
	updateFields := bson.M{}
	if data.Type != "" {
		updateFields["type"] = data.Type
	}
	if !data.RefId.IsZero() {
		updateFields["refId"] = data.RefId
	}
	// For boolean field, we always update it since it's either true or false
	// If you need to distinguish between "not provided" and "false",
	// consider using a pointer to bool in the DTO
	updateFields["isFav"] = data.IsFav

	if len(updateFields) == 0 {
		return nil, errors.New("no fields to update")
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *faveSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *faveSvcs) GetAll(ctx context.Context, query any) ([]models.Fave, error) {
	match := bson.M{"trash": bson.M{"$ne": true}}

	filters, err := helpers.ParseFilters[filter.FaveFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	var result []models.Fave
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

func (s *faveSvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.FavePagination, error) {
	match := bson.M{"trash": bson.M{"$ne": true}}

	filters, err := helpers.ParseFilters[filter.FaveFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Fave
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

	return &models.FavePagination{
		Faves:      result,
		Pagination: pagination,
	}, nil
}

func (s *faveSvcs) GetAllByType(ctx context.Context, faveType models.FaveType) ([]models.FaveItem, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	}

	// Determine the collection name based on the favorite type
	var collectionName string
	switch faveType {
	case models.FaveTypeProgram:
		collectionName = "tourismPrograms"
	case models.FaveTypeHotel:
		collectionName = "tourismHotels"
	case models.FaveTypeDiary:
		collectionName = "tourismDiaries"
	case models.FaveTypeExhibition:
		collectionName = "tourismExhibitions"
	case models.FaveTypeAgent:
		collectionName = "tourismAgents"
	default:
		return []models.FaveItem{}, nil
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":   faveType,
				"isFav":  true,
				"userId": userId,
				"trash":  bson.M{"$ne": true},
			},
		},
		{
			"$lookup": bson.M{
				"from":         collectionName,
				"localField":   "refId",
				"foreignField": "_id",
				"as":           "item",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$item",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	var result []models.FaveItem
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
