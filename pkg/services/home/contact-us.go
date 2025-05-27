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

	"git.larsa.io/mahdawi/microservices-commons/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactUsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ContactUs, error)
	GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error)
	Add(ctx context.Context, data *models.ContactUsDto) error
	AddMany(ctx context.Context, data []models.ContactUsDto) error
	Update(ctx context.Context, id string, data *models.ContactUsDto) error
	Delete(ctx context.Context, id string) error
}

type contactUssvcs struct {
	repo repo.ContactUsRepo
}

func NewContactUsSvcs(i *do.Injector) (ContactUsSvcs, error) {
	return &contactUssvcs{
		repo: do.MustInvoke[repo.ContactUsRepo](i),
	}, nil
}

func (l *contactUssvcs) GetOne(ctx context.Context, id string) (*models.ContactUs, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return l.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (l *contactUssvcs) GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error) {
	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := l.repo.Count(ctx, filterBody)
	if err != nil {
		return models.ContactUsPagination{}, err
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

	// Create aggregation pipeline for pagination
	pipeline := []bson.M{
		{"$match": filterBody},
		{"$skip": skip},
		{"$limit": limit},
	}

	var programs []models.ContactUs
	err = l.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &programs)
	})
	if err != nil {
		return models.ContactUsPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := models.ContactUsPagination{
		ContactUs: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *contactUssvcs) Add(ctx context.Context, data *models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	contactUs := &models.ContactUs{
		ContactUsDto: *data, // This preserves AdditionalFields
		Id:           primitive.NewObjectID(),
		Trash:        false,
		CreatedAt:    time.Now(),
		CreatedBy:    userId,
		UpdatedAt:    time.Now(),
		UpdatedBy:    userId,
	}

	if err := l.repo.Add(ctx, contactUs); err != nil {
		return err
	}

	return nil
}

func (l *contactUssvcs) AddMany(ctx context.Context, data []models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("empty data array")
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	var contactUsArray []any
	for _, flr := range data {
		contactUs := &models.ContactUs{
			ContactUsDto: flr, // This preserves AdditionalFields
			Id:           primitive.NewObjectID(),
			Trash:        false,
			CreatedAt:    time.Now(),
			CreatedBy:    userId,
			UpdatedAt:    time.Now(),
			UpdatedBy:    userId,
		}
		contactUsArray = append(contactUsArray, contactUs)
	}

	err = l.repo.AddMany(ctx, contactUsArray)
	if err != nil {
		return err
	}

	return nil
}

func (a *contactUssvcs) Update(ctx context.Context, id string, data *models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing contact to preserve created fields
	existing, err := a.GetOne(ctx, id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	contactUs := &models.ContactUs{
		ContactUsDto: *data, // This preserves AdditionalFields
		Id:           _id,
		Trash:        false,
		CreatedAt:    existing.CreatedAt,
		CreatedBy:    existing.CreatedBy,
		UpdatedBy:    userId,
		UpdatedAt:    time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": contactUs}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *contactUssvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
