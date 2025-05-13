package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type DelegationRepo interface {
	dbrepo.MainRepo[models.Delegation]
}

type delegationrepo struct {
	dbrepo.MainRepoImpl[models.Delegation]
}

func NewDelegationRepo(i *do.Injector) (DelegationRepo, error) {
	return &delegationrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Delegation]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismDelegationRequests",
		},
	}, nil
}
