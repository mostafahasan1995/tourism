package marketing

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/services/marketing/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"larsa-tourism-microservices/pkg/services/messaging"
	messagingenums "larsa-tourism-microservices/pkg/services/messaging/enums"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"
	messagingtpls "larsa-tourism-microservices/pkg/services/messaging/template"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
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
	repo        repo.NewsletterRepo
	messagesvcs messaging.MessageSvcs
}

func NewNewsletterSvcs(i *do.Injector) (NewsletterSvcs, error) {
	return &newsletterSvcs{
		repo:        do.MustInvoke[repo.NewsletterRepo](i),
		messagesvcs: do.MustInvoke[messaging.MessageSvcs](i),
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

	err = s.repo.Add(ctx, newsletter)
	if err != nil {
		return err
	}

	// Send welcome email
	s.SendWelcomeNewsletterEmail(ctx, newsletter.Email)

	return nil
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

func (s *newsletterSvcs) SendWelcomeNewsletterEmail(ctx context.Context, email string) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return
	}
	companyName := ""
	if cfg.Hp != nil {
		if val, err := common.GetOptionValue("BUSINESS_NAME", cfg.Hp); err == nil {
			if str, ok := val.(string); ok {
				companyName = str
			}
		}
	}
	tplData := &messagingtpls.NewsletterWelcomeTplData{
		Email:       email,
		CompanyName: companyName,
	}
	body, subject, err := s.messagesvcs.GetTemplateMessage(ctx, messagingenums.WELCOME_NEWSLETTER, tplData)
	if err == nil {
		emailMsg := &messagingmodels.Message{
			Type:        messagingenums.WELCOME_NEWSLETTER,
			Email:       email,
			Subject:     subject,
			Message:     body,
			MessageHtml: body,
			Target:      "email",
			Others:      map[string]any{},
		}
		s.messagesvcs.SendEmail(ctx, emailMsg)
	}
}
