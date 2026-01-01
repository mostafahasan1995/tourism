package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentFinancialRepo interface {
	dbrepo.MainRepo[models.AgentFinancialAccount]
}

type agentfinancialrepo struct {
	dbrepo.MainRepoImpl[models.AgentFinancialAccount]
}

func NewAgentFinancialRepo(i *do.Injector) (AgentFinancialRepo, error) {
	return &agentfinancialrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.AgentFinancialAccount]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismAgentFinancials",
		},
	}, nil
}
