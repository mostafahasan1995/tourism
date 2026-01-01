package di

import (
	transtest "larsa-tourism-microservices/pkg/services/trans-test"
	"larsa-tourism-microservices/pkg/services/trans-test/handler"
	"larsa-tourism-microservices/pkg/services/trans-test/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, transtest.NewTransTestSvcs)

	//repo
	do.Provide(i, repo.NewTransTestRepo)

	//handler
	handler.NewTransTestHandler(i, r)

	return i
}
