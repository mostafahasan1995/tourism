package messaging

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/nats"
	"larsa-tourism-microservices/pkg/services/messaging/enums"
	"larsa-tourism-microservices/pkg/services/messaging/models"
	"larsa-tourism-microservices/pkg/services/messaging/repo"
	"larsa-tourism-microservices/pkg/services/messaging/template"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageSvcs interface {
	GetAll(ctx context.Context) ([]models.Message, error)
	SendEmail(ctx context.Context, msg *models.Message) error
	GetTemplateMessage(ctx context.Context, msgType enums.MsgTyps, data any) (message string, subject string, err error)
}

type messagesvcs struct {
	repo repo.MessageRepo
}

func NewMessageSvcs(i *do.Injector) (MessageSvcs, error) {
	return &messagesvcs{
		repo: do.MustInvoke[repo.MessageRepo](i),
	}, nil
}

func (m *messagesvcs) GetAll(ctx context.Context) ([]models.Message, error) {
	pipeline := []bson.M{
		{"$match": bson.M{}},
	}

	var result []models.Message
	err := m.repo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *messagesvcs) firstOrUpdate(ctx context.Context, msg *models.Message) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	if msg.Id == primitive.NilObjectID {
		msg.Id = primitive.NewObjectID()
		msg.CreatedAt = time.Now()
		msg.CreatedBy = cfg.User.Id

		if err := m.repo.Add(ctx, msg); err != nil {
			return err
		}

	} else {
		filter := bson.M{"_id": msg.Id}
		update := bson.M{"$set": msg}

		if _, err := m.repo.Patch(ctx, filter, update); err != nil {
			return err
		}
	}

	return nil
}

func (m *messagesvcs) SendEmail(ctx context.Context, msg *models.Message) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	if err := m.firstOrUpdate(ctx, msg); err != nil {
		return err
	}

	email := &models.Email{
		To:      msg.Email,
		Subject: msg.Subject,
		Message: msg.MessageHtml,
		Headers: cfg.Hp,
		Others:  msg.Others,
	}

	nats.Publish("email:send", email)
	return nil
}

func (m *messagesvcs) GetTemplateMessage(ctx context.Context, msgType enums.MsgTyps, data any) (message string, subject string, err error) {
	tpl, ok := template.Templates[msgType]
	if !ok {
		return "", "", errors.New("template not found")
	}

	return tpl.Construct(ctx, data)
}
