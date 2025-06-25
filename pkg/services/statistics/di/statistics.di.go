package di

import (
	interactions "larsa-tourism-microservices/pkg/services/interactions"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	picklist "larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/statistics"
	"larsa-tourism-microservices/pkg/services/statistics/handler"
	"larsa-tourism-microservices/pkg/services/statistics/repository"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// Init registers all statistics-related dependencies and handlers
func Init(injector *do.Injector, router *chi.Mux) {
	// Register the manual statistics repository
	do.Provide(injector, func(i *do.Injector) (repository.ManualStatisticsRepository, error) {
		db := do.MustInvoke[*mongo.Database](i)
		return repository.NewManualStatisticsRepository(db), nil
	})

	// Register the statistics service
	do.Provide(injector, func(i *do.Injector) (statistics.StatisticsSvcs, error) {
		programSvcs := do.MustInvoke[ourservice.ProgramSvcs](i)
		travelRequestSvcs := do.MustInvoke[ourservice.TravelRequestSvcs](i)
		hotelsSvcs := do.MustInvoke[ourservice.HotelsSvcs](i)
		packageSvcs := do.MustInvoke[ourservice.PackageSvcs](i)
		destSvcs := do.MustInvoke[picklist.DestinationSvcs](i)
		reviewSvcs := do.MustInvoke[interactions.ReviewsSvcs](i)
		manualStatsRepo := do.MustInvoke[repository.ManualStatisticsRepository](i)

		return statistics.NewStatisticsSvcs(
			programSvcs,
			travelRequestSvcs,
			hotelsSvcs,
			packageSvcs,
			destSvcs,
			reviewSvcs,
			manualStatsRepo,
		), nil
	})

	// Register the statistics handler
	handler.NewStatisticsHandler(injector, router)
}
