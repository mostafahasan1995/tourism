package repo

// type GameCustomerRepo interface {
// 	dbrepo.MainRepo[models.GameCustomer]
// }

// type gamecustomerrepo struct {
// 	dbrepo.MainRepoImpl[models.GameCustomer]
// }

// func NewGameCustomerRepo(i *do.Injector) (GameCustomerRepo, error) {
// 	return &gamecustomerrepo{
// 		MainRepoImpl: dbrepo.MainRepoImpl[models.GameCustomer]{
// 			Db:       do.MustInvoke[*mongo.Client](i),
// 			CollName: "tourismGameCustomers",
// 		},
// 	}, nil
// }
