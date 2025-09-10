package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelDetailsRepo interface {
	dbrepo.MainRepo[models.HotelDetails]
}

type hoteldetailsrepo struct {
	dbrepo.MainRepoImpl[models.HotelDetails]
}

func NewHotelDetailsRepo(i *do.Injector) (HotelDetailsRepo, error) {
	return &hoteldetailsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.HotelDetails]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiHotelDetails",
		},
	}, nil
}
