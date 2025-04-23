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

	do.Provide(i, repo.NewFaveRepo)
	do.Provide(i, repo.NewTestimonialRepo)

	handler.NewFaveHandler(i, r)
	handler.NewTestimonialHandler(i, r)

	return i
}
