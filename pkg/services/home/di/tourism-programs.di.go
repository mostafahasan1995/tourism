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

	do.Provide(i, home.NewTourismProgramSvcs)
	do.Provide(i, home.NewHotelsSvcs)
	do.Provide(i, home.NewFlightTicketRequestSvcs)
	do.Provide(i, home.NewVipCarRequestSvcs)

	//repos

	do.Provide(i, repo.NewTourismProgramRepo)
	do.Provide(i, repo.NewHotelsRepo)
	do.Provide(i, repo.NewFlightTicketRequestRepo)
	do.Provide(i, repo.NewVipCarRequestRepo)

	//handlers

	handler.NewTourismProgramHandler(i, r)
	handler.NewHotelsHandler(i, r)
	handler.NewFlightTicketRequestHandler(i, r)
	handler.NewVipCarRequestHandler(i, r)

	return i
}
