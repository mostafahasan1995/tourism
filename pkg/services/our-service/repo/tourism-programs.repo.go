package repo

import (
	"context"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// deprecated
type TourismProgramRepo interface {
	GetOne(ctx context.Context, id string) (*models.TourismProgram, error)
	Add(ctx context.Context, data *models.TourismProgramDto) error
	GetAll(ctx context.Context, filter filter.TourismProgramFilter) (models.TourismProgramPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.TourismProgramDto) error
	Delete(ctx context.Context, id string) error
}

type tourismProgramrepo struct {
	db       *mongo.Client
	collName string
}

func NewTourismProgramRepo(i *do.Injector) (TourismProgramRepo, error) {
	return &tourismProgramrepo{
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismTourismProgram",
	}, nil
}

func (l *tourismProgramrepo) GetOne(ctx context.Context, id string) (*models.TourismProgram, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.TourismProgram
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil

}

func (l *tourismProgramrepo) GetAll(ctx context.Context, filter filter.TourismProgramFilter) (models.TourismProgramPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.TourismProgramPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.TourismProgramPagination{}, err
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
		return models.TourismProgramPagination{}, err
	}

	var programs []models.TourismProgram
	if err := cur.All(ctx, &programs); err != nil {
		return models.TourismProgramPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	result := models.TourismProgramPagination{
		TourismProgram: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *tourismProgramrepo) Add(ctx context.Context, data *models.TourismProgramDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	tourismProgram := &models.TourismProgram{
		TourismProgramDto: models.TourismProgramDto{ // Initialize the embedded ChatDto
			Name:        data.Name,
			Destination: data.Destination,
			Duration:    data.Duration,
			TravelType:  data.TravelType,
			GroupSize:   data.GroupSize,
			Image:       data.Image,
		},
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	if _, err := coll.InsertOne(ctx, tourismProgram); err != nil {
		return err
	}

	return nil
}

func (l *tourismProgramrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.TourismProgramDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preTourismProgram, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	tourismProgram := &models.TourismProgram{
		TourismProgramDto: models.TourismProgramDto{ // Initialize the embedded ChatDto
			Name:        data.Name,
			Destination: data.Destination,
			Duration:    data.Duration,
			TravelType:  data.TravelType,
			GroupSize:   data.GroupSize,
			Image:       data.Image,
		},
		Id:        id,
		CreatedAt: preTourismProgram.CreatedAt,
		CreatedBy: preTourismProgram.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": tourismProgram}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedTourismProgram models.TourismProgram
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedTourismProgram); err != nil {
		return err
	}

	return nil
}
func (l *tourismProgramrepo) Delete(ctx context.Context, id string) error {

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
