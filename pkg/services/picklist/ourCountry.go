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

// deprecated
type OurCountrySvcs interface {
	GetOne(ctx context.Context, id string) (*models.OurCountry, error)
	GetAll(ctx context.Context, filter filter.OurCountryFilter) (models.OurCountryPagination, error)
	Add(ctx context.Context, data *models.OurCountryDto) error
	Update(ctx context.Context, id string, data *models.OurCountryDto) error
	Delete(ctx context.Context, id string) error
}

type ourCountrysvcs struct {
	repo repo.OurCountryRepo
}

func NewOurCountrySvcs(i *do.Injector) (OurCountrySvcs, error) {
	return &ourCountrysvcs{
		repo: do.MustInvoke[repo.OurCountryRepo](i),
	}, nil
}

func (l *ourCountrysvcs) GetOne(ctx context.Context, id string) (*models.OurCountry, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *ourCountrysvcs) GetAll(ctx context.Context, filter filter.OurCountryFilter) (models.OurCountryPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.OurCountryPagination{}, err
	}

	return data, nil
}

func (l *ourCountrysvcs) Add(ctx context.Context, data *models.OurCountryDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	ourCountry := &models.OurCountry{
		OurCountryDto: models.OurCountryDto{
			Name:        data.Name,
			Image:       data.Image,
			Icon:        data.Icon,
			Galeres:     data.Galeres,
			Description: data.Description,
		},
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	if err := l.repo.Add(ctx, ourCountry); err != nil {
		return err
	}

	return nil

}

func (a *ourCountrysvcs) Update(ctx context.Context, id string, data *models.OurCountryDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *ourCountrysvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
