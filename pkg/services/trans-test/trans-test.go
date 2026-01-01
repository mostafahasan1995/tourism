package transtest

import (
	"context"
	"encoding/json"
	"errors"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/trans-test/models"
	"larsa-tourism-microservices/pkg/services/trans-test/repo"
	"larsa-tourism-microservices/pkg/types"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransTestSvcs interface {
	Add(ctx context.Context, transTest *models.TransTestDto) (*models.TransTest, error)
	GetOne(ctx context.Context, id string) (*models.TransTest, error)
	Get(ctx context.Context, skip, limit int64, queryStr string) (*models.TransTestPagination, error)
}

type transtestsvcs struct {
	repo    repo.TransTestRepo
	withtxn *db.WithTxn
}

func NewTransTestSvcs(i *do.Injector) (TransTestSvcs, error) {
	return &transtestsvcs{
		repo:    do.MustInvoke[repo.TransTestRepo](i),
		withtxn: do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (t *transtestsvcs) Add(ctx context.Context, transTest *models.TransTestDto) (*models.TransTest, error) {
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		transTest := &models.TransTest{
			Id:           primitive.NewObjectID(),
			TransTestDto: *transTest,
			Status:       "active",
			Trash:        false,
			CreatedAt:    time.Now(),
			//CreatedBy:    "system",
		}

		if err := t.repo.Add(ctx, transTest); err != nil {
			return nil, err
		}

		return transTest, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TransTest), nil
}

func (t *transtestsvcs) GetOne(ctx context.Context, id string) (*models.TransTest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	result, err := t.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *transtestsvcs) Get(ctx context.Context, skip, limit int64, queryStr string) (*models.TransTestPagination, error) {

	var conditions query.Conditions

	if err := json.Unmarshal([]byte(queryStr), &conditions); err != nil {
		return nil, err
	}

	if err := conditions.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := conditions.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.TransTest
	errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.TransTestPagination{
		Trans:      result,
		Pagination: pagination,
	}, nil
}
