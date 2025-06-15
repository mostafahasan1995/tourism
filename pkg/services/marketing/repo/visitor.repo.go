package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/marketing/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitorRepo interface {
	dbrepo.MainRepo[models.Visitor]
}

type visitorRepo struct {
	dbrepo.MainRepoImpl[models.Visitor]
}

func NewVisitorRepo(i *do.Injector) (VisitorRepo, error) {
	return &visitorRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Visitor]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "marketingVisitors",
		},
	}, nil
}

type VisitorActivityRepo interface {
	dbrepo.MainRepo[models.VisitorActivity]
}

type visitorActivityRepo struct {
	dbrepo.MainRepoImpl[models.VisitorActivity]
}

func NewVisitorActivityRepo(i *do.Injector) (VisitorActivityRepo, error) {
	return &visitorActivityRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.VisitorActivity]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "marketingVisitorActivities",
		},
	}, nil
}
