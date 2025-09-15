package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/liteapi/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelReviewRepo interface {
	dbrepo.MainRepo[models.HotelReview]
}

type hotelreviewrepo struct {
	dbrepo.MainRepoImpl[models.HotelReview]
}

func NewHotelReviewRepo(i *do.Injector) (HotelReviewRepo, error) {
	return &hotelreviewrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.HotelReview]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "liteApiHotelReviews",
		},
	}, nil
}
