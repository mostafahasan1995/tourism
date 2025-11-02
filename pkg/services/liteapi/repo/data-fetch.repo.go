package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataFetchRepo interface {
	dbrepo.MainRepo[models.DataFetch]
	EnsureIndexes(ctx context.Context) error
}

type datafetchrepo struct {
	dbrepo.MainRepoImpl[models.DataFetch]
}

func NewDataFetchRepo(i *do.Injector) (DataFetchRepo, error) {
	return &datafetchrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.DataFetch]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiDataFetches",
		},
	}, nil
}

func (d *datafetchrepo) EnsureIndexes(ctx context.Context) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := d.Db.Database(cfg.Db).Collection(d.CollName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "placeId", Value: 1}, {Key: "language", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("placeId_language_unique"),
	}

	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}
