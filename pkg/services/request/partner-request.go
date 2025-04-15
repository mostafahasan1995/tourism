package request

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/request/filter"
	"larsa-tourism-microservices/pkg/services/request/models"
	"larsa-tourism-microservices/pkg/services/request/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartnerRequestSvcs interface {
	GetOne(ctx context.Context, id string) (*models.PartnerRequest, error)
	GetAll(ctx context.Context, filter filter.PartnerRequestFilter) (models.PartnerRequestPagination, error)
	Add(ctx context.Context, data *models.PartnerRequestDto) error
	AddMany(ctx context.Context, data []models.PartnerRequestDto) error
	Update(ctx context.Context, id string, data *models.PartnerRequestDto) error
	Delete(ctx context.Context, id string) error
}

type partnerRequestsvcs struct {
	repo repo.PartnerRequestRepo
}

func NewPartnerRequestSvcs(i *do.Injector) (PartnerRequestSvcs, error) {
	return &partnerRequestsvcs{
		repo: do.MustInvoke[repo.PartnerRequestRepo](i),
	}, nil
}

func (l *partnerRequestsvcs) GetOne(ctx context.Context, id string) (*models.PartnerRequest, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *partnerRequestsvcs) GetAll(ctx context.Context, filter filter.PartnerRequestFilter) (models.PartnerRequestPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.PartnerRequestPagination{}, err
	}

	return data, nil
}

func (l *partnerRequestsvcs) Add(ctx context.Context, data *models.PartnerRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	partnerRequest := &models.PartnerRequest{
		PartnerRequestDto: models.PartnerRequestDto{
			CompanyName:     data.CompanyName,
			BusinessType:    data.BusinessType,
			Website:         data.Website,
			CompanyLocation: data.CompanyLocation,
			ContactDetail: models.ContactDetail{
				FullName:    data.ContactDetail.FullName,
				Position:    data.ContactDetail.Position,
				PhoneNumber: data.ContactDetail.PhoneNumber,
				Email:       data.ContactDetail.Email,
			},
		},
		
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	if err := l.repo.Add(ctx, partnerRequest); err != nil {
		return err
	}

	return nil

}
func (l *partnerRequestsvcs) AddMany(ctx context.Context, data []models.PartnerRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	var writeOps []mongo.WriteModel

	for _, flr := range data {
		partnerRequest := &models.PartnerRequest{
			PartnerRequestDto: models.PartnerRequestDto{
				CompanyName:     flr.CompanyName,
				BusinessType:    flr.BusinessType,
				Website:         flr.Website,
				CompanyLocation: flr.CompanyLocation,
				ContactDetail: models.ContactDetail{
					FullName:    flr.ContactDetail.FullName,
					Position:    flr.ContactDetail.Position,
					PhoneNumber: flr.ContactDetail.PhoneNumber,
					Email:       flr.ContactDetail.Email,
				},
			},			
			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}
		writeOp := mongo.NewInsertOneModel()
		writeOp.SetDocument(partnerRequest)
		writeOps = append(writeOps, writeOp)
		if len(writeOps) == 0 {
			return errors.New("empty write ops")
		}
	}
	_, errInsrt := l.repo.BulkWrite(ctx, writeOps)
	if errInsrt != nil {
		return err
	}

	return nil
}
func (a *partnerRequestsvcs) Update(ctx context.Context, id string, data *models.PartnerRequestDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *partnerRequestsvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
