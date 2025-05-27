package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProgramRepo interface {
	dbrepo.MainRepo[models.Program]
}

type programrepo struct {
	dbrepo.MainRepoImpl[models.Program]
}

func NewProgramRepo(i *do.Injector) (ProgramRepo, error) {
	return &programrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Program]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismPrograms",
		},
	}, nil
}
