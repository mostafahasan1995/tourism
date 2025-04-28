package repo

import (
	"larsa-tourism-microservices/pkg/services/custom-form/models"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomFormRepo interface {
	dbrepo.MainRepo[models.CustomForm]
}

type customFormRepo struct {
	dbrepo.MainRepoImpl[models.CustomForm]
}

func NewCustomFormRepo(i *do.Injector) (CustomFormRepo, error) {
	return &customFormRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.CustomForm]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCustomForms",
		},
	}, nil
}
