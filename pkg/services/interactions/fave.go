package interactions

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FaveSvcs interface {
	Fav(ctx context.Context, data *models.FaveDto) (string, error)
	GetAllByType(ctx context.Context, faveType models.FaveType) ([]models.FaveItem, error)
}

type faveSvcs struct {
	repo repo.FaveRepo
}

func NewFaveSvcs(i *do.Injector) (FaveSvcs, error) {
	return &faveSvcs{
		repo: do.MustInvoke[repo.FaveRepo](i),
	}, nil
}

func (s *faveSvcs) Fav(ctx context.Context, data *models.FaveDto) (string, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return "", err
	}
	existingFave, err := s.repo.GetByFilter(ctx, bson.M{"userId": cfg.User.Id, "refId": data.RefId, "type": data.Type})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	if existingFave != nil {
		err := s.repo.DeleteMain(ctx, bson.M{"_id": existingFave.Id})
		if err != nil {
			return "", err
		}
		return "Removed From Favorites", nil
	}

	data.IsFav = true
	data.Type = models.FaveType(data.Type)
	fave := &models.Fave{
		Id:        primitive.NewObjectID(),
		UserId:    cfg.User.Id,
		FaveDto:   *data,
		CreatedAt: time.Now(),
	}

	err = s.repo.Add(ctx, fave)

	if err != nil {
		return "", err
	}

	return "Added To Favorites", nil
}

func (s *faveSvcs) GetAllByType(ctx context.Context, faveType models.FaveType) ([]models.FaveItem, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	}

	// Determine the collection name based on the favorite type
	var collectionName string
	switch faveType {
	case models.FaveTypeProgram:
		collectionName = "tourismPrograms"
	case models.FaveTypeHotel:
		collectionName = "tourismHotels"
	case models.FaveTypeTravelerStory:
		collectionName = "tourismTravelerStories"
	case models.FaveTypeClientStory:
		collectionName = "tourismClientStories"
	case models.FaveTypeExhibition:
		collectionName = "tourismExhibitions"
	case models.FaveTypeAgent:
		collectionName = "tourismAgents"
	default:
		return []models.FaveItem{}, nil
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":   faveType,
				"isFav":  true,
				"userId": userId,
				"trash":  bson.M{"$ne": true},
			},
		},
		{
			"$lookup": bson.M{
				"from":         collectionName,
				"localField":   "refId",
				"foreignField": "_id",
				"as":           "item",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$item",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	var result []models.FaveItem
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
