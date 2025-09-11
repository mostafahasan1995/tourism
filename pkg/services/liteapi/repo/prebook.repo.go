package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type PreBookRepo interface {
	dbrepo.MainRepo[models.PreBook]
}

type prebookrepo struct {
	dbrepo.MainRepoImpl[models.PreBook]
}

func NewPreBookRepo(i *do.Injector) (PreBookRepo, error) {
	return &prebookrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.PreBook]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiPreBooks",
		},
	}, nil
}
