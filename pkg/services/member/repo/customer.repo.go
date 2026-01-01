package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/member/models"

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
