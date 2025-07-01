package repo

// type DiaryRepo interface {
// 	dbrepo.MainRepo[models.Diary]
// }

// type diaryrepo struct {
// 	dbrepo.MainRepoImpl[models.Diary]
// }

// func NewDiaryRepo(i *do.Injector) (DiaryRepo, error) {
// 	return &diaryrepo{
// 		MainRepoImpl: dbrepo.MainRepoImpl[models.Diary]{
// 			Db:       do.MustInvoke[*mongo.Client](i),
// 			CollName: "tourismDiaries",
// 		},
// 	}, nil
// }
