package di

import (
	"larsa-tourism-microservices/pkg/services/messaging"
	"larsa-tourism-microservices/pkg/services/messaging/handler"
	"larsa-tourism-microservices/pkg/services/messaging/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, messaging.NewMessageSvcs)

	//repos
	do.Provide(i, repo.NewMessageRepo)

	//handlers
	handler.NewMessagingHandler(i, r)

	return i
}
