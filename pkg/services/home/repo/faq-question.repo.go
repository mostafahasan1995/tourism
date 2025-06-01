package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// FaqQuestionRepo interface for FAQ questions
type FaqQuestionRepo interface {
	dbrepo.MainRepo[models.FaqQuestion]
}

// Repository implementation
type faqQuestionRepo struct {
	dbrepo.MainRepoImpl[models.FaqQuestion]
}

// Constructor
func NewFaqQuestionRepo(i *do.Injector) (FaqQuestionRepo, error) {
	return &faqQuestionRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqQuestion]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqQuestions",
		},
	}, nil
}
