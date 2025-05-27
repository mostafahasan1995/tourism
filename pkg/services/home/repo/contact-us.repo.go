package repo

import (
	"context"
	"fmt"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"git.larsa.io/mahdawi/microservices-commons/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContactUsRepo interface {
	dbrepo.MainRepo[models.ContactUs]
	GetOne(ctx context.Context, id string) (*models.ContactUs, error)
	GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.ContactUsDto) error
	Delete(ctx context.Context, id string) error
}

type contactUsrepo struct {
	dbrepo.MainRepoImpl[models.ContactUs]

	db       *mongo.Client
	collName string
}

func NewContactUsRepo(i *do.Injector) (ContactUsRepo, error) {
	return &contactUsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.ContactUs]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismContactUs",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismContactUs",
	}, nil
}

func (l *contactUsrepo) GetOne(ctx context.Context, id string) (*models.ContactUs, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.ContactUs
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *contactUsrepo) GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.ContactUsPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
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

	// Query options with pagination
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.ContactUsPagination{}, err
	}

	var programs []models.ContactUs
	if err := cur.All(ctx, &programs); err != nil {
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

func (l *contactUsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preContactUs, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	contactUs := &models.ContactUs{
		ContactUsDto: models.ContactUsDto{
			FullName:        data.FullName,
			EmailAddress:    data.EmailAddress,
			PhoneNumber:     data.PhoneNumber,
			HowDidYouFindUs: data.HowDidYouFindUs,
			Message:         data.Message,
		},

		Id:        id,
		Trash:     false,
		CreatedAt: preContactUs.CreatedAt,
		CreatedBy: preContactUs.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": contactUs}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedContactUs models.ContactUs
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedContactUs); err != nil {
		return err
	}

	return nil
}
func (l *contactUsrepo) Delete(ctx context.Context, id string) error {

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

func (l *contactUsrepo) Add(ctx context.Context, data *models.ContactUs, opts ...*options.InsertOneOptions) error {
	fmt.Printf("Starting repo Add method with data: %+v\n", data)

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		fmt.Printf("Error getting config in repo: %v\n", err)
		return err
	}
	fmt.Printf("Got config with DB in repo: %s\n", cfg.Db)

	coll := l.db.Database(cfg.Db).Collection(l.collName)
	fmt.Printf("Using collection: %s\n", l.collName)

	result, err := coll.InsertOne(ctx, data, opts...)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return err
	}
	fmt.Printf("Successfully inserted document with ID: %v\n", result.InsertedID)

	return nil
}
