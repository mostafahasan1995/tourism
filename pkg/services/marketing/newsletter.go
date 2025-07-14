package marketing

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/services/marketing/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsletterSvcs interface {
	Subscribe(ctx context.Context, email string) error
	Unsubscribe(ctx context.Context, email string) error
}

type newsletterSvcs struct {
	repo repo.NewsletterRepo
}

func NewNewsletterSvcs(i *do.Injector) (NewsletterSvcs, error) {
	return &newsletterSvcs{
		repo: do.MustInvoke[repo.NewsletterRepo](i),
	}, nil
}

func (s *newsletterSvcs) Subscribe(ctx context.Context, email string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	existing, err := s.repo.GetByFilter(ctx, bson.M{"email": email})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if existing != nil {
		return helpers.BadRequest("Email already exists")
	}

	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}
	newsletter := &models.Newsletter{
		NewsletterDTO: models.NewsletterDTO{
			Email: email,
		},
		Status:    models.NewsletterStatus(models.NewsletterStatusActive),
		UpdatedBy: userId,
		UpdatedAt: time.Now(),
		CreatedBy: userId,
		CreatedAt: time.Now(),
	}
	fmt.Println("newsletter cre", newsletter)
	return s.repo.Add(ctx, newsletter)
}
func (s *newsletterSvcs) Unsubscribe(ctx context.Context, email string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	existing, err := s.repo.GetByFilter(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	if existing == nil {
		return helpers.BadRequest("Email not found")
	}
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}
	existing.Status = models.NewsletterStatusInactive
	existing.UpdatedBy = userId
	existing.UpdatedAt = time.Now()
	filter := bson.M{"_id": existing.Id[0]}
	update := bson.M{"$set": existing}
	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
