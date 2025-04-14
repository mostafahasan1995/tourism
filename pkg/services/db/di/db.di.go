package di

import (
	"larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/db/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, db.NewSortingSvcs)

	//repos
	do.Provide(i, repo.NewSortingRepo)

	return i
}
