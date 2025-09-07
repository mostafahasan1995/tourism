package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CurrencyRepo interface {
	dbrepo.MainRepo[models.Currency]
}

type currencyrepo struct {
	dbrepo.MainRepoImpl[models.Currency]
}

func NewCurrencyRepo(i *do.Injector) (CurrencyRepo, error) {
	return &currencyrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Currency]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiCurrencies",
		},
	}, nil
}
