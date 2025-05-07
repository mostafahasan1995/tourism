package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameRepo interface {
	dbrepo.MainRepo[models.Game]
}

type gamerepo struct {
	dbrepo.MainRepoImpl[models.Game]
}

func NewGameRepo(i *do.Injector) (GameRepo, error) {
	return &gamerepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Game]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismGames",
		},
	}, nil
}
