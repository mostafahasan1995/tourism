package customform

import (
	"context"
	"larsa-tourism-microservices/pkg/services/custom-form/models"
	"larsa-tourism-microservices/pkg/services/custom-form/repo"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomFormSvcs interface {
	Add(ctx context.Context, data *models.CustomFormDto) (*models.CustomForm, error)
}

type customFormSvcs struct {
	repo repo.CustomFormRepo
}

func NewCustomFormSvcs(i *do.Injector) (CustomFormSvcs, error) {
	return &customFormSvcs{
		repo: do.MustInvoke[repo.CustomFormRepo](i),
	}, nil
}

func (s *customFormSvcs) Add(ctx context.Context, data *models.CustomFormDto) (*models.CustomForm, error) {
	form := &models.CustomForm{
		Id:            primitive.NewObjectID(),
		CustomFormDto: *data,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Add(ctx, form); err != nil {
		return nil, err
	}

	return form, nil
}
