package home

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrustedPartnersSvcs interface {
	GetOne(ctx context.Context, id string) (*models.TrustedPartner, error)
	GetAll(ctx context.Context) ([]models.TrustedPartner, error)
	Add(ctx context.Context, data *models.TrustedPartnerDto) (*models.TrustedPartner, error)
	AddMany(ctx context.Context, data []models.TrustedPartnerDto) error
	Update(ctx context.Context, id string, data *models.TrustedPartnerDto) (*models.TrustedPartner, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.TrustedPartner, error)
	Delete(ctx context.Context, id string) error
}

type trustedPartnerssvcs struct {
	repo repo.TrustedPartnersRepo
}

func NewTrustedPartnersSvcs(i *do.Injector) (TrustedPartnersSvcs, error) {
	return &trustedPartnerssvcs{
		repo: do.MustInvoke[repo.TrustedPartnersRepo](i),
	}, nil
}

func (s *trustedPartnerssvcs) GetOne(ctx context.Context, id string) (*models.TrustedPartner, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *trustedPartnerssvcs) GetAll(ctx context.Context) ([]models.TrustedPartner, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	var result []models.TrustedPartner
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *trustedPartnerssvcs) Add(ctx context.Context, data *models.TrustedPartnerDto) (*models.TrustedPartner, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	partner := &models.TrustedPartner{
		TrustedPartnerDto: models.TrustedPartnerDto{
			Title:        data.Title,
			ProgramTitle: data.ProgramTitle,
			ProgramId:    data.ProgramId,
			Description:  data.Description,
			Image:        data.Image,
			URL:          data.URL,
			DisplayOrder: data.DisplayOrder,
			IsActive:     data.IsActive,
		},

		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	if err := s.repo.Add(ctx, partner); err != nil {
		return nil, err
	}

	return partner, nil
}

func (s *trustedPartnerssvcs) AddMany(ctx context.Context, data []models.TrustedPartnerDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	var writeOps []mongo.WriteModel

	for _, item := range data {
		partner := &models.TrustedPartner{
			TrustedPartnerDto: models.TrustedPartnerDto{
				Title:        item.Title,
				ProgramTitle: item.ProgramTitle,
				ProgramId:    item.ProgramId,
				Description:  item.Description,
				Image:        item.Image,
				URL:          item.URL,
				DisplayOrder: item.DisplayOrder,
				IsActive:     item.IsActive,
			},

			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}

		writeOp := mongo.NewInsertOneModel()
		writeOp.SetDocument(partner)
		writeOps = append(writeOps, writeOp)
	}

	if len(writeOps) == 0 {
		return errors.New("empty write ops")
	}

	_, errInsrt := s.repo.BulkWrite(ctx, writeOps)
	if errInsrt != nil {
		return errInsrt
	}

	return nil
}

func (s *trustedPartnerssvcs) Update(ctx context.Context, id string, data *models.TrustedPartnerDto) (*models.TrustedPartner, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get the existing partner first to preserve creation info
	existingPartner, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	partner := &models.TrustedPartner{
		TrustedPartnerDto: *data,
		Id:                _id,
		Trash:             false,
		CreatedAt:         existingPartner.CreatedAt,
		CreatedBy:         existingPartner.CreatedBy,
		UpdatedAt:         time.Now(),
		UpdatedBy:         cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": partner}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *trustedPartnerssvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.TrustedPartner, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}

	for key, value := range updates {
		updateDoc[key] = value
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	result, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *trustedPartnerssvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}
