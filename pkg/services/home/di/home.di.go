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

	//do.Provide(i, home.NewOurAgentsSvcs)
	do.Provide(i, home.NewTrustedPartnersSvcs)
	do.Provide(i, home.NewCustomerPersonaSvcs)
	do.Provide(i, home.NewContactUsSvcs)
	do.Provide(i, home.NewPartnerRequestSvcs)
	// Reviews system moved to interactions service - removed from here

	// FAQ services
	do.Provide(i, home.NewFaqPageSvcs)
	do.Provide(i, home.NewFaqGroupSvcs)
	do.Provide(i, home.NewFaqQuestionSvcs)

	//repos

	do.Provide(i, repo.NewContactUsRepo)
	// do.Provide(i, repo.NewOurAgentsRepo)
	do.Provide(i, repo.NewTrustedPartnersRepo)
	do.Provide(i, repo.NewCustomerPersonaRepo)
	do.Provide(i, repo.NewPartnerRequestRepo)
	// Reviews repo moved to interactions service - removed from here

	// FAQ repos
	do.Provide(i, repo.NewFaqPageRepo)
	do.Provide(i, repo.NewFaqGroupRepo)
	do.Provide(i, repo.NewFaqQuestionRepo)

	//handlers

	handler.NewContactUsHandler(i, r)
	//handler.NewOurAgentsHandler(i, r)
	handler.NewTrustedPartnersHandler(i, r)
	handler.NewCustomerPersonaHandler(i, r)
	handler.NewPartnerRequestHandler(i, r)
	// Reviews handler moved to interactions service - removed from here

	// FAQ handler
	handler.NewFaqHandler(i, r)

	// Register repositories
	do.Provide(i, repo.NewSpiritualGiftRepo)

	// Register services
	do.Provide(i, home.NewSpiritualGiftSvcs)

	// Register handlers
	handler.NewSpiritualGiftHandler(i, r)

	return i
}
