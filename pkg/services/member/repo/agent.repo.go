package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/member/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentRepo interface {
	dbrepo.MainRepo[models.Agent]
}

type agentrepo struct {
	dbrepo.MainRepoImpl[models.Agent]
}

func NewAgentRepo(i *do.Injector) (AgentRepo, error) {
	return &agentrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Agent]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismAgents",
		},
	}, nil
}
