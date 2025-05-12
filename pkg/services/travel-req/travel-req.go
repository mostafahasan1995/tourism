package travelreq

import (
	"context"
	"encoding/json"
	"errors"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/services/travel-req/models"
	"larsa-tourism-microservices/pkg/services/travel-req/repo"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelReqSvcs interface {
	GetRelatedReq(ctx context.Context, id string) (any, error)
	GetAll(ctx context.Context) ([]models.TravelReq, error)
	Add(ctx context.Context, reqType string, data json.RawMessage) (*models.TravelReq, error)
}

type travelreqsvcs struct {
	repo     repo.TravelReqRepo
	reqTypes ReqTypes
	withtxn  *db.WithTxn
}

func NewTravelReqSvcs(i *do.Injector) (TravelReqSvcs, error) {
	return &travelreqsvcs{
		repo:     do.MustInvoke[repo.TravelReqRepo](i),
		reqTypes: do.MustInvoke[ReqTypes](i),
		withtxn:  do.MustInvoke[*db.WithTxn](i),
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
		reqTypeSvcs, ok := t.reqTypes[reqType]
		if !ok {
			return nil, errors.New("invalid request type")
		}

		result, err := reqTypeSvcs.Svcs.Add(ctx, data)
		if err != nil {
			return nil, err
		}

		travelReq := &models.TravelReq{
			Id:          primitive.NewObjectID(),
			ReqId:       "req-0001",
			ServiceType: reqType,
			Date:        time.Now(),
			Status:      "pending",
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
