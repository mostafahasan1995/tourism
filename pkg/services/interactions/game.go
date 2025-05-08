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

type GameSvcs interface {
	GetGame(ctx context.Context) (*models.Game, error)
	UpdateBox(ctx context.Context, boxId string, data models.MysteryBox) (*models.Game, error)
	UpdateSettings(ctx context.Context, data *models.Attempts) (*models.Game, error)
	OpenBox(ctx context.Context, boxId string) (*models.MysteryBox, error)
}

type gamesvcs struct {
	repo             repo.GameRepo
	gamecustomerrepo repo.GameCustomerRepo
}

func NewGameSvcs(i *do.Injector) (GameSvcs, error) {
	return &gamesvcs{
		repo:             do.MustInvoke[repo.GameRepo](i),
		gamecustomerrepo: do.MustInvoke[repo.GameCustomerRepo](i),
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
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 2",
				Discount: models.Discount{
					Value: 0,
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 3",
				Discount: models.Discount{
					Value: 0,
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 4",
				Discount: models.Discount{
					Value: 0,
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 5",
				Discount: models.Discount{
					Value: 0,
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
			{
				Id:   primitive.NewObjectID(),
				Name: "box 6",
				Discount: models.Discount{
					Value: 0,
					Unit:  "%",
				},
				ValidityPeriod: models.Validity{
					Value: 3,
					Unit:  models.ValidityUnitDay,
				},
				ExpiryDate: time.Now().AddDate(0, 1, 0),
			},
		},
		Attempts: models.Attempts{
			Value: 1,
			Unit:  models.AttemptsUnitPerDay,
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

func (g *gamesvcs) OpenBox(ctx context.Context, boxId string) (*models.MysteryBox, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(boxId)
	if err != nil {
		return nil, err
	}

	game, err := g.repo.GetByFilter(ctx, bson.M{"name": "game"})
	if err != nil {
		return nil, errors.New("error get game")
	}

	attempts := game.Attempts

	var numOfAttempts int

	switch attempts.Unit {
	case models.AttemptsUnitPerDay:
		{
			// Get attempts for current day
			now := time.Now().UTC()
			startOfDay := now.Truncate(24 * time.Hour)
			filter := bson.M{
				"userId": cfg.User.Id,
				"date": bson.M{
					"$gte": startOfDay,
				},
			}
			count, err := g.gamecustomerrepo.Count(ctx, filter)
			if err != nil {
				return nil, err
			}
			numOfAttempts = int(count)
		}

	case models.AttemptsUnitPerWeek:
		{
			// Get attempts for current week
			now := time.Now().UTC()
			startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
			startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

			filter := bson.M{
				"userId": cfg.User.Id,
				"date": bson.M{
					"$gte": startOfWeek,
				},
			}
			count, err := g.gamecustomerrepo.Count(ctx, filter)
			if err != nil {
				return nil, err
			}
			numOfAttempts = int(count)
		}

	case models.AttemptsUnitPerMonth:
		{
			// Get attempts for current month
			now := time.Now().UTC()
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

			filter := bson.M{
				"userId": cfg.User.Id,
				"date": bson.M{
					"$gte": startOfMonth,
				},
			}
			count, err := g.gamecustomerrepo.Count(ctx, filter)
			if err != nil {
				return nil, err
			}
			numOfAttempts = int(count)
		}
	}

	// Check if user has exceeded their attempts
	if numOfAttempts >= attempts.Value {
		return nil, errors.New("maximum attempts reached for this period")
	}

	var selectedBox *models.MysteryBox
	for _, box := range game.Boxes {
		if box.Id == _id {
			selectedBox = &box
			break
		}
	}

	if selectedBox == nil {
		return nil, errors.New("box not found")
	}

	rec := &models.GameCustomer{
		Id:     primitive.NewObjectID(),
		UserId: cfg.User.Id,
		Email:  cfg.User.UserData.Email,
		Box:    *selectedBox,
		Date:   time.Now(),
	}

	if err := g.gamecustomerrepo.Add(ctx, rec); err != nil {
		return nil, err
	}

	return selectedBox, nil

}
