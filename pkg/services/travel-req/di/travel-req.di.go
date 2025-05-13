package di

import (
	travelreq "larsa-tourism-microservices/pkg/services/travel-req"
	"larsa-tourism-microservices/pkg/services/travel-req/handler"
	"larsa-tourism-microservices/pkg/services/travel-req/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	do.Provide(i, travelreq.NewTravelReqSvcs)
	do.Provide(i, travelreq.NewReqTypes)
	//
	do.Provide(i, travelreq.NewBusinessManSvcs)
	do.Provide(i, travelreq.NewCustomPlanSvcs)
	do.Provide(i, travelreq.NewDelegationSvcs)
	do.Provide(i, travelreq.NewFlightTicketRequestSvcs)
	do.Provide(i, travelreq.NewPartnerRequestSvcs)
	do.Provide(i, travelreq.NewVipCarRequestSvcs)

	do.Provide(i, repo.NewTravelReqRepo)
	//
	do.Provide(i, repo.NewBusinessManRepo)
	do.Provide(i, repo.NewCustomPlanRepo)
	do.Provide(i, repo.NewDelegationRepo)
	do.Provide(i, repo.NewFlightTicketRequestRepo)
	do.Provide(i, repo.NewPartnerRequestRepo)
	do.Provide(i, repo.NewVipCarRequestRepo)

	handler.NewTravelReqHandler(i, r)

	return i
}
