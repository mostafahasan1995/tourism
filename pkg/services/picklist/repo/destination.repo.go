package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/picklist/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationRepo interface {
	dbrepo.MainRepo[models.Destination]
}

type destinationrepo struct {
	dbrepo.MainRepoImpl[models.Destination]
}

func NewDestinationRepo(i *do.Injector) (DestinationRepo, error) {
	return &destinationrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Destination]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismDestinations",
		},
	}, nil
}
