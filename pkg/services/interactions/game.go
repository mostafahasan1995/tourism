package interactions

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameSvcs interface {
	GetGame(ctx context.Context) (*models.Game, error)
	UpdateBox(ctx context.Context, boxId string, data models.MysteryBox) (*models.Game, error)
	UpdateSettings(ctx context.Context, data *models.Attempts) (*models.Game, error)
}

type gamesvcs struct {
	repo repo.GameRepo
}

func NewGameSvcs(i *do.Injector) (GameSvcs, error) {
	return &gamesvcs{
		repo: do.MustInvoke[repo.GameRepo](i),
	}, nil
}

func (g *gamesvcs) Init(ctx context.Context) (*models.Game, error) {
	game := &models.Game{
		Id:   primitive.NewObjectID(),
		Name: "game",
		Boxes: []models.MysteryBox{
			{
				Id:   primitive.NewObjectID(),
				Name: "box 1",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 2",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 3",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 4",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 5",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 6",
				Discount: models.Discount{
					Value: 0,
					Unit:  "percent",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  "day",
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
		},
		Attempts: models.Attempts{
			Value: 1,
			Unit:  "day",
		},
	}

	err := g.repo.Add(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *gamesvcs) GetGame(ctx context.Context) (*models.Game, error) {

	game, err := g.repo.GetByFilter(ctx, bson.M{"name": "game"})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			game, err := g.Init(ctx)
			if err != nil {
				return nil, err
			}
			return game, nil

		} else {
			return nil, err
		}

	}

	return game, nil

}

func (g *gamesvcs) UpdateBox(ctx context.Context, boxId string, data models.MysteryBox) (*models.Game, error) {
	_id, err := primitive.ObjectIDFromHex(boxId)
	if err != nil {
		return nil, err
	}
	data.Id = _id
	filter := bson.M{"name": "game", "boxes._id": _id}
	update := bson.M{"$set": bson.M{"boxes.$": data}}

	updatedGame, err := g.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil

}

func (g *gamesvcs) UpdateSettings(ctx context.Context, data *models.Attempts) (*models.Game, error) {

	filter := bson.M{"name": "game"}
	update := bson.M{"$set": bson.M{"attempts": data}}

	updatedGame, err := g.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil

}
