package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewsRepo interface {
	dbrepo.MainRepo[models.Review]
}

type reviewsRepo struct {
	dbrepo.MainRepoImpl[models.Review]
}

func NewReviewsRepo(i *do.Injector) (ReviewsRepo, error) {
	return &reviewsRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Review]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "unifiedReviews",
		},
	}, nil
}
