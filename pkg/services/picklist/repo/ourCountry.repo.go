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

type OurCountryRepo interface {
	dbrepo.MainRepo[models.OurCountry]
	GetOne(ctx context.Context, id string) (*models.OurCountry, error)
	GetAll(ctx context.Context, filter filter.OurCountryFilter) (models.OurCountryPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.OurCountryDto) error
	Delete(ctx context.Context, id string) error
}

type ourCountryrepo struct {
	dbrepo.MainRepoImpl[models.OurCountry]

	db       *mongo.Client
	collName string
}

func NewOurCountryRepo(i *do.Injector) (OurCountryRepo, error) {
	return &ourCountryrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.OurCountry]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismOurCountry",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismOurCountry",
	}, nil
}

func (l *ourCountryrepo) GetOne(ctx context.Context, id string) (*models.OurCountry, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.OurCountry
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *ourCountryrepo) GetAll(ctx context.Context, filter filter.OurCountryFilter) (models.OurCountryPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.OurCountryPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.OurCountryPagination{}, err
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
		return models.OurCountryPagination{}, err
	}

	var programs []models.OurCountry
	if err := cur.All(ctx, &programs); err != nil {
		return models.OurCountryPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}


	result := models.OurCountryPagination {
		OurCountry:programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}


	return result, nil
}

func (l *ourCountryrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.OurCountryDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preOurCountry, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	ourCountry := &models.OurCountry{
		OurCountryDto: models.OurCountryDto{
			Name:                          data.Name,
			Image:                          data.Image,

		},
		Id:        id,
		Trash:     false,
		CreatedAt: preOurCountry.CreatedAt,
		CreatedBy: preOurCountry.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": ourCountry}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedOurCountry models.OurCountry
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedOurCountry); err != nil {
		return err
	}

	return nil
}
func (l *ourCountryrepo) Delete(ctx context.Context, id string) error {

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
