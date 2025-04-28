package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoxTrackingRepo interface {
	dbrepo.MainRepo[models.BoxTracking]
}

type boxTrackingRepo struct {
	dbrepo.MainRepoImpl[models.BoxTracking]
}

func NewBoxTrackingRepo(i *do.Injector) (BoxTrackingRepo, error) {
	return &boxTrackingRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.BoxTracking]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismBoxTrackings",
		},
	}, nil
}
