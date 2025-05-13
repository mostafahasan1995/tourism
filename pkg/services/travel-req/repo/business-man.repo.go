package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type BusinessManRepo interface {
	dbrepo.MainRepo[models.BusinessMan]
}

type businessmanrepo struct {
	dbrepo.MainRepoImpl[models.BusinessMan]
}

func NewBusinessManRepo(i *do.Injector) (BusinessManRepo, error) {
	return &businessmanrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.BusinessMan]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismBusinessmenRequests",
		},
	}, nil
}
