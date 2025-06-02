package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/filter"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewsRepo interface {
	dbrepo.MainRepo[models.Review]
	GetByFilter(ctx context.Context, filter bson.M) (*models.Review, error)
	GetAllByFilter(ctx context.Context, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error)
	GetEntityReviews(ctx context.Context, entityType string, refId string, page, perPage int) (models.ReviewPagination, error)
	UpdateReviewStatus(ctx context.Context, reviewId string, status string) error
}

type reviewsRepo struct {
	dbrepo.MainRepoImpl[models.Review]
	db       *mongo.Client
	collName string
}

func NewReviewsRepo(i *do.Injector) (ReviewsRepo, error) {
	return &reviewsRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Review]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "unifiedReviews",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "unifiedReviews",
	}, nil
}

func (r *reviewsRepo) GetByFilter(ctx context.Context, filter bson.M) (*models.Review, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := r.db.Database(cfg.Db).Collection(r.collName)
	var review models.Review
	err = coll.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewsRepo) GetAllByFilter(ctx context.Context, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.ReviewPagination{}, err
	}

	coll := r.db.Database(cfg.Db).Collection(r.collName)
	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.ReviewPagination{}, err
	}

	// Validate pagination values
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	// Query options with pagination and sorting
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{"date", -1}, {"createdAt", -1}}) // Sort by review date, newest first (ordered)

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.ReviewPagination{}, err
	}

	var reviews []models.Review
	if err := cur.All(ctx, &reviews); err != nil {
		return models.ReviewPagination{}, err
	}

	// Prepare pagination result
	totalPages := int64(0)
	if perPage > 0 {
		totalPages = (totalCount + int64(perPage) - 1) / int64(perPage)
	}

	result := models.ReviewPagination{
		Reviews: reviews,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(perPage),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (r *reviewsRepo) GetEntityReviews(ctx context.Context, entityType string, refId string, page, perPage int) (models.ReviewPagination, error) {
	refObjectId, err := primitive.ObjectIDFromHex(refId)
	if err != nil {
		return models.ReviewPagination{}, err
	}

	filter := filter.ReviewsFilter{
		Type: entityType,
		Ref:  refObjectId,
	}

	return r.GetAllByFilter(ctx, filter, page, perPage)
}

func (r *reviewsRepo) UpdateReviewStatus(ctx context.Context, reviewId string, status string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	reviewObjectId, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return err
	}

	coll := r.db.Database(cfg.Db).Collection(r.collName)
	filter := bson.M{"_id": reviewObjectId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err = coll.UpdateOne(ctx, filter, update)
	return err
}
