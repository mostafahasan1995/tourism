package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/member/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentJoinRepo interface {
	dbrepo.MainRepo[models.AgentJoin]
}

type agentjoinrepo struct {
	dbrepo.MainRepoImpl[models.AgentJoin]
}

func NewAgentJoinRepo(i *do.Injector) (AgentJoinRepo, error) {
	return &agentjoinrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.AgentJoin]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismAgentjoins",
		},
	}, nil
}
