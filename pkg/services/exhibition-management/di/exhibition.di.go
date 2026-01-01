package di

import (
	exhibition_management "larsa-tourism-microservices/pkg/services/exhibition-management"
	"larsa-tourism-microservices/pkg/services/exhibition-management/handler"
	"larsa-tourism-microservices/pkg/services/exhibition-management/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	// Register repositories
	do.Provide(i, repo.NewExhibitionRepo)

	// Register services
	do.Provide(i, exhibition_management.NewExhibitionSvcs)

	// Register handlers
	handler.NewExhibitionHandler(i, r)

	return i
}
