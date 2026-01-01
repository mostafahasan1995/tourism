package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CityRepo interface {
	dbrepo.MainRepo[models.City]
}

type cityrepo struct {
	dbrepo.MainRepoImpl[models.City]
}

func NewCityRepo(i *do.Injector) (CityRepo, error) {
	return &cityrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.City]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiCities",
		},
	}, nil
}
