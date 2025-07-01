package di

import (
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/handler"
	"larsa-tourism-microservices/pkg/services/interactions/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	do.Provide(i, interactions.NewFaveSvcs)
	do.Provide(i, interactions.NewTestimonialSvcs)
	//do.Provide(i, interactions.NewDiarySvcs) // deprecated

	do.Provide(i, interactions.NewGameSvcs)
	do.Provide(i, interactions.NewTravelExperSvcs)

	do.Provide(i, interactions.NewReviewsSvcs)

	do.Provide(i, repo.NewFaveRepo)
	do.Provide(i, repo.NewTestimonialRepo)
	//do.Provide(i, repo.NewDiaryRepo) //deprecated

	do.Provide(i, repo.NewGameRepo)
	//do.Provide(i, repo.NewGameCustomerRepo)
	do.Provide(i, repo.NewCouponRepo)
	do.Provide(i, repo.NewTravelerStoryRepo)
	do.Provide(i, repo.NewClientStoryRepo)

	do.Provide(i, repo.NewReviewsRepo)

	handler.NewFaveHandler(i, r)
	handler.NewTestimonialHandler(i, r)
	//handler.NewDiaryHandler(i, r) //deprecated

	handler.NewGameHandler(i, r)
	handler.NewTravelExperHandler(i, r)

	handler.NewReviewsHandler(i, r)

	return i
}
