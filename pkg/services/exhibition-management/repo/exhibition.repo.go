package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/exhibition-management/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExhibitionRepo interface {
	dbrepo.MainRepo[models.Exhibition]
}

type exhibitionRepo struct {
	dbrepo.MainRepoImpl[models.Exhibition]
}

func NewExhibitionRepo(i *do.Injector) (ExhibitionRepo, error) {
	return &exhibitionRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Exhibition]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "exhibitions",
		},
	}, nil
}
