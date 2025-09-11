package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo interface {
	dbrepo.MainRepo[models.Booking]
}

type bookingrepo struct {
	dbrepo.MainRepoImpl[models.Booking]
}

func NewBookingRepo(i *do.Injector) (BookingRepo, error) {
	return &bookingrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Booking]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiBookings",
		},
	}, nil
}
