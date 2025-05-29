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

// FaqGroupRepo interface for FAQ groups
type FaqGroupRepo interface {
	dbrepo.MainRepo[models.FaqGroup]
}

// FaqQuestionRepo interface for FAQ questions
type FaqQuestionRepo interface {
	dbrepo.MainRepo[models.FaqQuestion]
}

// Repository implementations
type faqPageRepo struct {
	dbrepo.MainRepoImpl[models.FaqPage]
}

type faqGroupRepo struct {
	dbrepo.MainRepoImpl[models.FaqGroup]
}

type faqQuestionRepo struct {
	dbrepo.MainRepoImpl[models.FaqQuestion]
}

// Constructors
func NewFaqPageRepo(i *do.Injector) (FaqPageRepo, error) {
	return &faqPageRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqPage]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqPages",
		},
	}, nil
}

func NewFaqGroupRepo(i *do.Injector) (FaqGroupRepo, error) {
	return &faqGroupRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqGroup]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqGroups",
		},
	}, nil
}

func NewFaqQuestionRepo(i *do.Injector) (FaqQuestionRepo, error) {
	return &faqQuestionRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.FaqQuestion]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismFaqQuestions",
		},
	}, nil
}
