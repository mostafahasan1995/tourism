package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type MysteryBoxRepo interface {
	dbrepo.MainRepo[models.MysteryBox]
}

type mysteryBoxRepo struct {
	dbrepo.MainRepoImpl[models.MysteryBox]
}

func NewMysteryBoxRepo(i *do.Injector) (MysteryBoxRepo, error) {
	return &mysteryBoxRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.MysteryBox]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismMysteryBoxes",
		},
	}, nil
}
