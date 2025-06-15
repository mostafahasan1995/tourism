package ourservice

import (
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

// func (s *settingssvcs) GetByName(ctx context.Context, name string) (any, error) {
// 	result , err := s.repo.GetByFilter(ctx, bson.M{"name": name})
// 	if err != nil {
// 		return nil,err
// 	}
// 	return result.Value, nil
// }

// func (s *settingssvcs) Update(ctx context.Context, name string) (*models.Settings, error) {

// }
