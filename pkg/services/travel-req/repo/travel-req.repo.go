package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelReqRepo interface {
	dbrepo.MainRepo[models.TravelReq]
}

type travelreq struct {
	dbrepo.MainRepoImpl[models.TravelReq]
}

func NewTravelReqRepo(i *do.Injector) (TravelReqRepo, error) {
	return &travelreq{
		MainRepoImpl: dbrepo.MainRepoImpl[models.TravelReq]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTravelrequests",
		},
	}, nil
}
