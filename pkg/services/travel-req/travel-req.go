package travelreq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/travel-req/enums"
	"larsa-tourism-microservices/pkg/services/travel-req/models"
	"larsa-tourism-microservices/pkg/services/travel-req/repo"
	"larsa-tourism-microservices/pkg/types"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelReqSvcs interface {
	GetRelatedReq(ctx context.Context, id string) (any, error)
	Get(ctx context.Context, skip, limit int64) (*models.TravelReqWithPagination, error)
	GetAll(ctx context.Context) ([]models.TravelReq, error)
	Add(ctx context.Context, reqType string, data json.RawMessage) (*models.TravelReq, error)
}

type travelreqsvcs struct {
	repo        repo.TravelReqRepo
	reqTypes    ReqTypes
	sortingsvcs dbsvcs.SortingSvcs
	withtxn     *db.WithTxn
}

func NewTravelReqSvcs(i *do.Injector) (TravelReqSvcs, error) {
	return &travelreqsvcs{
		repo:        do.MustInvoke[repo.TravelReqRepo](i),
		reqTypes:    do.MustInvoke[ReqTypes](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (t *travelreqsvcs) GetRelatedReq(ctx context.Context, id string) (any, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	travelReq, err := t.repo.GetByFilter(ctx, bson.M{"_id": _id})
	if err != nil {
		return nil, errors.New("travel request not found")
	}

	relatedReq, err := t.reqTypes[travelReq.ServiceType].Svcs.GetByFilter(ctx, bson.M{"_id": travelReq.Ref})
	if err != nil {
		return nil, err
	}

	return relatedReq, nil
}

func (t *travelreqsvcs) Get(ctx context.Context, skip, limit int64) (*models.TravelReqWithPagination, error) {
	match := bson.M{}

	count, err := t.repo.Count(ctx, match)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}
	var result []models.TravelReq
	errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
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

	return &models.TravelReqWithPagination{
		TravelReqs: result,
		Pagination: pagination,
	}, nil
}

func (t *travelreqsvcs) GetAll(ctx context.Context) ([]models.TravelReq, error) {

	pipeline := []bson.M{
		{"$match": bson.M{}},
	}

	var result []models.TravelReq
	err := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *travelreqsvcs) Add(ctx context.Context, reqType string, data json.RawMessage) (*models.TravelReq, error) {
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {

		reqType := enums.ServiceType(reqType)

		reqTypeSvcs, ok := t.reqTypes[reqType]
		if !ok {
			return nil, errors.New("invalid request type")
		}

		result, err := reqTypeSvcs.Svcs.Add(ctx, data)
		if err != nil {
			return nil, err
		}

		seq, err := t.sortingsvcs.GetAndUpdateSourceSeq(ctx, "travelreq")
		if err != nil {
			return nil, err
		}

		reqId := fmt.Sprintf("RQ-%d-%d", time.Now().Year(), seq)

		travelReq := &models.TravelReq{
			Id:          primitive.NewObjectID(),
			ReqId:       reqId,
			ServiceType: reqType,
			Date:        time.Now(),
			Status:      enums.StatusPending,
			Ref:         result.Id,
		}

		if err := t.repo.Add(ctx, travelReq); err != nil {
			return nil, err
		}

		return travelReq, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TravelReq), nil
}

// func (t *travelreqsvcs) Update(ctx context.Context, id string, data json.RawMessage) (*models.TravelReq, error) {

// 	t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {

// 		return nil, nil
// 	})

// }
