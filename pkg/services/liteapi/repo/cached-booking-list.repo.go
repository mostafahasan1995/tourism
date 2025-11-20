package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CachedBookingListRepo interface {
	dbrepo.MainRepo[models.CachedBookingList]
}

type cachedbookinglistrepo struct {
	dbrepo.MainRepoImpl[models.CachedBookingList]
}

func NewCachedBookingListRepo(i *do.Injector) (CachedBookingListRepo, error) {
	return &cachedbookinglistrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.CachedBookingList]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiCachedBookingLists",
		},
	}, nil
}
