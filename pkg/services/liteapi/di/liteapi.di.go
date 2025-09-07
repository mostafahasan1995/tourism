package di

import (
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/services/liteapi/handler"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, liteapi.NewHotelSvcs)
	do.Provide(i, liteapi.NewReferenceDataSvcs)

	//repos
	do.Provide(i, repo.NewHotelRepo)
	do.Provide(i, repo.NewCityRepo)
	do.Provide(i, repo.NewCountryRepo)
	do.Provide(i, repo.NewCurrencyRepo)
	do.Provide(i, repo.NewIataRepo)
	do.Provide(i, repo.NewHotelChainRepo)
	do.Provide(i, repo.NewHotelTypeRepo)

	//handlers
	handler.NewDataHandler(i, r)

	return i
}
