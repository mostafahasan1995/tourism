package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/travel-req/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartnerRequestRepo interface {
	dbrepo.MainRepo[models.PartnerRequest]
}

type partnerRequestrepo struct {
	dbrepo.MainRepoImpl[models.PartnerRequest]
}

func NewPartnerRequestRepo(i *do.Injector) (PartnerRequestRepo, error) {
	return &partnerRequestrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.PartnerRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismPartnerRequest",
		},
	}, nil
}
