package di

import (
	"larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/handler"
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services

	do.Provide(i, ourService.NewTourismProgramSvcs)
	do.Provide(i, ourService.NewHotelsSvcs)


	//repos

	do.Provide(i, repo.NewTourismProgramRepo)
	do.Provide(i, repo.NewHotelsRepo)


	//handlers

	handler.NewTourismProgramHandler(i, r)
	handler.NewHotelsHandler(i, r)


	return i
}
