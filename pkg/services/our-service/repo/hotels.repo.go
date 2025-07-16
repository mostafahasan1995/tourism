package repo

import (
	"context"

	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelsRepo interface {
	dbrepo.MainRepo[models.Hotels]
	GetOne(ctx context.Context, id string) (*models.Hotels, error)
	GetAll(ctx context.Context, filter filter.HotelsFilter, page, perPage int64) (models.HotelsPagination, error)
	GetAllHotels(ctx context.Context, hotelFilter filter.HotelsFilter) ([]models.Hotels, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.HotelsDto) (*models.Hotels, error)
	Delete(ctx context.Context, id string) error
}

type hotelsrepo struct {
	dbrepo.MainRepoImpl[models.Hotels]
	db       *mongo.Client
	collName string
}

func NewHotelsRepo(i *do.Injector) (HotelsRepo, error) {
	return &hotelsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Hotels]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismHotels",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismHotels",
	}, nil
}

func (l *hotelsrepo) GetOne(ctx context.Context, id string) (*models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.Hotels
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	data.CalculateAverageRating()
	return &data, nil
}

func (l *hotelsrepo) GetAllHotels(ctx context.Context, hotelFilter filter.HotelsFilter) ([]models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	// Use the ToBsonFilter method to build the MongoDB filter
	filterBody := hotelFilter.ToBsonFilter()

	findOptions := options.Find().
		SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {

		return nil, err
	}

	var hotels []models.Hotels
	if err := cur.All(ctx, &hotels); err != nil {

		return nil, err
	}

	// Calculate average rating for each hotel
	for i := range hotels {
		hotels[i].CalculateAverageRating()
	}

	// result := models.Hotels{
	// 	Hotels: hotels,
	// 	Pagination: common.Pagination{
	// 		TotalPages: float64(totalPages),
	// 		PerPage:    perPage,
	// 		TotalCount: totalCount,
	// 	},
	// }

	return hotels, nil
}

func (l *hotelsrepo) GetAll(ctx context.Context, hotelFilter filter.HotelsFilter, page, perPage int64) (models.HotelsPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.HotelsPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	// Use the ToBsonFilter method to build the MongoDB filter
	filterBody := hotelFilter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {

		return models.HotelsPagination{}, err
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	skip := (page - 1) * perPage

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(perPage).
		SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {

		return models.HotelsPagination{}, err
	}

	var hotels []models.Hotels
	if err := cur.All(ctx, &hotels); err != nil {

		return models.HotelsPagination{}, err
	}

	// Calculate average rating for each hotel
	for i := range hotels {
		hotels[i].CalculateAverageRating()
	}

	// Prepare pagination result
	totalPages := int64(0)
	if perPage > 0 {
		totalPages = (totalCount + perPage - 1) / perPage
	}

	result := models.HotelsPagination{
		Hotels: hotels,
		Pagination: types.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    perPage,
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *hotelsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.HotelsDto) (*models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preHotels, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return nil, err
	}

	hotels := &models.Hotels{
		HotelsDto: *data,
		Id:        id,
		Trash:     false,
		CreatedAt: preHotels.CreatedAt,
		CreatedBy: preHotels.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": hotels}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedHotels models.Hotels
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedHotels); err != nil {
		return nil, err
	}

	updatedHotels.CalculateAverageRating()
	return &updatedHotels, nil
}

func (l *hotelsrepo) Delete(ctx context.Context, id string) error {
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
