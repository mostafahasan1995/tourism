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

type FinancialSettingsSvcs interface {
	Get(ctx context.Context) (*models.FinancialSettings, error)
	Update(ctx context.Context, data *models.FinancialSettingsDto) (*models.FinancialSettings, error)
}

type financialsettingssvcs struct {
	repo repo.FinancialSettingsRepo
}

func NewFinancialSettingsSvcs(i *do.Injector) (FinancialSettingsSvcs, error) {
	return &financialsettingssvcs{
		repo: do.MustInvoke[repo.FinancialSettingsRepo](i),
	}, nil
}

func (s *financialsettingssvcs) init(ctx context.Context) (*models.FinancialSettings, error) {
	settings := &models.FinancialSettings{
		Id:   primitive.NewObjectID(),
		Name: "financialsettings",
		FinancialSettingsDto: models.FinancialSettingsDto{
			ProfitRatio: 0,
			BankAccount: models.BankAccount{
				BankName:          "",
				AccountNumber:     "",
				AccountHolderName: "",
				IBAN:              "",
				SwiftCode:         "",
				Currencies:        []string{},
			},
		},
	}

	if err := s.repo.Add(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *financialsettingssvcs) Get(ctx context.Context) (*models.FinancialSettings, error) {
	settings, err := s.repo.GetByFilter(ctx, bson.M{"name": "financialsettings"})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return s.init(ctx)
		}
		return nil, err
	}

	return settings, nil
}

func (s *financialsettingssvcs) Update(ctx context.Context, data *models.FinancialSettingsDto) (*models.FinancialSettings, error) {
	settings, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}

	settings.FinancialSettingsDto = *data

	filter := bson.M{"_id": settings.Id}
	update := bson.M{"$set": settings}

	updatedSettings, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedSettings, nil

}
