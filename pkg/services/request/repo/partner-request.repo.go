package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/request/filter"
	"larsa-tourism-microservices/pkg/services/request/models"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartnerRequestRepo interface {
	dbrepo.MainRepo[models.PartnerRequest]
	GetOne(ctx context.Context, id string) (*models.PartnerRequest, error)
	GetAll(ctx context.Context, filter filter.PartnerRequestFilter) (models.PartnerRequestPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.PartnerRequestDto) error
	Delete(ctx context.Context, id string) error
}

type partnerRequestrepo struct {
	dbrepo.MainRepoImpl[models.PartnerRequest]

	db       *mongo.Client
	collName string
}

func NewPartnerRequestRepo(i *do.Injector) (PartnerRequestRepo, error) {
	return &partnerRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.PartnerRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismPartnerRequest",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismPartnerRequest",
	}, nil
}

func (l *partnerRequestrepo) GetOne(ctx context.Context, id string) (*models.PartnerRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.PartnerRequest
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *partnerRequestrepo) GetAll(ctx context.Context, filter filter.PartnerRequestFilter) (models.PartnerRequestPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.PartnerRequestPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.PartnerRequestPagination{}, err
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

	// Query options with pagination
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.PartnerRequestPagination{}, err
	}

	var programs []models.PartnerRequest
	if err := cur.All(ctx, &programs); err != nil {
		return models.PartnerRequestPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := models.PartnerRequestPagination{
		PartnerRequest: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *partnerRequestrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.PartnerRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	prePartnerRequest, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	partnerRequest := &models.PartnerRequest{
		PartnerRequestDto: models.PartnerRequestDto{
			CompanyName:     data.CompanyName,
			BusinessType:    data.BusinessType,
			Website:         data.Website,
			CompanyLocation: data.CompanyLocation,
			ContactDetail: models.ContactDetail{
				FullName:    data.ContactDetail.FullName,
				Position:    data.ContactDetail.Position,
				PhoneNumber: data.ContactDetail.PhoneNumber,
				Email:       data.ContactDetail.Email,
			},
		},

		Id:        id,
		Trash:     false,
		CreatedAt: prePartnerRequest.CreatedAt,
		CreatedBy: prePartnerRequest.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": partnerRequest}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedPartnerRequest models.PartnerRequest
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedPartnerRequest); err != nil {
		return err
	}

	return nil
}
func (l *partnerRequestrepo) Delete(ctx context.Context, id string) error {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": primitive.NilObjectID,
	}}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&opt,
	)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
