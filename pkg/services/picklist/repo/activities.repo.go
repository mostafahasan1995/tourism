package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/picklist/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivitiesRepo interface {
	dbrepo.MainRepo[models.Activities]
}

type activitiesrepo struct {
	dbrepo.MainRepoImpl[models.Activities]
}

func NewActivitiesRepo(i *do.Injector) (ActivitiesRepo, error) {
	return &activitiesrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Activities]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismActivities",
		},
	}, nil
}
