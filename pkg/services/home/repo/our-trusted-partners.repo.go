package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"

	"go.mongodb.org/mongo-driver/mongo"
)

type TrustedPartnersRepo interface {
	dbrepo.MainRepo[models.TrustedPartner]
}

type trustedPartnersRepo struct {
	dbrepo.MainRepoImpl[models.TrustedPartner]

	db       *mongo.Client
	collName string
}

func NewTrustedPartnersRepo(i *do.Injector) (TrustedPartnersRepo, error) {
	return &trustedPartnersRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.TrustedPartner]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismTrustedPartners",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismTrustedPartners",
	}, nil
}
