package repo

import (
	"larsa-tourism-microservices/pkg/services/interactions/models"

	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type FaveRepo interface {
	dbrepo.MainRepo[models.Fave]
}

type faverepo struct {
	dbrepo.MainRepoImpl[models.Fave]
}

func NewFaveRepo(i *do.Injector) (FaveRepo, error) {
	return &faverepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Fave]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFavorites",
		},
	}, nil
}
