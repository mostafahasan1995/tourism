package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type SettingsRepo interface {
	dbrepo.MainRepo[models.Settings]
}

type settingsrepo struct {
	dbrepo.MainRepoImpl[models.Settings]
}

func NewSettingsRepo(i *do.Injector) (SettingsRepo, error) {
	return &settingsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Settings]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismSettings",
		},
	}, nil
}
