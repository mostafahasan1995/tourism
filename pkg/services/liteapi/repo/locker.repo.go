package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LockerRepo interface {
	dbrepo.MainRepo[models.Locker]
	EnsureIndexes(ctx context.Context) error
}

type lockerrepo struct {
	dbrepo.MainRepoImpl[models.Locker]
}

func NewLockerRepo(i *do.Injector) (LockerRepo, error) {
	return &lockerrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Locker]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiLocker",
			DbName:   "liteApi",
		},
	}, nil
}

// EnsureIndexes creates a unique compound index on country + language
func (l *lockerrepo) EnsureIndexes(ctx context.Context) error {

	coll := l.Db.Database(l.DbName).Collection(l.CollName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("key_unique"),
	}

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}
