package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacilityRepo interface {
	dbrepo.MainRepo[models.Facility]
}

type facilityrepo struct {
	dbrepo.MainRepoImpl[models.Facility]
}

func NewFacilityRepo(i *do.Injector) (FacilityRepo, error) {
	return &facilityrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Facility]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiFacilities",
		},
	}, nil
}
