package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SettingsSvcs interface {
	Get(ctx context.Context) (*models.Settings, error)
	Update(ctx context.Context, data *models.SettingsDto) (*models.Settings, error)
}

type settingssvcs struct {
	repo repo.SettingsRepo
}

func NewSettingsSvcs(i *do.Injector) (SettingsSvcs, error) {
	return &settingssvcs{
		repo: do.MustInvoke[repo.SettingsRepo](i),
	}, nil
}

func (s *settingssvcs) init(ctx context.Context) (*models.Settings, error) {
	settings := &models.Settings{
		Id:   primitive.NewObjectID(),
		Name: "settings",
		SettingsDto: models.SettingsDto{
			ProfitRatio: 0,
		},
	}

	if err := s.repo.Add(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *settingssvcs) Get(ctx context.Context) (*models.Settings, error) {
	settings, err := s.repo.GetByFilter(ctx, bson.M{"name": "settings"})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s.init(ctx)
		}
		return nil, err
	}

	return settings, nil
}

func (s *settingssvcs) Update(ctx context.Context, data *models.SettingsDto) (*models.Settings, error) {
	settings, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}

	settings.SettingsDto = *data

	filter := bson.M{"_id": settings.Id}
	update := bson.M{"$set": settings}

	updatedSettings, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedSettings, nil

}
