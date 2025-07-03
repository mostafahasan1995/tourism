package interactions

// type DiarySvcs interface {
// 	GetOne(ctx context.Context, id string) (*models.Diary, error)
// 	GetAll(ctx context.Context, query string) ([]models.Diary, error)
// 	Get(ctx context.Context, skip, limit int64, query string) (*models.DiaryWithPagination, error)
// 	Add(ctx context.Context, data *models.DiaryDto) (*models.Diary, error)
// }

// type diarysvcs struct {
// 	repo repo.DiaryRepo
// }

// func NewDiarySvcs(i *do.Injector) (DiarySvcs, error) {
// 	return &diarysvcs{
// 		repo: do.MustInvoke[repo.DiaryRepo](i),
// 	}, nil
// }

// func (d *diarysvcs) GetOne(ctx context.Context, id string) (*models.Diary, error) {
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return d.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
// }

// func (d *diarysvcs) GetAll(ctx context.Context, query string) ([]models.Diary, error) {
// 	match := bson.M{"trash": false}

// 	f, err := helpers.ParseFilters[filter.DiaryFilter](query)
// 	if err != nil {
// 		return nil, errors.New("invalid query")
// 	}

// 	pipeline := f.BuildPipeline(match)

// 	var result []models.Diary
// 	err = d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
// 		if err := cur.All(ctx, &result); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (d *diarysvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.DiaryWithPagination, error) {
// 	match := bson.M{"trash": false}

// 	f, err := helpers.ParseFilters[filter.DiaryFilter](query)
// 	if err != nil {
// 		return nil, errors.New("invalid query")
// 	}

// 	pipeline := f.BuildPipeline(match)

// 	countPipeline := make([]bson.M, len(pipeline))
// 	copy(countPipeline, pipeline)

// 	count, err := d.repo.Count(ctx, countPipeline)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
// 	pipeline = append(pipeline, bson.M{"$skip": skip})
// 	pipeline = append(pipeline, bson.M{"$limit": limit})

// 	var result []models.Diary
// 	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
// 		if err := cur.All(ctx, &result); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if errAg != nil {
// 		return nil, errAg
// 	}

// 	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
// 	pagination := types.Pagination{
// 		TotalPages: totalPages,
// 		PerPage:    limit,
// 		TotalCount: count,
// 	}

// 	return &models.DiaryWithPagination{
// 		Diaries:    result,
// 		Pagination: pagination,
// 	}, nil
// }

// func (d *diarysvcs) Add(ctx context.Context, data *models.DiaryDto) (*models.Diary, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	diary := &models.Diary{
// 		Id:        primitive.NewObjectID(),
// 		DiaryDto:  *data,
// 		CreatedAt: time.Now(),
// 		CreatedBy: cfg.User.Id,
// 	}

// 	if err := d.repo.Add(ctx, diary); err != nil {
// 		return nil, err
// 	}

// 	return diary, nil
// }
