package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerPersonaRepo interface {
	dbrepo.MainRepo[models.CustomerPersona]
}

type customerPersonaRepo struct {
	dbrepo.MainRepoImpl[models.CustomerPersona]
}

func NewCustomerPersonaRepo(i *do.Injector) (CustomerPersonaRepo, error) {
	return &customerPersonaRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.CustomerPersona]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCustomerPersonas",
		},
	}, nil
}
