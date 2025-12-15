package di

import (
	"larsa-tourism-microservices/pkg/services/exchange"
	"larsa-tourism-microservices/pkg/services/exchange/handler"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Register(i *do.Injector) *do.Injector {
	do.Provide(i, exchange.NewExchangeService)
	return i
}

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	// Register service (already done via Register, but keeping for consistency)
	do.Provide(i, exchange.NewExchangeService)

	// Register handler
	handler.NewExchangeHandler(i, r)

	return i
}
