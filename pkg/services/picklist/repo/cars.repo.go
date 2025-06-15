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
	GetAll(ctx context.Context, filter filter.CarsFilter) (*models.CarsPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.CarsDto) (*models.Cars, error)
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

func (l *carsrepo) GetAll(ctx context.Context, f filter.CarsFilter) (*models.CarsPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	// Build filter using pipeline pattern for consistency
	pipeline := f.BuildPipeline(bson.M{})
	match := pipeline[0]["$match"].(bson.M)

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, match)
	if err != nil {
		return nil, err
	}

	// Pagination defaults and limits
	page := f.Page
	if page <= 0 {
		page = 1
	}
	size := f.Size
	if size <= 0 {
		size = 10 // Default page size
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Query options with pagination
	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"_id": -1})

	cur, err := coll.Find(ctx, match, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var cars []models.Cars
	if err := cur.All(ctx, &cars); err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := &models.CarsPagination{
		Cars: cars,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *carsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.CarsDto) (*models.Cars, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	preCars, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return nil, err
	}

	cars := &models.Cars{
		CarsDto: models.CarsDto{
			CarType: data.CarType,
			Image:   data.Image,
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

	coll := l.db.Database(cfg.Db).Collection(l.collName)
	var updatedCars models.Cars
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedCars); err != nil {
		return nil, err
	}

	return &updatedCars, nil
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
		"updatedBy": cfg.User.Id,
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
