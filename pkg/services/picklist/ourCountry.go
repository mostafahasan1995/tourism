package picklist

// deprecated
// type OurCountrySvcs interface {
// 	GetOne(ctx context.Context, id string) (*models.OurCountry, error)
// 	GetAll(ctx context.Context, query any) (*models.OurCountryPagination, error)
// 	Add(ctx context.Context, data *models.OurCountryDto) (*models.OurCountry, error)
// 	Update(ctx context.Context, id string, data *models.OurCountryDto) (*models.OurCountry, error)
// 	Delete(ctx context.Context, id string) error
// }

// type ourCountrySvcs struct {
// 	repo repo.OurCountryRepo
// }

// func NewOurCountrySvcs(i *do.Injector) (OurCountrySvcs, error) {
// 	return &ourCountrySvcs{
// 		repo: do.MustInvoke[repo.OurCountryRepo](i),
// 	}, nil
// }

// func (o *ourCountrySvcs) GetOne(ctx context.Context, id string) (*models.OurCountry, error) {
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, helpers.InvalidObjectId()
// 	}

// 	return o.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
// }

// func (o *ourCountrySvcs) GetAll(ctx context.Context, query any) (*models.OurCountryPagination, error) {
// 	filters, err := helpers.ParseFilters[filter.OurCountryFilter](query)
// 	if err != nil {
// 		return nil, errors.New("invalid query")
// 	}

// 	return o.repo.GetAllPaginated(ctx, *filters)
// }

// func (o *ourCountrySvcs) Add(ctx context.Context, data *models.OurCountryDto) (*models.OurCountry, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ourCountry := &models.OurCountry{
// 		OurCountryDto: *data,
// 		Id:            primitive.NewObjectID(),
// 		Trash:         false,
// 		CreatedAt:     time.Now(),
// 		CreatedBy:     cfg.User.Id,
// 		UpdatedAt:     time.Now(),
// 		UpdatedBy:     cfg.User.Id,
// 	}

// 	if err := o.repo.Add(ctx, ourCountry); err != nil {
// 		return nil, err
// 	}

// 	return ourCountry, nil
// }

// func (o *ourCountrySvcs) Update(ctx context.Context, id string, data *models.OurCountryDto) (*models.OurCountry, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, helpers.InvalidObjectId()
// 	}

// 	ourCountry := &models.OurCountry{
// 		Id:            _id,
// 		OurCountryDto: *data,
// 		UpdatedAt:     time.Now(),
// 		UpdatedBy:     cfg.User.Id,
// 	}

// 	filter := bson.M{"_id": _id}
// 	update := bson.M{"$set": ourCountry}

// 	updatedCountry, err := o.repo.Patch(ctx, filter, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return updatedCountry, nil
// }

// func (o *ourCountrySvcs) Delete(ctx context.Context, id string) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return helpers.InvalidObjectId()
// 	}

// 	filter := bson.M{"_id": _id}
// 	update := bson.M{"$set": bson.M{
// 		"trash":     true,
// 		"updatedAt": time.Now(),
// 		"updatedBy": cfg.User.Id,
// 	}}

// 	_, err = o.repo.Patch(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
