package home

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrustedPartnersSvcs interface {
	GetOne(ctx context.Context, id string) (*models.TrustedPartner, error)
	GetAll(ctx context.Context) ([]models.TrustedPartner, error)
	Get(ctx context.Context, skip int64, limit int64, queryString string) (*mongo.Cursor, int64, error)
	Add(ctx context.Context, data *models.TrustedPartnerDto) (*models.TrustedPartner, error)

	Update(ctx context.Context, id string, data *models.TrustedPartnerDto) (*models.TrustedPartner, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.TrustedPartner, error)
	Delete(ctx context.Context, id string) error

	// V2
	GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.TrustedPartnerPagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.TrustedPartner, error)
}

type trustedPartnerssvcs struct {
	repo repo.TrustedPartnersRepo
}

func NewTrustedPartnersSvcs(i *do.Injector) (TrustedPartnersSvcs, error) {
	return &trustedPartnerssvcs{
		repo: do.MustInvoke[repo.TrustedPartnersRepo](i),
	}, nil
}

func (s *trustedPartnerssvcs) GetOne(ctx context.Context, id string) (*models.TrustedPartner, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *trustedPartnerssvcs) GetAll(ctx context.Context) ([]models.TrustedPartner, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	var result []models.TrustedPartner
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

func (s *trustedPartnerssvcs) Get(ctx context.Context, skip int64, limit int64, queryString string) (*mongo.Cursor, int64, error) {
	// Create filter from query string
	filter, err := filter.NewTrustedPartnersFilter(queryString)
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

func (s *trustedPartnerssvcs) Add(ctx context.Context, data *models.TrustedPartnerDto) (*models.TrustedPartner, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	partner := &models.TrustedPartner{
		TrustedPartnerDto: models.TrustedPartnerDto{
			Title:        data.Title,
			ProgramTitle: data.ProgramTitle,
			ProgramId:    data.ProgramId,
			Description:  data.Description,
			Image:        data.Image,
			URL:          data.URL,
			DisplayOrder: data.DisplayOrder,
			IsActive:     data.IsActive,
		},

		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	if err := s.repo.Add(ctx, partner); err != nil {
		return nil, err
	}

	return partner, nil
}

func (s *trustedPartnerssvcs) Update(ctx context.Context, id string, data *models.TrustedPartnerDto) (*models.TrustedPartner, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get the existing partner first to preserve creation info
	existingPartner, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	partner := &models.TrustedPartner{
		TrustedPartnerDto: *data,
		Id:                _id,
		Trash:             false,
		CreatedAt:         existingPartner.CreatedAt,
		CreatedBy:         existingPartner.CreatedBy,
		UpdatedAt:         time.Now(),
		UpdatedBy:         cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": partner}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *trustedPartnerssvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.TrustedPartner, error) {
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

func (s *trustedPartnerssvcs) Delete(ctx context.Context, id string) error {
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

func (s *trustedPartnerssvcs) GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.TrustedPartnerPagination, error) {
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
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.TrustedPartner
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}
	return &models.TrustedPartnerPagination{
		Partners: result,
		Pagination: common.Pagination{
			TotalPages: math.Ceil(float64(count) / float64(limit)),
			PerPage:    limit,
			TotalCount: count,
		},
	}, nil
}
func (s *trustedPartnerssvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.TrustedPartner, error) {
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

	var result []models.TrustedPartner
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
