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

	do.Provide(i, home.NewFlightTicketRequestSvcs)


	//repos


	do.Provide(i, repo.NewFlightTicketRequestRepo)


	//handlers

	handler.NewFlightTicketRequestHandler(i, r)


	return i
}
