package repo

// type OurCountryRepo interface {
// 	dbrepo.MainRepo[models.OurCountry]
// 	GetOne(ctx context.Context, id string) (*models.OurCountry, error)
// 	GetAll(ctx context.Context, match bson.M, skip, limit int64) ([]models.OurCountry, int64, error)
// 	GetAllPaginated(ctx context.Context, f filter.OurCountryFilter) (*models.OurCountryPagination, error)
// 	Update(ctx context.Context, id primitive.ObjectID, data *models.OurCountryDto) (*models.OurCountry, error)
// 	Delete(ctx context.Context, id string) error
// }

// type ourCountryrepo struct {
// 	dbrepo.MainRepoImpl[models.OurCountry]
// 	db       *mongo.Client
// 	collName string
// }

// func NewOurCountryRepo(i *do.Injector) (OurCountryRepo, error) {
// 	return &ourCountryrepo{
// 		MainRepoImpl: dbrepo.MainRepoImpl[models.OurCountry]{
// 			Db:       do.MustInvoke[*mongo.Client](i),
// 			CollName: "tourismOurCountry",
// 		},
// 		db:       do.MustInvoke[*mongo.Client](i),
// 		collName: "tourismOurCountry",
// 	}, nil
// }

// func (l *ourCountryrepo) GetOne(ctx context.Context, id string) (*models.OurCountry, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	var data models.OurCountry
// 	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (l *ourCountryrepo) GetAll(ctx context.Context, match bson.M, skip, limit int64) ([]models.OurCountry, int64, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	totalCount, err := coll.CountDocuments(ctx, match)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"_id": -1})
// 	cur, err := coll.Find(ctx, match, findOptions)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer cur.Close(ctx)

// 	var countries []models.OurCountry
// 	if err := cur.All(ctx, &countries); err != nil {
// 		return nil, 0, err
// 	}

// 	return countries, totalCount, nil
// }

// func (l *ourCountryrepo) GetAllPaginated(ctx context.Context, f filter.OurCountryFilter) (*models.OurCountryPagination, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	// Build filter using pipeline pattern for consistency
// 	pipeline := f.BuildPipeline(bson.M{})
// 	match := pipeline[0]["$match"].(bson.M)

// 	// Count total documents matching the filter
// 	totalCount, err := coll.CountDocuments(ctx, match)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Pagination defaults and limits
// 	page := f.Page
// 	if page <= 0 {
// 		page = 1
// 	}
// 	size := f.Size
// 	if size <= 0 {
// 		size = 10 // Default page size
// 	}
// 	skip := int64((page - 1) * size)
// 	limit := int64(size)

// 	// Query options with pagination
// 	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"_id": -1})

// 	cur, err := coll.Find(ctx, match, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cur.Close(ctx)

// 	var countries []models.OurCountry
// 	if err := cur.All(ctx, &countries); err != nil {
// 		return nil, err
// 	}

// 	// Calculate total pages
// 	totalPages := float64(0)
// 	if size > 0 {
// 		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
// 	}

// 	result := &models.OurCountryPagination{
// 		OurCountry: countries,
// 		Pagination: common.Pagination{
// 			TotalPages: totalPages,
// 			PerPage:    int64(size),
// 			TotalCount: totalCount,
// 		},
// 	}

// 	return result, nil
// }

// func (l *ourCountryrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.OurCountryDto) (*models.OurCountry, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	preOurCountry, err := l.GetOne(ctx, id.Hex())
// 	if err != nil {
// 		return nil, err
// 	}

// 	ourCountry := &models.OurCountry{
// 		OurCountryDto: models.OurCountryDto{
// 			Name:        data.Name,
// 			Image:       data.Image,
// 			Icon:        data.Icon,
// 			Galeres:     data.Galeres,
// 			Description: data.Description,
// 		},
// 		Id:        id,
// 		Trash:     false,
// 		CreatedAt: preOurCountry.CreatedAt,
// 		CreatedBy: preOurCountry.CreatedBy,
// 		UpdatedBy: cfg.User.Id,
// 		UpdatedAt: time.Now(),
// 	}

// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": ourCountry}

// 	upsert := false
// 	after := options.After
// 	opts := &options.FindOneAndUpdateOptions{
// 		ReturnDocument: &after,
// 		Upsert:         &upsert,
// 	}

// 	coll := l.db.Database(cfg.Db).Collection(l.collName)
// 	var updatedOurCountry models.OurCountry
// 	if err := coll.FindOneAndUpdate(
// 		ctx,
// 		filter,
// 		update,
// 		opts,
// 	).Decode(&updatedOurCountry); err != nil {
// 		return nil, err
// 	}

// 	return &updatedOurCountry, nil
// }

// func (l *ourCountryrepo) Delete(ctx context.Context, id string) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	filter := bson.M{"_id": _id}
// 	update := bson.M{"$set": bson.M{
// 		"trash":     true,
// 		"updatedAt": time.Now(),
// 		"updatedBy": cfg.User.Id,
// 	}}

// 	upsert := false
// 	after := options.After
// 	opt := options.FindOneAndUpdateOptions{
// 		ReturnDocument: &after,
// 		Upsert:         &upsert,
// 	}
// 	result := coll.FindOneAndUpdate(
// 		ctx,
// 		filter,
// 		update,
// 		&opt,
// 	)

// 	if result.Err() != nil {
// 		return result.Err()
// 	}

// 	return nil
// }
