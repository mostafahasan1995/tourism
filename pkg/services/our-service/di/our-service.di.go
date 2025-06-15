package di

import (
	ourService "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/handler"
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, ourService.NewHotelsSvcs)
	do.Provide(i, ourService.NewPackageSvcs)
	do.Provide(i, ourService.NewTravelRequestSvcs)
	do.Provide(i, ourService.NewProgramSvcs)
	do.Provide(i, ourService.NewInvoiceSvcs)

	//repos
	//do.Provide(i, repo.NewTourismProgramRepo)
	do.Provide(i, repo.NewHotelsRepo)
	do.Provide(i, repo.NewPackageRepo)
	do.Provide(i, repo.NewTravelRequestRepo)
	do.Provide(i, repo.NewProgramRepo)
	do.Provide(i, repo.NewInvoiceRepo)

	//handlers
	handler.NewHotelsHandler(i, r)
	handler.NewPackageHandler(i, r)
	handler.NewTravelRequestHandler(i, r)
	handler.NewProgramHandler(i, r)
	handler.NewInvoiceHandler(i, r)

	return i
}
