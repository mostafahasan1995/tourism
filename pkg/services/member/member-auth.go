package member

import (
	"context"
	"errors"
	"fmt"
	"io"
	"larsa-tourism-microservices/pkg/gateway"
	"larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/messaging"
	messagingenums "larsa-tourism-microservices/pkg/services/messaging/enums"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"
	messagingtpls "larsa-tourism-microservices/pkg/services/messaging/template"
	"larsa-tourism-microservices/pkg/util"
	"strings"

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

	// Get all roles from auth service
	resp2, err2 := m.gateway.Request(ctx, "users", "roles/all", "GET", "", map[string]any{})
	if err2 != nil {
		return primitive.NilObjectID, "", fmt.Errorf("error requesting roles: %w", err2)
	}
	defer resp2.Body.Close()

	// Read response body into bytes (only once - body can only be read once)
	bodyBytes, errRead := io.ReadAll(resp2.Body)
	if errRead != nil {
		return primitive.NilObjectID, "", fmt.Errorf("error reading roles response: %w", errRead)
	}

	// Check status code after reading body
	if resp2.StatusCode != 200 {
		bodyStr := string(bodyBytes)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return primitive.NilObjectID, "", fmt.Errorf("error getting roles from auth service (status %d): %s", resp2.StatusCode, bodyStr)
	}

	// Check if body is empty
	if len(bodyBytes) == 0 {
		return primitive.NilObjectID, "", errors.New("empty response body from roles service")
	}

	// Trim whitespace and check if it's just whitespace
	bodyStr := string(bodyBytes)
	bodyStr = strings.TrimSpace(bodyStr)
	if bodyStr == "" || bodyStr == "null" || bodyStr == "[]" {
		return primitive.NilObjectID, "", fmt.Errorf("empty or null response body from roles service (length: %d)", len(bodyBytes))
	}

	// Parse roles response - it's an array of role objects
	type Capability struct {
		Id       primitive.ObjectID `json:"_id"`
		Name     string             `json:"name"`
		Label    string             `json:"label"`
		Reserved bool               `json:"reserved"`
	}

	type Role struct {
		Id           primitive.ObjectID `json:"_id"`
		Name         string             `json:"name"`
		Label        string             `json:"label"`
		Reserved     bool               `json:"reserved"`
		Capabilities []Capability       `json:"capabilities"`
		CreatedBy    primitive.ObjectID `json:"createdBy,omitempty"`
	}

	// Try to parse as direct array first
	var roles []Role
	// Use the trimmed string for parsing
	trimmedBytes := []byte(bodyStr)
	if errDec := json.Unmarshal(trimmedBytes, &roles); errDec != nil {
		// If that fails, try parsing as wrapped object
		var wrappedResponse struct {
			Roles []Role `json:"roles"`
			Data  []Role `json:"data"`
		}
		if errDec2 := json.Unmarshal(trimmedBytes, &wrappedResponse); errDec2 != nil {
			// Show first 500 chars of body for debugging
			bodyPreview := bodyStr
			if len(bodyPreview) > 500 {
				bodyPreview = bodyPreview[:500] + "..."
			}
			return primitive.NilObjectID, "", fmt.Errorf("error parsing roles response (body length: %d, status: %d, preview: %s, original error: %v, wrapped error: %v)", len(bodyBytes), resp2.StatusCode, bodyPreview, errDec, errDec2)
		}
		if len(wrappedResponse.Roles) > 0 {
			roles = wrappedResponse.Roles
		} else if len(wrappedResponse.Data) > 0 {
			roles = wrappedResponse.Data
		} else {
			return primitive.NilObjectID, "", errors.New("no roles found in response")
		}
	}

	// Find agent role IDs by searching for "agent" in role name (case-insensitive)
	var agentRoleIDs []string
	for _, role := range roles {
		roleNameLower := strings.ToLower(role.Name)
		if strings.Contains(roleNameLower, "agent") {
			agentRoleIDs = append(agentRoleIDs, role.Id.Hex())
		}
	}

	// If no agent roles found, return error
	if len(agentRoleIDs) == 0 {
		return primitive.NilObjectID, "", errors.New("no agent roles found in auth service")
	}

	switch member := data.(type) {
	case *models.Agent:
		user = map[string]any{
			"firstName": member.Name.GetContentByLang("en"),
			"lastName":  "-",
			"email":     member.Security.Email,
			"roles":     agentRoleIDs, // Use dynamically found agent role IDs
		}
		password = member.Security.NewPassword
	case *models.Customer:
		user = map[string]any{
			"firstName": member.Name.GetContentByLang("en"),
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
			"firstName": member.Name.GetContentByLang("en"),
			"lastName":  "-",
			"email":     member.Security.Email,
		}
		id = member.Id.Hex()
	case *models.Customer:
		user = map[string]any{
			"firstName": member.Name.GetContentByLang("en"),
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
			MemberName:   member.Name.GetContentByLang("en"),
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}

		email = member.Security.Email
	case *models.Customer:
		invitationTplData = &messagingtpls.InvetationTplData{
			MemberName:   member.Name.GetContentByLang("en"),
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
			MemberName:   member.Name.GetContentByLang("en"),
			CompanyName:  businessName,
			PlatformName: platformName,
			Email:        member.Security.Email,
			Password:     password,
			Link:         frontEndUrl + "/auth/login",
		}

		email = member.Security.Email
	case *models.Customer:
		accountUpdatedTplData = &messagingtpls.AccountUpdatedTplData{
			MemberName:   member.Name.GetContentByLang("en"),
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
