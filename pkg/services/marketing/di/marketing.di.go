package di

import (
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/handler"
	"larsa-tourism-microservices/pkg/services/marketing/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, marketing.NewExhibitorProfileSvcs)
	do.Provide(i, marketing.NewVisitorSvcs)
	do.Provide(i, marketing.NewInquirySvcs)
	do.Provide(i, marketing.NewNewsletterSvcs)
	//repos
	do.Provide(i, repo.NewExhibitorProfileRepo)
	do.Provide(i, repo.NewVisitorRepo)
	do.Provide(i, repo.NewVisitorActivityRepo)
	do.Provide(i, repo.NewInquiryRepo)
	do.Provide(i, repo.NewNewsletterRepo)
	//handlers
	handler.NewExhibitorProfileHandler(i, r)
	handler.NewVisitorHandler(i, r)
	handler.NewInquiryHandler(i, r)
	handler.NewNewsletterHandler(i, r)
	return i
}
