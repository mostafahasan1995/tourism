package ourservice

import (
	"context"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/samber/do"
)

type SettingsSvcs interface {
}

type settingssvcs struct {
	repo repo.SettingsRepo
}

func NewSettingsSvcs(i *do.Injector) (SettingsSvcs, error) {
	return &settingssvcs{
		repo: do.MustInvoke[repo.SettingsRepo](i),
	}, nil
}

func (s *settingssvcs) Init(ctx context.Context) (*models.Settings, error) {

	// settings := &models.Settings{
	// 	Name:        "settings",
	// 	ProfitRatio: 0,
	// }

	return nil, nil
	// settings := &models.Settings{
	// 	Name:        "settings",
	// 	ProfitRatio: 0,
	// }

	// if err := s.repo.Add(ctx, settings); err != nil {
	// 	return nil, err
	// }
}
