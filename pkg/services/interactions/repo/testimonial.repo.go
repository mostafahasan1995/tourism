package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestimonialRepo interface {
	dbrepo.MainRepo[models.Testimonial]
}

type testimonial struct {
	dbrepo.MainRepoImpl[models.Testimonial]
}

func NewTestimonialRepo(i *do.Injector) (TestimonialRepo, error) {
	return &testimonial{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Testimonial]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTestimonials",
		},
	}, nil
}
