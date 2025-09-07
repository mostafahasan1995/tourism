package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type IataRepo interface {
	dbrepo.MainRepo[models.Iata]
}

type iatarepo struct {
	dbrepo.MainRepoImpl[models.Iata]
}

func NewIataRepo(i *do.Injector) (IataRepo, error) {
	return &iatarepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Iata]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiIatas",
		},
	}, nil
}
