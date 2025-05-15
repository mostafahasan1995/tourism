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
	do.Provide(i, ourService.NewTourismProgramSvcs)
	do.Provide(i, ourService.NewHotelsSvcs)
	do.Provide(i, ourService.NewPackageSvcs)

	//repos
	do.Provide(i, repo.NewTourismProgramRepo)
	do.Provide(i, repo.NewHotelsRepo)
	do.Provide(i, repo.NewPackageRepo)

	//handlers
	handler.NewTourismProgramHandler(i, r)
	handler.NewHotelsHandler(i, r)
	handler.NewPackageHandler(i, r)

	return i
}
