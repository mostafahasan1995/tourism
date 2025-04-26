package interactions

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FaveSvcs interface {
	Add(ctx context.Context, data *models.FaveDto) (*models.Fave, error)
	GetAllByType(ctx context.Context, t models.FaveType) ([]models.FaveItem, error)
}

type faveSvcs struct {
	repo repo.FaveRepo
}

func NewFaveSvcs(i *do.Injector) (FaveSvcs, error) {
	return &faveSvcs{
		repo: do.MustInvoke[repo.FaveRepo](i),
	}, nil
}

func (f *faveSvcs) Add(ctx context.Context, data *models.FaveDto) (*models.Fave, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	fav := &models.Fave{
		UserId:    cfg.User.Id,
		FaveDto:   *data,
		CreatedAt: time.Now(),
	}

	filter := bson.M{
		"userId": cfg.User.Id,
		"refId":  data.RefId,
	}

	update := bson.M{"$set": fav}

	upsert := true
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	fav, err = f.repo.Patch(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return fav, nil
}

func (f *faveSvcs) GetAllByType(ctx context.Context, typ models.FaveType) ([]models.FaveItem, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":   typ,
				"isFav":  true,
				"userId": cfg.User.Id,
			},
		},
	}

	var collName string

	switch typ {
	case models.FaveTypeProgram:
		{
			collName = "tourismTourismProgram"
		}
	case models.FaveTypeHotel:
		{
			collName = "tourismHotels"
		}
	case models.FaveTypeExhibition:
		{
			collName = "tourismExhibition"
		}
	case models.FaveTypeDiary:
		{
			collName = "tourismDiaries"
		}
	}

	itemPipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": collName,
				"let": bson.M{
					"itemId": "$refId",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []any{
									bson.M{"$eq": []any{"$_id", "$$itemId"}},
									bson.M{"$eq": []any{"$trash", false}},
								},
							},
						},
					},
				},
				"as": "item",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$item",
				"preserveNullAndEmptyArrays": false,
			},
		},
	}

	pipeline = append(pipeline, itemPipeline...)

	var result []models.FaveItem
	f.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})

	return result, nil

}
