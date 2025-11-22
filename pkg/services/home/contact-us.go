package home

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"larsa-tourism-microservices/pkg/services/messaging"
	messagingenums "larsa-tourism-microservices/pkg/services/messaging/enums"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"
	"larsa-tourism-microservices/pkg/services/messaging/template"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactUsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ContactUs, error)
	GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error)
	Add(ctx context.Context, data *models.ContactUsDto) error
	AddMany(ctx context.Context, data []models.ContactUsDto) error
	Update(ctx context.Context, id string, data *models.ContactUsDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetSettings(ctx context.Context) (*models.ContactUsSettings, error)
	AddOrUpdateSettings(ctx context.Context, settings *models.ContactUsSettingsDto) error
	// V2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ContactUsPagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.ContactUs, error)
}

type contactUssvcs struct {
	repo         repo.ContactUsRepo
	settingsRepo repo.ContactUsSettingsRepo
	messagesvcs  messaging.MessageSvcs
}

func NewContactUsSvcs(i *do.Injector) (ContactUsSvcs, error) {
	return &contactUssvcs{
		repo:         do.MustInvoke[repo.ContactUsRepo](i),
		settingsRepo: do.MustInvoke[repo.ContactUsSettingsRepo](i),
		messagesvcs:  do.MustInvoke[messaging.MessageSvcs](i),
	}, nil
}

