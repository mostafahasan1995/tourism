package home

import (
	"context"
	"errors"
	"fmt"
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

type ContactUsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ContactUs, error)
	GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error)
	Add(ctx context.Context, data *models.ContactUsDto) error
	AddMany(ctx context.Context, data []models.ContactUsDto) error
	Update(ctx context.Context, id string, data *models.ContactUsDto) error
	Delete(ctx context.Context, id string) error
}

type contactUssvcs struct {
	repo repo.ContactUsRepo
}

func NewContactUsSvcs(i *do.Injector) (ContactUsSvcs, error) {
	return &contactUssvcs{
		repo: do.MustInvoke[repo.ContactUsRepo](i),
	}, nil
}

func (l *contactUssvcs) GetOne(ctx context.Context, id string) (*models.ContactUs, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *contactUssvcs) GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.ContactUsPagination{}, err
	}

	return data, nil
}

func (l *contactUssvcs) Add(ctx context.Context, data *models.ContactUsDto) error {
	fmt.Printf("Starting Add method with data: %+v\n", data)
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		return err
	}
	fmt.Printf("Got config with DB: %s\n", cfg.Db)

	contactUs := &models.ContactUs{
		ContactUsDto: models.ContactUsDto{
			FullName:        data.FullName,
			EmailAddress:    data.EmailAddress,
			PhoneNumber:     data.PhoneNumber,
			HowDidYouFindUs: data.HowDidYouFindUs,
			Message:         data.Message,
		},

		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}
	fmt.Printf("Created contact us object: %+v\n", contactUs)

	if err := l.repo.Add(ctx, contactUs); err != nil {
		fmt.Printf("Error adding to repo: %v\n", err)
		return err
	}

	fmt.Println("Successfully added contact form submission")
	return nil
}

func (l *contactUssvcs) AddMany(ctx context.Context, data []models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	var writeOps []mongo.WriteModel

	for _, flr := range data {
		contactUs := &models.ContactUs{
			ContactUsDto: models.ContactUsDto{
				FullName:        flr.FullName,
				EmailAddress:    flr.EmailAddress,
				PhoneNumber:     flr.PhoneNumber,
				HowDidYouFindUs: flr.HowDidYouFindUs,
				Message:         flr.Message,
			},

			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}
		writeOp := mongo.NewInsertOneModel()
		writeOp.SetDocument(contactUs)
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
func (a *contactUssvcs) Update(ctx context.Context, id string, data *models.ContactUsDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *contactUssvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
