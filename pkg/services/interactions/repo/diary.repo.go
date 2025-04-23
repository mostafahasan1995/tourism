package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiaryRepo interface {
	dbrepo.MainRepo[models.Diary]
}

type diaryrepo struct {
	dbrepo.MainRepoImpl[models.Diary]
}

func NewDiaryRepo(i *do.Injector) (DiaryRepo, error) {
	return &diaryrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Diary]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismDiaries",
		},
	}, nil
}
