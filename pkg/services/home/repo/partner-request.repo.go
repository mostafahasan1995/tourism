package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartnerRequestRepo interface {
	dbrepo.MainRepo[models.PartnerRequest]
}

type partnerRequestRepo struct {
	dbrepo.MainRepoImpl[models.PartnerRequest]
}

func NewPartnerRequestRepo(i *do.Injector) (PartnerRequestRepo, error) {
	return &partnerRequestRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.PartnerRequest]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismPartnerRequests",
		},
	}, nil
}
