package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo interface {
	dbrepo.MainRepo[models.UserBooking]
}

type bookingrepo struct {
	dbrepo.MainRepoImpl[models.UserBooking]
}

func NewBookingRepo(i *do.Injector) (BookingRepo, error) {
	return &bookingrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.UserBooking]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiBookings",
		},
	}, nil
}
