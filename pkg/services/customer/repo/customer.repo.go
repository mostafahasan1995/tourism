package repo

import (
	"larsa-tourism-microservices/pkg/services/customer/models"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepo interface {
	dbrepo.MainRepo[models.Customer]
}

type customerrepo struct {
	dbrepo.MainRepoImpl[models.Customer]
}

func NewCustomerRepo(i *do.Injector) (CustomerRepo, error) {
	return &customerrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Customer]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCustomers",
		},
	}, nil
}
