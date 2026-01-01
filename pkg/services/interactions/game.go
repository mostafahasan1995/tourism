package interactions

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/db"
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
	//OpenBox(ctx context.Context, boxId string) (*models.MysteryBox, error)
	GetCustomers(ctx context.Context) (any, error)
	//
	TryBox(ctx context.Context, data *models.TryBoxDto) (*models.Coupon, error)
}

type gamesvcs struct {
	repo       repo.GameRepo
	couponrepo repo.CouponRepo
	withtxn    *db.WithTxn
}

func NewGameSvcs(i *do.Injector) (GameSvcs, error) {
	return &gamesvcs{
		repo:       do.MustInvoke[repo.GameRepo](i),
		couponrepo: do.MustInvoke[repo.CouponRepo](i),
		withtxn:    do.MustInvoke[*db.WithTxn](i),
	}, nil
}

// game
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

func (g *gamesvcs) TryBox(ctx context.Context, data *models.TryBoxDto) (*models.Coupon, error) {
	result, err := g.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		//check if already won
		_, errW := g.couponrepo.GetByFilter(ctx, bson.M{
			"email": data.Email,
			"won":   true,
		})

		if errW == nil {
			return nil, errors.New("already won")
		} else {
			if !errors.Is(errW, mongo.ErrNoDocuments) {
				return nil, errors.New("error check for prev won")
			}
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
					"email": data.Email,
					"date": bson.M{
						"$gte": startOfDay,
					},
				}
				count, err := g.couponrepo.Count(ctx, filter)
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
					"email": data.Email,
					"date": bson.M{
						"$gte": startOfWeek,
					},
				}
				count, err := g.couponrepo.Count(ctx, filter)
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
					"email": data.Email,
					"date": bson.M{
						"$gte": startOfMonth,
					},
				}
				count, err := g.couponrepo.Count(ctx, filter)
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
			if box.Id == data.BoxId {
				selectedBox = &box
				break
			}
		}

		if selectedBox == nil {
			return nil, errors.New("box not found")
		}

		won := selectedBox.Discount.Value > 0
		var code string
		if won {
			code = util.GenerateUniqueString(10)
		}

		coupon := &models.Coupon{
			Id:       primitive.NewObjectID(),
			Email:    data.Email,
			BoxId:    selectedBox.Id,
			Code:     code,
			Discount: selectedBox.Discount,
			Validity: selectedBox.ValidityPeriod,
			Won:      won,
			Date:     time.Now(),
		}

		if err := g.couponrepo.Add(ctx, coupon); err != nil {
			return nil, err
		}

		return coupon, nil

	})

	if err != nil {
		return nil, err
	}
	return result.(*models.Coupon), nil

}

// player
// func (g *gamesvcs) OpenBox(ctx context.Context, boxId string) (*models.MysteryBox, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//check if already won
// 	_, errW := g.gamecustomerrepo.GetByFilter(ctx, bson.M{
// 		"userId": cfg.User.Id,
// 		"won":    true,
// 	})

// 	if errW == nil {
// 		return nil, errors.New("already won")
// 	} else {
// 		if !errors.Is(errW, mongo.ErrNoDocuments) {
// 			return nil, errors.New("error check for prev won")
// 		}
// 	}

// 	_id, err := primitive.ObjectIDFromHex(boxId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	game, err := g.repo.GetByFilter(ctx, bson.M{"name": "game"})
// 	if err != nil {
// 		return nil, errors.New("error get game")
// 	}

// 	attempts := game.Attempts

// 	var numOfAttempts int

// 	switch attempts.Unit {
// 	case models.AttemptsUnitPerDay:
// 		{
// 			// Get attempts for current day
// 			now := time.Now().UTC()
// 			startOfDay := now.Truncate(24 * time.Hour)
// 			filter := bson.M{
// 				"userId": cfg.User.Id,
// 				"date": bson.M{
// 					"$gte": startOfDay,
// 				},
// 			}
// 			count, err := g.gamecustomerrepo.Count(ctx, filter)
// 			if err != nil {
// 				return nil, err
// 			}
// 			numOfAttempts = int(count)
// 		}