func (l *contactUssvcs) GetOne(ctx context.Context, id string) (*models.ContactUs, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return l.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (l *contactUssvcs) GetAll(ctx context.Context, filter filter.ContactUsFilter) (models.ContactUsPagination, error) {
	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := l.repo.Count(ctx, filterBody)
	if err != nil {
		return models.ContactUsPagination{}, err
	}

	// Pagination defaults and limits
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = 10 // Default page size
	}
	if size > 100 {
		size = 100 // Maximum page size limit
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Create aggregation pipeline for pagination
	pipeline := []bson.M{
		{"$match": filterBody},
		{"$sort": bson.M{"createdAt": -1}}, // Sort by creation date, newest first
		{"$skip": skip},
		{"$limit": limit},
	}

	var programs []models.ContactUs
	err = l.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &programs)
	})
	if err != nil {
		return models.ContactUsPagination{}, err
	}

	// Prepare pagination result
	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	result := models.ContactUsPagination{
		ContactUs: programs,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *contactUssvcs) Add(ctx context.Context, data *models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	contactUs := &models.ContactUs{
		ContactUsDto: *data, // This preserves AdditionalFields
		Id:           primitive.NewObjectID(),
		Trash:        false,
		Status:       "pending",
		CreatedAt:    time.Now(),
		CreatedBy:    userId,
		UpdatedAt:    time.Now(),
		UpdatedBy:    userId,
	}

	if err := l.repo.Add(ctx, contactUs); err != nil {
		return err
	}

	err = l.SendContactUsEmail(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (l *contactUssvcs) AddMany(ctx context.Context, data []models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("empty data array")
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	var contactUsArray []any
	for _, flr := range data {
		contactUs := &models.ContactUs{
			ContactUsDto: flr, // This preserves AdditionalFields
			Id:           primitive.NewObjectID(),
			Trash:        false,
			Status:       "pending",
			CreatedAt:    time.Now(),
			CreatedBy:    userId,
			UpdatedAt:    time.Now(),
			UpdatedBy:    userId,
		}
		err = l.SendContactUsEmail(ctx, &flr)
		if err != nil {
			return err
		}
		contactUsArray = append(contactUsArray, contactUs)
	}

	err = l.repo.AddMany(ctx, contactUsArray)
	if err != nil {
		return err
	}

	return nil
}
func (l *contactUssvcs) SendContactUsEmail(ctx context.Context, data *models.ContactUsDto) error {
	settings, err := l.GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get contact us settings: %w", err)
	}

	// Validate that email recipients are configured
	if len(settings.Emails) == 0 {
		return fmt.Errorf("no email recipients configured in contact us settings")
	}

	tplData := template.ContactUsTplData{
		FullName:        data.FullName,
		EmailAddress:    data.EmailAddress,
		PhoneNumber:     data.PhoneNumber,
		HowDidYouFindUs: data.HowDidYouFindUs,
		Message:         data.Message,
		Additional:      data.AdditionalFields,
	}
	body, subject, err := l.messagesvcs.GetTemplateMessage(ctx, messagingenums.CONTACTUS, &tplData)
	if err != nil {
		return fmt.Errorf("failed to generate email template: %w", err)
	}

	var emailErrors []error
	for _, email := range settings.Emails {
		emailMsg := &messagingmodels.Message{
			Type:        messagingenums.CONTACTUS,
			Email:       email,
			Subject:     subject,
			Message:     body,
			MessageHtml: body,
			Target:      "email",
			Others: map[string]any{
				// Add Reply-To header so recipients can reply directly to the submitter
				"replyTo":      data.EmailAddress,
				"replyToName":  data.FullName,
				"fromName":     "Contact Form", // Friendly sender name
				"contactEmail": data.EmailAddress,
				"contactName":  data.FullName,
			},
		}
		err = l.messagesvcs.SendEmail(ctx, emailMsg)
		if err != nil {
			emailErrors = append(emailErrors, fmt.Errorf("failed to send email to %s: %w", email, err))
		}
	}

	// If all emails failed, return error
	if len(emailErrors) == len(settings.Emails) {
		return fmt.Errorf("failed to send emails to all recipients: %v", emailErrors)
	}

	// If some emails failed, log but don't fail the request
	if len(emailErrors) > 0 {
		// Log partial failures (you might want to use a proper logger here)
		fmt.Printf("Warning: Failed to send emails to some recipients: %v\n", emailErrors)
	}

	return nil
}
func (a *contactUssvcs) Update(ctx context.Context, id string, data *models.ContactUsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing contact to preserve created fields
	existing, err := a.GetOne(ctx, id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	contactUs := &models.ContactUs{
		ContactUsDto: *data, // This preserves AdditionalFields
		Id:           _id,
		Trash:        false,
		Status:       existing.Status, // Preserve existing status
		CreatedAt:    existing.CreatedAt,
		CreatedBy:    existing.CreatedBy,
		UpdatedBy:    userId,
		UpdatedAt:    time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": contactUs}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *contactUssvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	// Build the update document with only the fields that are being updated
	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	// Add the specific fields to update
	for key, value := range updates {
		switch key {
		case "status":
			updateDoc["status"] = value
		case "fullName":
			updateDoc["fullName"] = value
		case "emailAddress":
			updateDoc["emailAddress"] = value
		case "phoneNumber":
			// Handle phoneNumber as PhoneNumber struct
			if phoneMap, ok := value.(map[string]interface{}); ok {
				phoneNumber := types.PhoneNumber{}
				if pre, ok := phoneMap["pre"].(string); ok {
					phoneNumber.Pre = pre
				}
				if content, ok := phoneMap["content"].(string); ok {
					phoneNumber.Content = content
				}
				updateDoc["phoneNumber"] = phoneNumber
			} else {
				updateDoc["phoneNumber"] = value
			}
		case "howDidYouFindUs":
			updateDoc["howDidYouFindUs"] = value
		case "message":
			updateDoc["message"] = value
		default:
			// Handle additional fields
			updateDoc["additionalFields."+key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *contactUssvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated (public endpoint)
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Use nil ObjectID for anonymous users
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *contactUssvcs) GetSettings(ctx context.Context) (*models.ContactUsSettings, error) {
	var settings []models.ContactUsSettings
	err := a.settingsRepo.Aggregate(ctx, []bson.M{}, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &settings)
	})
	if err != nil {
		return nil, err
	}
	if len(settings) == 0 {
		return nil, errors.New("settings not found")
	}
	return &settings[0], nil
}

func (a *contactUssvcs) AddOrUpdateSettings(ctx context.Context, settings *models.ContactUsSettingsDto) error {
	existing, err := a.GetSettings(ctx)
	fmt.Println("Existing: ", settings)
	if existing == nil || err != nil {
		err = a.settingsRepo.Add(ctx, &models.ContactUsSettings{
			Id:     primitive.NewObjectID(),
			Emails: settings.Emails,
		})
		if err != nil {
			return err
		}
	} else {
		existing.Emails = settings.Emails
		_, err = a.settingsRepo.Patch(ctx, bson.M{"_id": existing.Id}, bson.M{"$set": existing})
		if err != nil {
			return err
		}
	}
	return nil
}

// V2
func (a *contactUssvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ContactUsPagination, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)
	count, err := a.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.ContactUs
	err = a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}
	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pg := common.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}
	return &models.ContactUsPagination{
		ContactUs:  result,
		Pagination: pg,
	}, nil
}

func (a *contactUssvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.ContactUs, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}
	pipline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	var result []models.ContactUs
	err = a.repo.Aggregate(ctx, pipline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
