package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/picklist/filter"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CarsRepo interface {
	dbrepo.MainRepo[models.Cars]
	GetOne(ctx context.Context, id string) (*models.Cars, error)
	GetAll(ctx context.Context, filter filter.CarsFilter) (models.CarsPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.CarsDto) error
	Delete(ctx context.Context, id string) error
}

type carsrepo struct {
	dbrepo.MainRepoImpl[models.Cars]

	db       *mongo.Client
	collName string
}

func NewCarsRepo(i *do.Injector) (CarsRepo, error) {
	return &carsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Cars]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCars",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismCars",
	}, nil
}

func (l *carsrepo) GetOne(ctx context.Context, id string) (*models.Cars, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.Cars
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *carsrepo) GetAll(ctx context.Context, filter filter.CarsFilter) (models.CarsPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.CarsPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.CarsPagination{}, err
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
		return models.CarsPagination{}, err
	}

	var programs []models.Cars
	if err := cur.All(ctx, &programs); err != nil {
		return models.CarsPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}


	result := models.CarsPagination {
		Cars:programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}


	return result, nil
}

func (l *carsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.CarsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preCars, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	cars := &models.Cars{
		CarsDto: models.CarsDto{
			CarType:                          data.CarType,
			Image:                          data.Image,

		},
		Id:        id,
		Trash:     false,
		CreatedAt: preCars.CreatedAt,
		CreatedBy: preCars.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": cars}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedCars models.Cars
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedCars); err != nil {
		return err
	}

	return nil
}
func (l *carsrepo) Delete(ctx context.Context, id string) error {

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
