package di

import (
	"larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/picklist/handler"
	"larsa-tourism-microservices/pkg/services/picklist/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services

	do.Provide(i, picklist.NewCarsSvcs)

	//repos
	do.Provide(i, repo.NewCarsRepo)

	//handlers

	handler.NewCarsHandler(i, r)

	return i
}
