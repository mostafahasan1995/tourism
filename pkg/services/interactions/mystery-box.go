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

type MysteryBoxSvcs interface {
	Add(ctx context.Context, data *models.MysteryBoxDto) (*models.MysteryBox, error)
	GetMysteryBox(ctx context.Context) (any, error)
	OpenBox(ctx context.Context, mysteryBoxId, boxId string) (*models.Box, error)
}

type mysteryboxsvcs struct {
	repo         repo.MysteryBoxRepo
	boxtrackrepo repo.BoxTrackingRepo
	withtxn      *db.WithTxn
}

func NewMysteryBoxSvcs(i *do.Injector) (MysteryBoxSvcs, error) {
	return &mysteryboxsvcs{
		repo:         do.MustInvoke[repo.MysteryBoxRepo](i),
		boxtrackrepo: do.MustInvoke[repo.BoxTrackingRepo](i),
		withtxn:      do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (m *mysteryboxsvcs) Add(ctx context.Context, data *models.MysteryBoxDto) (*models.MysteryBox, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var boxes []models.Box
	for _, box := range data.Boxes {
		box.Id = primitive.NewObjectID()
		boxes = append(boxes, box)
	}

	data.Boxes = boxes
	if len(data.Boxes) > data.NumOfBox {
		return nil, errors.New("number of boxes cannot be greater than the number of boxes")
	}

	mysBox := &models.MysteryBox{
		Id:            primitive.NewObjectID(),
		MysteryBoxDto: *data,
		CreatedAt:     time.Now(),
		CreatedBy:     cfg.User.Id,
	}

	if err := m.repo.Add(ctx, mysBox); err != nil {
		return nil, err
	}

	return mysBox, nil
}

func (m *mysteryboxsvcs) GetMysteryBox(ctx context.Context) (any, error) {
	mysteryBox, err := m.repo.GetByFilter(ctx, bson.M{"isActive": true, "trash": false})
	if err != nil {
		return nil, errors.New("mystery box not found")
	}

	numOfBox := mysteryBox.NumOfBox //6
	boxes := mysteryBox.Boxes

	boxIds := []primitive.ObjectID{}

	for _, box := range boxes {
		boxIds = append(boxIds, box.Id)
	}
	for range make([]struct{}, numOfBox-len(boxes)) {
		boxIds = append(boxIds, primitive.NewObjectID())
	}

	result := struct {
		Id    primitive.ObjectID   `json:"_id"`
		Boxes []primitive.ObjectID `json:"boxes"`
	}{
		Id:    mysteryBox.Id,
		Boxes: boxIds,
	}

	return result, nil
}

func (m *mysteryboxsvcs) OpenBox(ctx context.Context, mysteryBoxId, boxId string) (*models.Box, error) {
	result, err := m.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		mysId, err := primitive.ObjectIDFromHex(mysteryBoxId)
		if err != nil {
			return nil, err
		}

		bId, err := primitive.ObjectIDFromHex(boxId)
		if err != nil {
			return nil, err
		}

		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}

		attempts, err := m.boxtrackrepo.Count(ctx, bson.M{"userId": cfg.User.Id, "mysteryBoxId": mysId})
		if err != nil {
			return nil, err
		}

		mysteryBox, err := m.repo.GetByFilter(ctx, bson.M{"_id": mysId, "trash": false})
		if err != nil {
			return nil, errors.New("mystery box not found")
		}

		if attempts >= mysteryBox.Attempts {
			return nil, errors.New("you have reached the maximum number of attempts")
		}

		var selectedBox *models.Box
		for _, box := range mysteryBox.Boxes {
			if box.Id == bId {
				selectedBox = &box
				break
			}
		}

		boxTracking := &models.BoxTracking{
			UserId:       cfg.User.Id,
			MysteryBoxId: mysId,
			Box:          selectedBox,
			IsWinner:     selectedBox != nil,
		}

		if err := m.boxtrackrepo.Add(ctx, boxTracking); err != nil {
			return nil, err
		}

		if selectedBox == nil {
			return nil, nil
		}

		return selectedBox, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Box), nil

}
