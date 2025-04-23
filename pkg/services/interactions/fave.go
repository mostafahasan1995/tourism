package interactions

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FaveSvcs interface {
	Add(ctx context.Context, data *models.FaveDto) (*models.Fave, error)
	GetAllByType(ctx context.Context, t string) ([]*models.Fave, error)
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
		"ref":    data.Ref,
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

func (f *faveSvcs) GetAllByType(ctx context.Context, t string) ([]*models.Fave, error) {
	return nil, nil
}
