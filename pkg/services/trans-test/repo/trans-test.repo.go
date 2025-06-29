package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/trans-test/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransTestRepo interface {
	dbrepo.MainRepo[models.TransTest]
}

type transtestrepo struct {
	dbrepo.MainRepoImpl[models.TransTest]
}

func NewTransTestRepo(i *do.Injector) (TransTestRepo, error) {
	return &transtestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.TransTest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTransTests",
		},
	}, nil
}
