package member

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/gateway"
	"larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/messaging"
	messagingenums "larsa-tourism-microservices/pkg/services/messaging/enums"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"
	messagingtpls "larsa-tourism-microservices/pkg/services/messaging/template"
	"larsa-tourism-microservices/pkg/util"

	"github.com/goccy/go-json"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberAuthSvcs interface {
	AddCredentials(ctx context.Context, data any) (userId primitive.ObjectID, pass string, err error)
	UpdateCredentials(ctx context.Context, password string, data any) error
	GetInvitationEmail(ctx context.Context, password string, data any) (*messagingmodels.Message, error)
	GetAccountUpdatedEmail(ctx context.Context, password string, data any) (*messagingmodels.Message, error)
	DeleteCredentials(ctx context.Context, userId string) (err error)
}

type memberAuthSvcs struct {
	gateway     gateway.Gateway
	messagesvcs messaging.MessageSvcs
}

func NewMemberAuthSvcs(i *do.Injector) (MemberAuthSvcs, error) {
	return &memberAuthSvcs{
		gateway:     do.MustInvoke[gateway.Gateway](i),
		messagesvcs: do.MustInvoke[messaging.MessageSvcs](i),
	}, nil
}

func (m *memberAuthSvcs) AddCredentials(ctx context.Context, data any) (userId primitive.ObjectID, pass string, err error) {
	var user map[string]any
	var password string

	switch member := data.(type) {
	case *models.Agent:
		user = map[string]any{
			"firstName": member.Name,
			"lastName":  "-",
			"email":     member.Security.Email,
		}
		password = member.Security.NewPassword
	case *models.Customer:
		user = map[string]any{
			"firstName": member.Name,
			"lastName":  "-",
			"email":     member.Security.Email,
		}
		password = member.Security.NewPassword
	}

	if password == "" {
		password = util.GeneratePassword(8, 2, 2, 2)
	}

	user["password"] = password

	resp, err := m.gateway.Request(ctx, "users", "users/", "POST", "", user)

	zeroId := primitive.NilObjectID

	if err != nil {
		return zeroId, "", errors.New("error adding user")
	} else if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 409:
			return zeroId, "", ErrDupliateEmail
		case 401:
			return zeroId, "", ErrUnauthorized
		case 403:
			return zeroId, "", ErrForbidden
		default:
			return zeroId, "", errors.New("error adding user")
		}

	}

	type aux struct {
		User struct {
			Id primitive.ObjectID `json:"_id"`
		} `json:"user"`
	}

	var _data aux
	if errDec := json.NewDecoder(resp.Body).Decode(&_data); errDec != nil {
		return zeroId, "", errDec
	}

	return _data.User.Id, password, nil

}

func (m *memberAuthSvcs) UpdateCredentials(ctx context.Context, password string, data any) error {
	var user map[string]any
	var id string

	switch member := data.(type) {
	case *models.Agent:
		user = map[string]any{
			"firstName": member.Name,
			"lastName":  "-",
			"email":     member.Security.Email,
		}
		id = member.Id.Hex()
	case *models.Customer:
		user = map[string]any{
			"firstName": member.Name,
			"lastName":  "-",
			"email":     member.Security.Email,
		}
		id = member.Id.Hex()
	}

	if password != "" {
		user["password"] = password
	}

	resp, err := m.gateway.Request(ctx, "users", "users/"+id, "PATCH", "", user)

	if err != nil {
		return errors.New("error update user")
	} else if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 409:
			return ErrDupliateEmail
		case 401:
			return ErrUnauthorized
		case 403:
			return ErrForbidden
		default:
			return errors.New("error adding user")
		}

	}

	return nil
}

func (m *memberAuthSvcs) DeleteCredentials(ctx context.Context, userId string) (err error) {
	resp, err := m.gateway.Request(ctx, "users", "users/"+userId, "DELETE", "", map[string]any{})
	if err != nil || resp.StatusCode != 200 {
		return errors.New("error deleting user")
	}
	return nil
}

func (m *memberAuthSvcs) GetInvitationEmail(ctx context.Context, password string, data any) (*messagingmodels.Message, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	val, err := common.GetOptionValue("BUSINESS_NAME", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting business name")
	}
	businessName, ok := val.(string)
	if !ok {
		businessName = "[Business Name]"
	}

	plat, err := common.GetOptionValue("PLATFORM_NAME", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting platform name")
	}
	platformName, ok := plat.(string)
	if !ok {
		platformName = "[Platform Name]"
	}

	fUrl, err := common.GetOptionValue("FRONTEND_URL", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting frontend url")
	}

	frontEndUrl, ok := fUrl.(string)
	if !ok {
		return nil, errors.New("error getting frontend url")
	}

	var invitationTplData *messagingtpls.InvetationTplData
	var email string

	switch member := data.(type) {
	case *models.Agent:
		invitationTplData = &messagingtpls.InvetationTplData{
			MemberName:   member.Name,
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}

		email = member.Security.Email
	case *models.Customer:
		invitationTplData = &messagingtpls.InvetationTplData{
			MemberName:   member.Name,
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}
		email = member.Security.Email
	}

	message, subject, err := m.messagesvcs.GetTemplateMessage(ctx, messagingenums.INVITATION, invitationTplData)
	if err != nil {
		return nil, err
	}

	msg := &messagingmodels.Message{
		Type:        messagingenums.INVITATION,
		Email:       email,
		Subject:     subject,
		Message:     message,
		MessageHtml: message,
		Target:      "email", //todo: must set in message service
		Others:      map[string]any{},
	}

	return msg, nil
}

func (m *memberAuthSvcs) GetAccountUpdatedEmail(ctx context.Context, password string, data any) (*messagingmodels.Message, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	val, err := common.GetOptionValue("BUSINESS_NAME", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting business name")
	}
	businessName, ok := val.(string)
	if !ok {
		businessName = "[Business Name]"
	}

	plat, err := common.GetOptionValue("PLATFORM_NAME", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting platform name")
	}
	platformName, ok := plat.(string)
	if !ok {
		platformName = "[Platform Name]"
	}

	fUrl, err := common.GetOptionValue("FRONTEND_URL", cfg.Hp)
	if err != nil {
		return nil, errors.New("error getting frontend url")
	}

	frontEndUrl, ok := fUrl.(string)
	if !ok {
		return nil, errors.New("error getting frontend url")
	}

	var accountUpdatedTplData *messagingtpls.AccountUpdatedTplData
	var email string

	switch member := data.(type) {
	case *models.Agent:
		accountUpdatedTplData = &messagingtpls.AccountUpdatedTplData{
			MemberName:   member.Name,
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}

		email = member.Security.Email
	case *models.Customer:
		accountUpdatedTplData = &messagingtpls.AccountUpdatedTplData{
			MemberName:   member.Name,
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}
		email = member.Security.Email
	}

	message, subject, err := m.messagesvcs.GetTemplateMessage(ctx, messagingenums.ACCOUNTUPDATED, accountUpdatedTplData)
	if err != nil {
		return nil, err
	}

	msg := &messagingmodels.Message{
		Type:        messagingenums.INVITATION,
		Email:       email,
		Subject:     subject,
		Message:     message,
		MessageHtml: message,
		Target:      "email", //todo: must set in message service
		Others:      map[string]any{},
	}

	return msg, nil
}
