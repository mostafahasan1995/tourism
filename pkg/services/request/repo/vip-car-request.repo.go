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

type VipCarRequestRepo interface {
	dbrepo.MainRepo[models.VipCarRequest]
	GetOne(ctx context.Context, id string) (*models.VipCarRequest, error)
	GetAll(ctx context.Context, filter filter.VipCarRequestFilter) (models.VipCarRequestPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.VipCarRequestDto) error
	Delete(ctx context.Context, id string) error
}

type vipCarRequestrepo struct {
	dbrepo.MainRepoImpl[models.VipCarRequest]
	db       *mongo.Client
	collName string
}

func NewVipCarRequestRepo(i *do.Injector) (VipCarRequestRepo, error) {
	return &vipCarRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.VipCarRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismVipCarRequest",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismVipCarRequest",
	}, nil
}

func (l *vipCarRequestrepo) GetOne(ctx context.Context, id string) (*models.VipCarRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.VipCarRequest
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *vipCarRequestrepo) GetAll(ctx context.Context, filter filter.VipCarRequestFilter) (models.VipCarRequestPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.VipCarRequestPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.VipCarRequestPagination{}, err
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
		return models.VipCarRequestPagination{}, err
	}

	var programs []models.VipCarRequest
	if err := cur.All(ctx, &programs); err != nil {
		return models.VipCarRequestPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := models.VipCarRequestPagination{
		VipCarRequest: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *vipCarRequestrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.VipCarRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preVipCarRequest, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	vipCarRequest := &models.VipCarRequest{
		VipCarRequestDto: models.VipCarRequestDto{
			Location:              data.Location,
			Capacity:              data.Capacity,
			DriverLanguagesSpoken: data.DriverLanguagesSpoken,
			LuxuryFeatures:        data.LuxuryFeatures,
			StartDate:             data.StartDate,
			EndDate:               data.EndDate,
			StartTime:             data.StartTime,
			EndTime:               data.EndTime,
			CarTypeId:             data.CarTypeId,
		},

		Id:        id,
		Trash:     false,
		CreatedAt: preVipCarRequest.CreatedAt,
		CreatedBy: preVipCarRequest.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": vipCarRequest}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedVipCarRequest models.VipCarRequest
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedVipCarRequest); err != nil {
		return err
	}

	return nil
}
func (l *vipCarRequestrepo) Delete(ctx context.Context, id string) error {

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
