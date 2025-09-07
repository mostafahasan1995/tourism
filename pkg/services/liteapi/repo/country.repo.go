package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CountryRepo interface {
	dbrepo.MainRepo[models.Country]
}

type countryrepo struct {
	dbrepo.MainRepoImpl[models.Country]
}

func NewCountryRepo(i *do.Injector) (CountryRepo, error) {
	return &countryrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Country]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiCountries",
		},
	}, nil
}
