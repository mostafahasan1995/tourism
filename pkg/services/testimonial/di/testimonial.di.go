package di

import (
	"larsa-tourism-microservices/pkg/services/testimonial"
	"larsa-tourism-microservices/pkg/services/testimonial/handler"
	"larsa-tourism-microservices/pkg/services/testimonial/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {

	do.Provide(i, testimonial.NewTestimonialService)
	do.Provide(i, repo.NewTestimonialRepo)

	handler.NewTestimonialHandler(i, r)

	return i
}
