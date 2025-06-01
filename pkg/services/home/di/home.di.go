package di

import (
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/handler"
	"larsa-tourism-microservices/pkg/services/home/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services

	do.Provide(i, home.NewContactUsSvcs)
	//do.Provide(i, home.NewOurAgentsSvcs)
	do.Provide(i, home.NewReviewsSvcs)

	// FAQ services
	do.Provide(i, home.NewFaqPageSvcs)
	do.Provide(i, home.NewFaqGroupSvcs)
	do.Provide(i, home.NewFaqQuestionSvcs)

	//repos

	do.Provide(i, repo.NewContactUsRepo)
	do.Provide(i, repo.NewOurAgentsRepo)
	do.Provide(i, repo.NewReviewsRepo)

	// FAQ repos
	do.Provide(i, repo.NewFaqPageRepo)
	do.Provide(i, repo.NewFaqGroupRepo)
	do.Provide(i, repo.NewFaqQuestionRepo)

	//handlers

	handler.NewContactUsHandler(i, r)
	//handler.NewOurAgentsHandler(i, r)
	handler.NewReviewsHandler(i, r)

	// FAQ handler
	handler.NewFaqHandler(i, r)

	return i
}
