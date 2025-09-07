package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelChainRepo interface {
	dbrepo.MainRepo[models.HotelChain]
}

type hotelchainrepo struct {
	dbrepo.MainRepoImpl[models.HotelChain]
}

func NewHotelChainRepo(i *do.Injector) (HotelChainRepo, error) {
	return &hotelchainrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.HotelChain]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiHotelChains",
		},
	}, nil
}
