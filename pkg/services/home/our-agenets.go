package home

// type OurAgentsSvcs interface {
// 	GetOne(ctx context.Context, id string) (*models.OurAgents, error)
// 	GetAll(ctx context.Context, filter filter.OurAgentsFilter) (models.OurAgentsPagination, error)
// 	Add(ctx context.Context, data *models.OurAgentsDto) error
// 	AddMany(ctx context.Context, data []models.OurAgentsDto) error
// 	Update(ctx context.Context, id string, data *models.OurAgentsDto) error
// 	Delete(ctx context.Context, id string) error
// }

// type ourAgentssvcs struct {
// 	repo repo.OurAgentsRepo
// }

// func NewOurAgentsSvcs(i *do.Injector) (OurAgentsSvcs, error) {
// 	return &ourAgentssvcs{
// 		repo: do.MustInvoke[repo.OurAgentsRepo](i),
// 	}, nil
// }

// func (l *ourAgentssvcs) GetOne(ctx context.Context, id string) (*models.OurAgents, error) {
// 	return l.repo.GetOne(ctx, id)

// }

// func (l *ourAgentssvcs) GetAll(ctx context.Context, filter filter.OurAgentsFilter) (models.OurAgentsPagination, error) {

// 	data, err := l.repo.GetAll(ctx, filter)

// 	if err != nil {
// 		return models.OurAgentsPagination{}, err
// 	}

// 	return data, nil
// }

// func (l *ourAgentssvcs) Add(ctx context.Context, data *models.OurAgentsDto) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
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

// 		Id:        primitive.NewObjectID(),
// 		Trash:     false,
// 		CreatedAt: time.Now(),
// 		CreatedBy: cfg.User.Id,
// 		UpdatedAt: time.Now(),
// 		UpdatedBy: cfg.User.Id,
// 	}
// 	if err := l.repo.Add(ctx, ourAgents); err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (l *ourAgentssvcs) AddMany(ctx context.Context, data []models.OurAgentsDto) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	var writeOps []mongo.WriteModel

// 	for _, flr := range data {
// 		ourAgents := &models.OurAgents{
// 			OurAgentsDto: models.OurAgentsDto{
// 				FullName:          flr.FullName,
// 				Email:             flr.Email,
// 				PhoneNumber:       flr.PhoneNumber,
// 				Bio:               flr.Bio,
// 				Nationality:       flr.Nationality,
// 				LanguagesSpoken:   flr.LanguagesSpoken,
// 				CountriesYouServe: flr.CountriesYouServe,
// 				CompanyName:       flr.CompanyName,
// 				CompanyLogo:       flr.CompanyLogo,
// 			},

// 			Id:        primitive.NewObjectID(),
// 			Trash:     false,
// 			CreatedAt: time.Now(),
// 			CreatedBy: cfg.User.Id,
// 			UpdatedAt: time.Now(),
// 			UpdatedBy: cfg.User.Id,
// 		}
// 		writeOp := mongo.NewInsertOneModel()
// 		writeOp.SetDocument(ourAgents)
// 		writeOps = append(writeOps, writeOp)
// 		if len(writeOps) == 0 {
// 			return errors.New("empty write ops")
// 		}
// 	}
// 	_, errInsrt := l.repo.BulkWrite(ctx, writeOps)
// 	if errInsrt != nil {
// 		return err
// 	}

// 	return nil
// }

// func (a *ourAgentssvcs) Update(ctx context.Context, id string, data *models.OurAgentsDto) error {
// 	_id, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return helpers.InvalidObjectId()
// 	}
// 	return a.repo.Update(ctx, _id, data)
// }

// func (a *ourAgentssvcs) Delete(ctx context.Context, id string) error {

// 	return a.repo.Delete(ctx, id)
// }
