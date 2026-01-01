package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelRepo interface {
	dbrepo.MainRepo[models.Hotel]
}

type hotelrepo struct {
	dbrepo.MainRepoImpl[models.Hotel]
}

func NewHotelRepo(i *do.Injector) (HotelRepo, error) {
	return &hotelrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Hotel]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiHotels",
		},
	}, nil
}
