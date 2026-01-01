package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelTypeRepo interface {
	dbrepo.MainRepo[models.HotelType]
}

type hoteltyperepo struct {
	dbrepo.MainRepoImpl[models.HotelType]
}

func NewHotelTypeRepo(i *do.Injector) (HotelTypeRepo, error) {
	return &hoteltyperepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.HotelType]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiHotelTypes",
		},
	}, nil
}
