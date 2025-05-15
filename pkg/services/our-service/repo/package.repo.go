package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackageRepo interface {
	dbrepo.MainRepo[models.Package]
}

type packagerepo struct {
	dbrepo.MainRepoImpl[models.Package]
}

func NewPackageRepo(i *do.Injector) (PackageRepo, error) {
	return &packagerepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Package]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismPackages",
		},
	}, nil
}
