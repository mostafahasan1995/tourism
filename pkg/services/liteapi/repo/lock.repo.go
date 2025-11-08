package repo

// type LockRepo interface {
// 	dbrepo.MainRepo[models.Lock]
// 	EnsureIndexes(ctx context.Context) error
// }

// type lockrepo struct {
// 	dbrepo.MainRepoImpl[models.Lock]
// }

// func NewLockRepo(i *do.Injector) (LockRepo, error) {
// 	return &lockrepo{
// 		MainRepoImpl: dbrepo.MainRepoImpl[models.Lock]{
// 			Db:       do.MustInvoke[*mongo.Client](i),
// 			CollName: "liteApiLocks",
// 		},
// 	}, nil
// }

// // EnsureIndexes creates a unique compound index on country + language
// func (l *lockrepo) EnsureIndexes(ctx context.Context) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	coll := l.Db.Database(cfg.Db).Collection(l.CollName)

// 	indexModel := mongo.IndexModel{
// 		Keys:    bson.D{{Key: "placeId", Value: 1}, {Key: "language", Value: 1}},
// 		Options: options.Index().SetUnique(true).SetName("palceId_language_unique"),
// 	}

// 	_, err = coll.Indexes().CreateOne(ctx, indexModel)
// 	if err != nil {
// 		if !mongo.IsDuplicateKeyError(err) {
// 			return err
// 		}
// 	}

// 	return nil
// }
