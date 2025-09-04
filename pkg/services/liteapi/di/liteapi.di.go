package di

import (
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/services/liteapi/handler"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, liteapi.NewHotelSvcs)

	//repos
	do.Provide(i, repo.NewHotelRepo)

	//handlers
	handler.NewHotelHandler(i, r)

	return i
}
