package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactUsRepo interface {
	dbrepo.MainRepo[models.ContactUs]
}

type contactUsrepo struct {
	dbrepo.MainRepoImpl[models.ContactUs]
}

func NewContactUsRepo(i *do.Injector) (ContactUsRepo, error) {
	return &contactUsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.ContactUs]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismContactUs",
		},
	}, nil
}

type ContactUsSettingsRepo interface {
	dbrepo.MainRepo[models.ContactUsSettings]
}

type contactUsSettingsRepo struct {
	dbrepo.MainRepoImpl[models.ContactUsSettings]
}

func NewContactUsSettingsRepo(i *do.Injector) (ContactUsSettingsRepo, error) {
	return &contactUsSettingsRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.ContactUsSettings]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismContactUsSettings",
		},
	}, nil
}
