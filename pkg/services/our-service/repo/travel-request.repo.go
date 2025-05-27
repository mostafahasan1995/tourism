package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelRequestRepo interface {
	dbrepo.MainRepo[models.TravelRequest]
}

type travelrequestrepo struct {
	dbrepo.MainRepoImpl[models.TravelRequest]
}

func NewTravelRequestRepo(i *do.Injector) (TravelRequestRepo, error) {
	return &travelrequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.TravelRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTravelRequests",
		},
	}, nil
}