// 	case models.AttemptsUnitPerWeek:
// 		{
// 			// Get attempts for current week
// 			now := time.Now().UTC()
// 			startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
// 			startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

// 			filter := bson.M{
// 				"userId": cfg.User.Id,
// 				"date": bson.M{
// 					"$gte": startOfWeek,
// 				},
// 			}
// 			count, err := g.gamecustomerrepo.Count(ctx, filter)
// 			if err != nil {
// 				return nil, err
// 			}
// 			numOfAttempts = int(count)
// 		}

// 	case models.AttemptsUnitPerMonth:
// 		{
// 			// Get attempts for current month
// 			now := time.Now().UTC()
// 			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

// 			filter := bson.M{
// 				"userId": cfg.User.Id,
// 				"date": bson.M{
// 					"$gte": startOfMonth,
// 				},
// 			}
// 			count, err := g.gamecustomerrepo.Count(ctx, filter)
// 			if err != nil {
// 				return nil, err
// 			}
// 			numOfAttempts = int(count)
// 		}
// 	}

// 	// Check if user has exceeded their attempts
// 	if numOfAttempts >= attempts.Value {
// 		return nil, errors.New("maximum attempts reached for this period")
// 	}

// 	var selectedBox *models.MysteryBox
// 	for _, box := range game.Boxes {
// 		if box.Id == _id {
// 			selectedBox = &box
// 			break
// 		}
// 	}

// 	if selectedBox == nil {
// 		return nil, errors.New("box not found")
// 	}

// 	won := selectedBox.Discount.Value > 0

// 	rec := &models.GameCustomer{
// 		Id:       primitive.NewObjectID(),
// 		UserId:   cfg.User.Id,
// 		Email:    cfg.User.UserData.Email,
// 		BoxId:    selectedBox.Id,
// 		Discount: fmt.Sprintf("%d%s", selectedBox.Discount.Value, selectedBox.Discount.Unit),
// 		Validity: fmt.Sprintf("%d%s", selectedBox.ValidityPeriod.Value, selectedBox.ValidityPeriod.Unit),
// 		Won:      won,
// 		Date:     time.Now(),
// 	}

// 	if err := g.gamecustomerrepo.Add(ctx, rec); err != nil {
// 		return nil, err
// 	}

// 	return selectedBox, nil

// }

func (g *gamesvcs) GetCustomers(ctx context.Context) (any, error) {
	pipeline := []bson.M{
		{"$match": bson.M{}},
		{"$group": bson.M{
			"_id": "$email",
			"attempts": bson.M{
				"$sum": 1,
			},
			"won": bson.M{
				"$sum": bson.M{
					"$cond": bson.M{
						"if":   "$won",
						"then": 1,
						"else": 0,
					},
				},
			},
			"lastDate": bson.M{
				"$max": "$date",
			},
			"winningData": bson.M{
				"$push": bson.M{
					"$cond": bson.M{
						"if": "$won",
						"then": bson.M{
							"discount": "$discount",
							"validity": "$validity",
						},
						"else": "$$REMOVE",
					},
				},
			},
		}},
		{"$addFields": bson.M{
			"discount": bson.M{
				"$arrayElemAt": []interface{}{"$winningData.discount", 0},
			},
			"validity": bson.M{
				"$arrayElemAt": []interface{}{"$winningData.validity", 0},
			},
		}},
		{"$project": bson.M{
			"_id":      1,
			"attempts": 1,
			"won":      1,
			"lastDate": 1,
			"discount": 1,
			"validity": 1,
		}},
	}

	type aux struct {
		Email    string          `bson:"_id" json:"email"`
		Attempts int             `bson:"attempts" json:"attempts"`
		Won      int             `bson:"won" json:"won"`
		LastDate time.Time       `bson:"lastDate" json:"lastDate"`
		Discount models.Discount `bson:"discount" json:"discount"`
		Validity models.Validity `bson:"validity" json:"validity"`
	}

	var result []aux
	err := g.couponrepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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
