package di

import (
	interactions "larsa-tourism-microservices/pkg/services/interactions"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	picklist "larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/statistics"
	"larsa-tourism-microservices/pkg/services/statistics/handler"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

// Init registers all statistics-related dependencies and handlers
func Init(injector *do.Injector, router *chi.Mux) {
	// Register the statistics service
	do.Provide(injector, func(i *do.Injector) (statistics.StatisticsSvcs, error) {
		programSvcs := do.MustInvoke[ourservice.ProgramSvcs](i)
		travelRequestSvcs := do.MustInvoke[ourservice.TravelRequestSvcs](i)
		hotelsSvcs := do.MustInvoke[ourservice.HotelsSvcs](i)
		packageSvcs := do.MustInvoke[ourservice.PackageSvcs](i)
		destSvcs := do.MustInvoke[picklist.DestinationSvcs](i)
		reviewSvcs := do.MustInvoke[interactions.ReviewsSvcs](i)

		return statistics.NewStatisticsSvcs(
			programSvcs,
			travelRequestSvcs,
			hotelsSvcs,
			packageSvcs,
			destSvcs,
			reviewSvcs,
		), nil
	})

	// Register the statistics handler
	handler.NewStatisticsHandler(injector, router)
}
