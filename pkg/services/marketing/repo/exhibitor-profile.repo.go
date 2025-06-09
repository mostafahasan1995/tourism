package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/marketing/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExhibitorProfileRepo interface {
	dbrepo.MainRepo[models.ExhibitorProfile]
}

type exhibitorProfileRepo struct {
	dbrepo.MainRepoImpl[models.ExhibitorProfile]
}

func NewExhibitorProfileRepo(i *do.Injector) (ExhibitorProfileRepo, error) {
	return &exhibitorProfileRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.ExhibitorProfile]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "marketingExhibitorProfiles",
		},
	}, nil
}
