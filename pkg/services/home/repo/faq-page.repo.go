package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// FaqPageRepo interface for FAQ pages
type FaqPageRepo interface {
	dbrepo.MainRepo[models.FaqPage]
}

// Repository implementation
type faqPageRepo struct {
	dbrepo.MainRepoImpl[models.FaqPage]
}

// Constructor
func NewFaqPageRepo(i *do.Injector) (FaqPageRepo, error) {
	return &faqPageRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqPage]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqPages",
		},
	}, nil
}
