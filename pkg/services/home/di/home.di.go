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
	do.Provide(i, home.NewOurAgentsSvcs)

	//repos

	do.Provide(i, repo.NewContactUsRepo)
	do.Provide(i, repo.NewOurAgentsRepo)

	//handlers

	handler.NewContactUsHandler(i, r)
	handler.NewOurAgentsHandler(i, r)

	return i
}
