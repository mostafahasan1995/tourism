package repo

// deprecated
// type OurAgentsRepo interface {
// 	dbrepo.MainRepo[models.OurAgents]
// 	GetOne(ctx context.Context, id string) (*models.OurAgents, error)
// 	GetAll(ctx context.Context, filter filter.OurAgentsFilter) (models.OurAgentsPagination, error)
// 	Update(ctx context.Context, id primitive.ObjectID, data *models.OurAgentsDto) error
// 	Delete(ctx context.Context, id string) error
// }

// type ourAgentsrepo struct {
// 	dbrepo.MainRepoImpl[models.OurAgents]

// 	db       *mongo.Client
// 	collName string
// }

// func NewOurAgentsRepo(i *do.Injector) (OurAgentsRepo, error) {
// 	return &ourAgentsrepo{
// 		MainRepoImpl: dbrepo.MainRepoImpl[models.OurAgents]{
// 			Db:       do.MustInvoke[*mongo.Client](i),
// 			CollName: "tourismOurAgents",
// 		},
// 		db:       do.MustInvoke[*mongo.Client](i),
// 		collName: "tourismOurAgents",
// 	}, nil
// }

// func (l *ourAgentsrepo) GetOne(ctx context.Context, id string) (*models.OurAgents, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	var data models.OurAgents
// 	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
// 		return nil, err
// 	}
// 	return &data, nil

// }

// func (l *ourAgentsrepo) GetAll(ctx context.Context, filter filter.OurAgentsFilter) (models.OurAgentsPagination, error) {

// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return models.OurAgentsPagination{}, err
// 	}

// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	filterBody := filter.ToBsonFilter()

// 	// Count total documents matching the filter
// 	totalCount, err := coll.CountDocuments(ctx, filterBody)
// 	if err != nil {
// 		return models.OurAgentsPagination{}, err
// 	}

// 	// Pagination defaults and limits
// 	page := filter.Page
// 	if page <= 0 {
// 		page = 1
// 	}
// 	size := filter.Size
// 	if size <= 0 {
// 		size = int(totalCount) // return all if invalid
// 	}
// 	skip := int64((page - 1) * size)
// 	limit := int64(size)

// 	// Query options with pagination
// 	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

// 	cur, err := coll.Find(ctx, filterBody, findOptions)
// 	if err != nil {
// 		return models.OurAgentsPagination{}, err
// 	}

// 	var programs []models.OurAgents
// 	if err := cur.All(ctx, &programs); err != nil {
// 		return models.OurAgentsPagination{}, err
// 	}

// 	// Prepare pagination result
// 	totalPages := float64(0)
// 	if size > 0 {
// 		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
// 	}

// 	result := models.OurAgentsPagination{
// 		OurAgents: programs,
// 		Pagination: common.Pagination{
// 			TotalPages: totalPages,
// 			PerPage:    int64(size),
// 			TotalCount: totalCount,
// 		},
// 	}

// 	return result, nil
// }

// func (l *ourAgentsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.OurAgentsDto) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	coll := l.db.Database(cfg.Db).Collection(l.collName)

// 	preOurAgents, err := l.GetOne(ctx, id.Hex())
// 	if err != nil {
// 		return err
// 	}

// 	ourAgents := &models.OurAgents{
// 		OurAgentsDto: models.OurAgentsDto{
// 			FullName:          data.FullName,
// 			Email:             data.Email,
// 			PhoneNumber:       data.PhoneNumber,
// 			Bio:               data.Bio,
// 			Nationality:       data.Nationality,
// 			LanguagesSpoken:   data.LanguagesSpoken,
// 			CountriesYouServe: data.CountriesYouServe,
// 			CompanyName:       data.CompanyName,
// 			CompanyLogo:       data.CompanyLogo,
// 		},

// 		Id:        id,
// 		Trash:     false,
// 		CreatedAt: preOurAgents.CreatedAt,
// 		CreatedBy: preOurAgents.CreatedBy,
// 		UpdatedBy: cfg.User.Id,
// 		UpdatedAt: time.Now(),
// 	}

// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": ourAgents}

// 	upsert := false
// 	after := options.After
// 	opts := &options.FindOneAndUpdateOptions{
// 		ReturnDocument: &after,
// 		Upsert:         &upsert,
// 	}

// 	var updatedOurAgents models.OurAgents
// 	if err := coll.FindOneAndUpdate(
// 		ctx,
// 		filter,
// 		update,
// 		opts,
// 	).Decode(&updatedOurAgents); err != nil {
// 		return err
// 	}

// 	return nil
// }
// func (l *ourAgentsrepo) Delete(ctx context.Context, id string) error {

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
// 		"updatedBy": primitive.NilObjectID,
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
