package home

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VipCarRequestSvcs interface {
	GetOne(ctx context.Context, id string) (*models.VipCarRequest, error)
	GetAll(ctx context.Context, filter filter.VipCarRequestFilter) (models.VipCarRequestPagination, error)
	Add(ctx context.Context, data *models.VipCarRequestDto) error
	AddMany(ctx context.Context, data []models.VipCarRequestDto) error
	Update(ctx context.Context, id string, data *models.VipCarRequestDto) error
	Delete(ctx context.Context, id string) error
}

type vipCarRequestsvcs struct {
	repo repo.VipCarRequestRepo
}

func NewVipCarRequestSvcs(i *do.Injector) (VipCarRequestSvcs, error) {
	return &vipCarRequestsvcs{
		repo: do.MustInvoke[repo.VipCarRequestRepo](i),
	}, nil
}

func (l *vipCarRequestsvcs) GetOne(ctx context.Context, id string) (*models.VipCarRequest, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *vipCarRequestsvcs) GetAll(ctx context.Context, filter filter.VipCarRequestFilter) (models.VipCarRequestPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.VipCarRequestPagination{}, err
	}

	return data, nil
}

func (l *vipCarRequestsvcs) Add(ctx context.Context, data *models.VipCarRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	vipCarRequest := &models.VipCarRequest{
		VipCarRequestDto: models.VipCarRequestDto{
			Location:              data.Location,
			Capacity:              data.Capacity,
			DriverLanguagesSpoken: data.DriverLanguagesSpoken,
			LuxuryFeatures:        data.LuxuryFeatures,
			StartDate:             data.StartDate,
			EndDate:               data.EndDate,
			StartTime:             data.StartTime,
			EndTime:               data.EndTime,
			CarTypeId:             data.CarTypeId,
		},
		
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	if err := l.repo.Add(ctx, vipCarRequest); err != nil {
		return err
	}

	return nil

}
func (l *vipCarRequestsvcs) AddMany(ctx context.Context, data []models.VipCarRequestDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	var writeOps []mongo.WriteModel

	for _, flr := range data {
		vipCarRequest := &models.VipCarRequest{
			VipCarRequestDto: models.VipCarRequestDto{
				Location:              flr.Location,
				Capacity:              flr.Capacity,
				DriverLanguagesSpoken: flr.DriverLanguagesSpoken,
				LuxuryFeatures:        flr.LuxuryFeatures,
				StartDate:             flr.StartDate,
				EndDate:               flr.EndDate,
				StartTime:             flr.StartTime,
				EndTime:               flr.EndTime,
				CarTypeId:             flr.CarTypeId,
			},			
			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}
		writeOp := mongo.NewInsertOneModel()
		writeOp.SetDocument(vipCarRequest)
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
func (a *vipCarRequestsvcs) Update(ctx context.Context, id string, data *models.VipCarRequestDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *vipCarRequestsvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
