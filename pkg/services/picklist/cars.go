package picklist

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/picklist/filter"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/services/picklist/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Cars, error)
	GetAll(ctx context.Context, filter filter.CarsFilter) (models.CarsPagination, error)
	Add(ctx context.Context, data *models.CarsDto) error
	Update(ctx context.Context, id string, data *models.CarsDto) error
	Delete(ctx context.Context, id string) error
}

type carssvcs struct {
	repo repo.CarsRepo
}

func NewCarsSvcs(i *do.Injector) (CarsSvcs, error) {
	return &carssvcs{
		repo: do.MustInvoke[repo.CarsRepo](i),
	}, nil
}

func (l *carssvcs) GetOne(ctx context.Context, id string) (*models.Cars, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *carssvcs) GetAll(ctx context.Context, filter filter.CarsFilter) (models.CarsPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.CarsPagination{}, err
	}

	return data, nil
}

func (l *carssvcs) Add(ctx context.Context, data *models.CarsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	cars := &models.Cars{
		CarsDto: models.CarsDto{
			CarType: data.CarType,
			Image:   data.Image,
		},
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	if err := l.repo.Add(ctx, cars); err != nil {
		return err
	}

	return nil

}

func (a *carssvcs) Update(ctx context.Context, id string, data *models.CarsDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *carssvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
