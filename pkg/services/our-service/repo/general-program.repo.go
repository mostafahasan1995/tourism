package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type GeneralProgramRepo interface {
	dbrepo.MainRepo[models.GeneralProgram]
}

type generalprogramrepo struct {
	dbrepo.MainRepoImpl[models.GeneralProgram]
}

func NewGeneralProgramRepo(i *do.Injector) (GeneralProgramRepo, error) {
	return &generalprogramrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.GeneralProgram]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismGeneralPrograms",
		},
	}, nil
}
