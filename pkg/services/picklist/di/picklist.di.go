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
	do.Provide(i, picklist.NewOurCountrySvcs) // deprecated
	do.Provide(i, picklist.NewDestinationSvcs)

	//repos
	do.Provide(i, repo.NewCarsRepo)
	do.Provide(i, repo.NewOurCountryRepo)
	do.Provide(i, repo.NewDestinationRepo)

	//handlers

	handler.NewCarsHandler(i, r)
	handler.NewOurCountryHandler(i, r)
	handler.NewDestinationHandler(i, r)

	return i
}
