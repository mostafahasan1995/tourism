package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlaceRepo interface {
	dbrepo.MainRepo[models.Place]
}

type placerepo struct {
	dbrepo.MainRepoImpl[models.Place]
}

func NewPlaceRepo(i *do.Injector) (PlaceRepo, error) {
	return &placerepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Place]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiPlaces",
		},
	}, nil
}
