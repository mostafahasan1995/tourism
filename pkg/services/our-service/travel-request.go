package ourservice

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelRequestSvcs interface {
	Get(ctx context.Context, skip, limit int64, query string) (*models.TravelRequestPagination, error)
	GetOne(ctx context.Context, id string) (*models.TravelRequest, error)
	Add(ctx context.Context, data *models.TravelRequestDto) (*models.TravelRequest, error)
	Update(ctx context.Context, id string, data *models.TravelRequestDto) (*models.TravelRequest, error)
	MyRequests(ctx context.Context, status string) ([]models.TravelRequest, error)
	UpdateStatus(ctx context.Context, id string, status string) (*models.TravelRequest, error)
	Patch(ctx context.Context, filter, update bson.M) (*models.TravelRequest, error)
	BulkWrite(ctx context.Context, writes []mongo.WriteModel) (*mongo.BulkWriteResult, error)
}

type travelrequestsvcs struct {
	repo        repo.TravelRequestRepo
	sortingsvcs dbsvcs.SortingSvcs
	withtxn     *db.WithTxn
}

func NewTravelRequestSvcs(i *do.Injector) (TravelRequestSvcs, error) {
	return &travelrequestsvcs{
		repo:        do.MustInvoke[repo.TravelRequestRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (t *travelrequestsvcs) GetOne(ctx context.Context, id string) (*models.TravelRequest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return t.repo.GetByFilter(ctx, bson.M{"_id": _id})
}

func (t *travelrequestsvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.TravelRequestPagination, error) {
	match := bson.M{}

	filters, err := filter.NewTravelReqFilters(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.TravelRequest
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

	return &models.TravelRequestPagination{
		Requests:   result,
		Pagination: pagination,
	}, nil
}

func (t *travelrequestsvcs) Add(ctx context.Context, data *models.TravelRequestDto) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {

		seq, err := t.sortingsvcs.GetAndUpdateSourceSeq(ctx, "travelRequest")
		if err != nil {
			return nil, err
		}

		reqId := fmt.Sprintf("RQ-%d-%d", time.Now().Year(), seq)

		request := &models.TravelRequest{
			Id:               primitive.NewObjectID(),
			ReqId:            reqId,
			TravelRequestDto: *data,
			Date:             time.Now(),
			Status:           enums.TravelReqStatusPending,
			CustomerId:       cfg.User.Id,
			CreatedAt:        time.Now(),
			CreatedBy:        cfg.User.Id,
		}

		if err := t.repo.Add(ctx, request); err != nil {
			return nil, err
		}

		return request, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TravelRequest), nil
}

func (t *travelrequestsvcs) Update(ctx context.Context, id string, data *models.TravelRequestDto) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	request := &models.TravelRequest{
		Id:               _id,
		TravelRequestDto: *data,
		UpdatedAt:        time.Now(),
		UpdatedBy:        cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": request}

	updatedRequest, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

func (t *travelrequestsvcs) MyRequests(ctx context.Context, status string) ([]models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	match := bson.M{"customerId": cfg.User.Id}

	if status != "all" {
		match["status"] = status
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
	}

	var requests []models.TravelRequest
	err = t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &requests)
	})

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (t *travelrequestsvcs) UpdateStatus(ctx context.Context, id string, status string) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	request, err := t.repo.GetByFilter(ctx, bson.M{"_id": id, "trash": false})
	if err != nil {
		return nil, errors.New("travel request not found")
	}

	if request.CustomerId != cfg.User.Id {
		return nil, errors.New("unauthorized")
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now(), "updatedBy": cfg.User.Id}}

	updatedRequest, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedRequest, nil
}

func (t *travelrequestsvcs) Patch(ctx context.Context, filter, update bson.M) (*models.TravelRequest, error) {
	return t.repo.Patch(ctx, filter, update)
}

func (t *travelrequestsvcs) BulkWrite(ctx context.Context, writes []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return t.repo.BulkWrite(ctx, writes)
}
