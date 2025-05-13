package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomPlanRepo interface {
	dbrepo.MainRepo[models.CustomPlan]
}

type customplanrepo struct {
	dbrepo.MainRepoImpl[models.CustomPlan]
}

func NewCustomPlanRepo(i *do.Injector) (CustomPlanRepo, error) {
	return &customplanrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.CustomPlan]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCustomPlanRequests",
		},
	}, nil
}
