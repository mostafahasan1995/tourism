package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// FaqGroupRepo interface for FAQ groups
type FaqGroupRepo interface {
	dbrepo.MainRepo[models.FaqGroup]
}

// Repository implementation
type faqGroupRepo struct {
	dbrepo.MainRepoImpl[models.FaqGroup]
}

// Constructor
func NewFaqGroupRepo(i *do.Injector) (FaqGroupRepo, error) {
	return &faqGroupRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqGroup]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqGroups",
		},
	}, nil
}
