package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlightTicketRequestRepo interface {
	dbrepo.MainRepo[models.FlightTicketRequest]
}

type flightTicketRequestrepo struct {
	dbrepo.MainRepoImpl[models.FlightTicketRequest]
}

func NewFlightTicketRequestRepo(i *do.Injector) (FlightTicketRequestRepo, error) {
	return &flightTicketRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FlightTicketRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFlightTicketRequest",
		},
	}, nil
}
