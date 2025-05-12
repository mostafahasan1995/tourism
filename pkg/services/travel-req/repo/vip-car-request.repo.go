package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type VipCarRequestRepo interface {
	dbrepo.MainRepo[models.VipCarRequest]
}

type vipCarRequestrepo struct {
	dbrepo.MainRepoImpl[models.VipCarRequest]
}

func NewVipCarRequestRepo(i *do.Injector) (VipCarRequestRepo, error) {
	return &vipCarRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.VipCarRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismVipCarRequests",
		},
	}, nil
}
