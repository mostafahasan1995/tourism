package  home

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TourismProgramSvcs interface {
	GetOne(ctx context.Context, id string) (*models.TourismProgram, error)
	GetAll(ctx context.Context,filter filter.TourismProgramFilter) (models.TourismProgramPagination, error)
	Add(ctx context.Context, data *models.TourismProgramDto) error
	Update(ctx context.Context, id string, data *models.TourismProgramDto) error
	Delete(ctx context.Context, id string) error
}

type tourismProgramsvcs struct {
	repo repo.TourismProgramRepo
}

func NewTourismProgramSvcs(i *do.Injector) (TourismProgramSvcs, error) {
	return &tourismProgramsvcs{
		repo: do.MustInvoke[repo.TourismProgramRepo](i),
	}, nil
}

func (l *tourismProgramsvcs) GetOne(ctx context.Context, id string) (*models.TourismProgram, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *tourismProgramsvcs) GetAll(ctx context.Context,filter filter.TourismProgramFilter) (models.TourismProgramPagination, error) {

	data, err := l.repo.GetAll(ctx,filter)
	if err != nil {
		return models.TourismProgramPagination{}, err
	}
	return data, nil
}

func (l *tourismProgramsvcs) Add(ctx context.Context, data *models.TourismProgramDto) error {

	if err := l.repo.Add(ctx, data); err != nil {
		return err
	}

	return nil

}

func (a *tourismProgramsvcs) Update(ctx context.Context, id string, data *models.TourismProgramDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *tourismProgramsvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
