package di

import (
	"larsa-tourism-microservices/pkg/services/request"
	"larsa-tourism-microservices/pkg/services/request/handler"
	"larsa-tourism-microservices/pkg/services/request/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services

	do.Provide(i, request.NewFlightTicketRequestSvcs)
	do.Provide(i, request.NewVipCarRequestSvcs)
	do.Provide(i, request.NewPartnerRequestSvcs)
	do.Provide(i, request.NewTravelReqSvcs)

	//repos

	do.Provide(i, repo.NewFlightTicketRequestRepo)
	do.Provide(i, repo.NewVipCarRequestRepo)
	do.Provide(i, repo.NewPartnerRequestRepo)
	do.Provide(i, repo.NewTravelReqRepo)

	//handlers

	handler.NewFlightTicketRequestHandler(i, r)
	handler.NewVipCarRequestHandler(i, r)
	handler.NewPartnerRequestHandler(i, r)
	handler.NewTravelReqHandler(i, r)

	return i
}
