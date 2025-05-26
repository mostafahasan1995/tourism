package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomProgramRepo interface {
	dbrepo.MainRepo[models.CustomProgram]
}

type customprogramrepo struct {
	dbrepo.MainRepoImpl[models.CustomProgram]
}

func NewCustomProgramRepo(i *do.Injector) (CustomProgramRepo, error) {
	return &customprogramrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.CustomProgram]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCustomPrograms",
		},
	}, nil
}
