package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/marketing/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsletterRepo interface {
	dbrepo.MainRepo[models.Newsletter]
}

type newsletterRepo struct {
	dbrepo.MainRepoImpl[models.Newsletter]
}

func NewNewsletterRepo(i *do.Injector) (NewsletterRepo, error) {
	return &newsletterRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Newsletter]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismNewsletters",
		},
	}, nil
}
