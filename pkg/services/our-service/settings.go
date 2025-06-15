package ourservice

import (
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/samber/do"
)

type SettingsSvcs interface {
	// Init(ctx context.Context) error
	// GetSettingByName(ctx context.Context, name string) (*models.Settings, error)
	// Update(ctx context.Context, name string, value any) (*models.Settings, error)
}

type settingssvcs struct {
	repo repo.SettingsRepo
}

func NewSettingsSvcs(i *do.Injector) (SettingsSvcs, error) {
	return &settingssvcs{
		repo: do.MustInvoke[repo.SettingsRepo](i),
	}, nil
}

// func (s *settingssvcs) Init(ctx context.Context) error {
// 	allSettings := []any{
// 		models.Settings{
// 			Name:  "profitRatio",
// 			Value: 0,
// 		},
// 		// add more settings here
// 	}

// 	if err := s.repo.AddMany(ctx, allSettings); err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (s *settingssvcs) GetSettingByName(ctx context.Context, name string) (*models.Settings, error) {
// 	return s.repo.GetByFilter(ctx, bson.M{"name": name})
// }

// func (s *settingssvcs) Update(ctx context.Context, name string, value any) (*models.Settings, error) {

// 	filter := bson.M{"name": name}
// 	update := bson.M{"$set": bson.M{"value": value}}

// 	updatedSetting, err := s.repo.Patch(ctx, filter, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return updatedSetting, nil
// }
