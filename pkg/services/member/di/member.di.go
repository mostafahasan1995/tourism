package di

import (
	"larsa-tourism-microservices/pkg/services/member"
	"larsa-tourism-microservices/pkg/services/member/handler"
	"larsa-tourism-microservices/pkg/services/member/repo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	//services
	do.Provide(i, member.NewAgentSvcs)
	do.Provide(i, member.NewCustomerSvcs)
	do.Provide(i, member.NewMemberAuthSvcs)

	//repos
	do.Provide(i, repo.NewAgentRepo)
	do.Provide(i, repo.NewCustomerRepo)
	do.Provide(i, repo.NewAgentJoinRepo)
	//handlers
	handler.NewAgentHandler(i, r)
	handler.NewCustomerHandler(i, r)

	return i
}
