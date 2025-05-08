package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameCustomerRepo interface {
	dbrepo.MainRepo[models.GameCustomer]
}

type gamecustomerrepo struct {
	dbrepo.MainRepoImpl[models.GameCustomer]
}

func NewGameCustomerRepo(i *do.Injector) (GameCustomerRepo, error) {
	return &gamecustomerrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.GameCustomer]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismGameCustomers",
		},
	}, nil
}
