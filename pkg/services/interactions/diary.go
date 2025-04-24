package interactions

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiarySvcs interface {
	GetOne(ctx context.Context, id string) (*models.Diary, error)
	GetAll(ctx context.Context) ([]models.Diary, error)
	Add(ctx context.Context, data *models.DiaryDto) (*models.Diary, error)
}

type diarysvcs struct {
	repo repo.DiaryRepo
}

func NewDiarySvcs(i *do.Injector) (DiarySvcs, error) {
	return &diarysvcs{
		repo: do.MustInvoke[repo.DiaryRepo](i),
	}, nil
}

func (d *diarysvcs) GetOne(ctx context.Context, id string) (*models.Diary, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return d.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (d *diarysvcs) GetAll(ctx context.Context) ([]models.Diary, error) {
	match := bson.M{
		"$match": bson.M{"trash": false},
	}

	pipeline := []bson.M{
		match,
	}

	var result []models.Diary
	err := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

func (d *diarysvcs) Add(ctx context.Context, data *models.DiaryDto) (*models.Diary, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	diary := &models.Diary{
		Id:        primitive.NewObjectID(),
		DiaryDto:  *data,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
	}

	if err := d.repo.Add(ctx, diary); err != nil {
		return nil, err
	}

	return diary, nil
}
