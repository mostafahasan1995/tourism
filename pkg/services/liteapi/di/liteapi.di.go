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
	// do.Provide(i, liteapi.NewHotelSvcs)
	// do.Provide(i, liteapi.NewReferenceDataSvcs)
	do.Provide(i, liteapi.NewRatesSvcs)
	do.Provide(i, liteapi.NewDataSvcs)
	//do.Provide(i, liteapi.NewSearchV2Svcs)
	do.Provide(i, liteapi.NewSearchv3Svcs)

	//repos
	do.Provide(i, repo.NewHotelRepo)
	do.Provide(i, repo.NewCityRepo)
	do.Provide(i, repo.NewCountryRepo)
	do.Provide(i, repo.NewCurrencyRepo)
	do.Provide(i, repo.NewIataRepo)
	do.Provide(i, repo.NewHotelChainRepo)
	do.Provide(i, repo.NewHotelTypeRepo)
	do.Provide(i, repo.NewHotelDetailsRepo)
	do.Provide(i, repo.NewPreBookRepo)
	do.Provide(i, repo.NewBookingRepo)
	do.Provide(i, repo.NewFacilityRepo)
	do.Provide(i, repo.NewHotelReviewRepo)
	//do.Provide(i, repo.NewLockRepo)
	do.Provide(i, repo.NewDataFetchRepo)
	do.Provide(i, repo.NewLockerRepo)
	do.Provide(i, repo.NewPlaceRepo)

	//handlers
	handler.NewDataHandler(i, r)
	handler.NewRatesHandler(i, r)

	return i
}
