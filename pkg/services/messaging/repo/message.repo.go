package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/messaging/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepo interface {
	dbrepo.MainRepo[models.Message]
}

type messagerepo struct {
	dbrepo.MainRepoImpl[models.Message]
}

func NewMessageRepo(i *do.Injector) (MessageRepo, error) {
	return &messagerepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Message]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismMessages",
		},
	}, nil
}
