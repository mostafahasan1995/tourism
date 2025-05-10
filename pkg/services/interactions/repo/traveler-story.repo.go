package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelerStoryRepo interface {
	dbrepo.MainRepo[models.TravelerStory]
}

type travelerStory struct {
	dbrepo.MainRepoImpl[models.TravelerStory]
}

func NewTravelerStoryRepo(i *do.Injector) (TravelerStoryRepo, error) {
	return &travelerStory{
		MainRepoImpl: dbrepo.MainRepoImpl[models.TravelerStory]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTravelerStories",
		},
	}, nil
}
