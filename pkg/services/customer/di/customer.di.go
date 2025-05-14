package di

import (
	"larsa-tourism-microservices/pkg/services/customer"
	"larsa-tourism-microservices/pkg/services/customer/handler"
	"larsa-tourism-microservices/pkg/services/customer/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, customer.NewCustomerSvcs)

	//repos
	do.Provide(i, repo.NewCustomerRepo)

	//handlers
	handler.NewCustomerHandler(i, r)

	return i
}
