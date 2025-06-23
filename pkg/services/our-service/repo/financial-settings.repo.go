package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type FinancialSettingsRepo interface {
	dbrepo.MainRepo[models.FinancialSettings]
}

type financialsettingsrepo struct {
	dbrepo.MainRepoImpl[models.FinancialSettings]
}

func NewFinancialSettingsRepo(i *do.Injector) (FinancialSettingsRepo, error) {
	return &financialsettingsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FinancialSettings]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFinancialSettings",
		},
	}, nil
}
