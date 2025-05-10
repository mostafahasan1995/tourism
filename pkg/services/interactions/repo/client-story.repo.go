package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientStoryRepo interface {
	dbrepo.MainRepo[models.ClientStory]
}

type clientStory struct {
	dbrepo.MainRepoImpl[models.ClientStory]
}

func NewClientStoryRepo(i *do.Injector) (ClientStoryRepo, error) {
	return &clientStory{
		MainRepoImpl: dbrepo.MainRepoImpl[models.ClientStory]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismClientStories",
		},
	}, nil
}
