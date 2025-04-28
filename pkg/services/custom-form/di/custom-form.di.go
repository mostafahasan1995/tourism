package di

import (
	customform "larsa-tourism-microservices/pkg/services/custom-form"
	"larsa-tourism-microservices/pkg/services/custom-form/handler"
	"larsa-tourism-microservices/pkg/services/custom-form/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	do.Provide(i, customform.NewCustomFormSvcs)
	do.Provide(i, repo.NewCustomFormRepo)

	handler.NewCustomFormHandler(i, r)

	return i
}
