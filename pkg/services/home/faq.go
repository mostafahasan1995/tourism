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
	"strings"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// FAQ Page Service Interface
type FaqPageSvcs interface {
	GetOne(ctx context.Context, id string) (*models.FaqPage, error)
	GetAll(ctx context.Context, filter filter.FaqPageFilter) (models.FaqPagePagination, error)
	GetComplete(ctx context.Context, id string) (*models.FaqPageComplete, error)
	Add(ctx context.Context, data *models.FaqPageDto) error
	AddMany(ctx context.Context, data []models.FaqPageDto) error
	Update(ctx context.Context, id string, data *models.FaqPageDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	InitializeStaticPages(ctx context.Context) error
}

// FAQ Group Service Interface
type FaqGroupSvcs interface {
	GetOne(ctx context.Context, id string) (*models.FaqGroup, error)
	GetAll(ctx context.Context, filter filter.FaqGroupFilter) (models.FaqGroupPagination, error)
	GetWithQuestions(ctx context.Context, id string) (*models.FaqGroupWithQuestions, error)
	Add(ctx context.Context, data *models.FaqGroupDto) error
	AddMany(ctx context.Context, data []models.FaqGroupDto) error
	Update(ctx context.Context, id string, data *models.FaqGroupDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

// FAQ Question Service Interface
type FaqQuestionSvcs interface {
	GetOne(ctx context.Context, id string) (*models.FaqQuestion, error)
	GetAll(ctx context.Context, filter filter.FaqQuestionFilter) (models.FaqQuestionPagination, error)
	Add(ctx context.Context, data *models.FaqQuestionDto) error
	AddMany(ctx context.Context, data []models.FaqQuestionDto) error
	Update(ctx context.Context, id string, data *models.FaqQuestionDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

// Service implementations
type faqPageSvcs struct {
	repo         repo.FaqPageRepo
	groupRepo    repo.FaqGroupRepo
	questionRepo repo.FaqQuestionRepo
}

type faqGroupSvcs struct {
	repo         repo.FaqGroupRepo
	questionRepo repo.FaqQuestionRepo
}

type faqQuestionSvcs struct {
	repo repo.FaqQuestionRepo
}

// Constructors
func NewFaqPageSvcs(i *do.Injector) (FaqPageSvcs, error) {
	return &faqPageSvcs{
		repo:         do.MustInvoke[repo.FaqPageRepo](i),
		groupRepo:    do.MustInvoke[repo.FaqGroupRepo](i),
		questionRepo: do.MustInvoke[repo.FaqQuestionRepo](i),
	}, nil
}

func NewFaqGroupSvcs(i *do.Injector) (FaqGroupSvcs, error) {
	return &faqGroupSvcs{
		repo:         do.MustInvoke[repo.FaqGroupRepo](i),
		questionRepo: do.MustInvoke[repo.FaqQuestionRepo](i),
	}, nil
}

func NewFaqQuestionSvcs(i *do.Injector) (FaqQuestionSvcs, error) {
	return &faqQuestionSvcs{
		repo: do.MustInvoke[repo.FaqQuestionRepo](i),
	}, nil
}

// Helper function to generate slug from name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "&", "and")
	return slug
}

// Helper function to get user ID with nil check
func getUserId(ctx context.Context) primitive.ObjectID {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil || cfg.User == nil {
		return primitive.NilObjectID
	}
	return cfg.User.Id
}

// =============================================================================
// FAQ PAGE SERVICE IMPLEMENTATION
// =============================================================================

func (s *faqPageSvcs) GetOne(ctx context.Context, id string) (*models.FaqPage, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *faqPageSvcs) GetAll(ctx context.Context, filter filter.FaqPageFilter) (models.FaqPagePagination, error) {
	filterBody := filter.ToBsonFilter()

	// Count total documents
	totalCount, err := s.repo.Count(ctx, filterBody)
	if err != nil {
		return models.FaqPagePagination{}, err
	}

	// Pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Simple aggregation pipeline
	pipeline := []bson.M{
		{"$match": filterBody},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var pages []models.FaqPage
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &pages)
	})
	if err != nil {
		return models.FaqPagePagination{}, err
	}

	// Convert to FaqPageWithStats and get counts separately if needed
	var pagesWithStats []models.FaqPageWithStats
	for _, page := range pages {
		// For now, set counts to 0 to avoid performance issues
		// In the future, these can be calculated separately if needed
		pagesWithStats = append(pagesWithStats, models.FaqPageWithStats{
			FaqPage:       page,
			GroupCount:    0,
			QuestionCount: 0,
		})
	}

	// Calculate total pages
	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	return models.FaqPagePagination{
		FaqPages: pagesWithStats,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}, nil
}

func (s *faqPageSvcs) GetComplete(ctx context.Context, id string) (*models.FaqPageComplete, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Get the FAQ page
	page, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get general questions (not in any group)
	var generalQuestions []models.FaqQuestion
	generalPipeline := []bson.M{
		{
			"$match": bson.M{
				"faqPageId": _id,
				"trash":     bson.M{"$ne": true},
				"$or": []bson.M{
					{"faqGroupId": bson.M{"$exists": false}},
					{"faqGroupId": nil},
				},
			},
		},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
	}
	err = s.questionRepo.Aggregate(ctx, generalPipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &generalQuestions)
	})
	if err != nil {
		return nil, err
	}

	// Get groups with their questions
	var groups []models.FaqGroup
	groupPipeline := []bson.M{
		{
			"$match": bson.M{
				"faqPageId": _id,
				"trash":     bson.M{"$ne": true},
				"isActive":  true,
			},
		},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
	}
	err = s.groupRepo.Aggregate(ctx, groupPipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &groups)
	})
	if err != nil {
		return nil, err
	}

	var groupsWithQuestions []models.FaqGroupWithQuestions
	for _, group := range groups {
		var questions []models.FaqQuestion
		questionPipeline := []bson.M{
			{
				"$match": bson.M{
					"faqPageId":  _id,
					"faqGroupId": group.Id,
					"trash":      bson.M{"$ne": true},
					"isActive":   true,
				},
			},
			{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		}
		err := s.questionRepo.Aggregate(ctx, questionPipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &questions)
		})
		if err != nil {
			return nil, err
		}

		groupsWithQuestions = append(groupsWithQuestions, models.FaqGroupWithQuestions{
			FaqGroup:  group,
			Questions: questions,
		})
	}

	return &models.FaqPageComplete{
		FaqPage:          *page,
		GeneralQuestions: generalQuestions,
		Groups:           groupsWithQuestions,
	}, nil
}

func (s *faqPageSvcs) Add(ctx context.Context, data *models.FaqPageDto) error {
	// Get user ID, allowing for anonymous users
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID // Allow anonymous creation
	}

	// Generate slug if not provided
	slug := data.Slug
	if slug == "" {
		slug = generateSlug(data.Name)
	}

	// Create the FAQ page with proper ObjectID
	faqPage := &models.FaqPage{
		Id:          primitive.NewObjectID(),
		Name:        data.Name,
		Slug:        slug,
		Description: data.Description,
		IsActive:    data.IsActive,
		SortOrder:   data.SortOrder,
		Trash:       false,
		CreatedAt:   time.Now(),
		CreatedBy:   userId,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userId,
	}

	// Add to repository
	return s.repo.Add(ctx, faqPage)
}

func (s *faqPageSvcs) AddMany(ctx context.Context, data []models.FaqPageDto) error {
	if len(data) == 0 {
		return errors.New("empty data array")
	}

	userId := getUserId(ctx)
	var faqPages []any

	for _, dto := range data {
		slug := dto.Slug
		if slug == "" {
			slug = generateSlug(dto.Name)
		}

		faqPage := &models.FaqPage{
			Id:          primitive.NewObjectID(),
			Name:        dto.Name,
			Slug:        slug,
			Description: dto.Description,
			IsActive:    dto.IsActive,
			SortOrder:   dto.SortOrder,
			Trash:       false,
			CreatedAt:   time.Now(),
			CreatedBy:   userId,
			UpdatedAt:   time.Now(),
			UpdatedBy:   userId,
		}
		faqPages = append(faqPages, faqPage)
	}

	return s.repo.AddMany(ctx, faqPages)
}

func (s *faqPageSvcs) Update(ctx context.Context, id string, data *models.FaqPageDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing page
	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	// Generate slug if not provided
	slug := data.Slug
	if slug == "" {
		slug = generateSlug(data.Name)
	}

	faqPage := &models.FaqPage{
		Id:          _id,
		Name:        data.Name,
		Slug:        slug,
		Description: data.Description,
		IsActive:    data.IsActive,
		SortOrder:   data.SortOrder,
		Trash:       false,
		CreatedAt:   existing.CreatedAt,
		CreatedBy:   existing.CreatedBy,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userId,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": faqPage}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqPageSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	userId := getUserId(ctx)

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	for key, value := range updates {
		switch key {
		case "name":
			updateDoc["name"] = value
			// Auto-generate slug if name is updated and slug is not provided
			if _, hasSlug := updates["slug"]; !hasSlug {
				if name, ok := value.(string); ok {
					updateDoc["slug"] = generateSlug(name)
				}
			}
		case "slug", "description", "isActive", "sortOrder":
			updateDoc[key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqPageSvcs) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqPageSvcs) InitializeStaticPages(ctx context.Context) error {
	staticPages := []models.FaqPageDto{
		{Name: "About", IsActive: true, SortOrder: 1, Description: "About our company and services"},
		{Name: "Our Agents", IsActive: true, SortOrder: 2, Description: "Information about our travel agents"},
		{Name: "Delegation Travel Package", IsActive: true, SortOrder: 3, Description: "Delegation travel packages and services"},
		{Name: "Delegation Travel Service", IsActive: true, SortOrder: 4, Description: "Delegation travel service information"},
		{Name: "Religious Tourism Package", IsActive: true, SortOrder: 5, Description: "Religious tourism packages and pilgrimages"},
		{Name: "Imkan B2b", IsActive: true, SortOrder: 6, Description: "B2B services and partnerships"},
	}

	successCount := 0
	for _, pageDto := range staticPages {
		// Check if page already exists by name
		existing, _ := s.repo.GetByFilter(ctx, bson.M{
			"name":  pageDto.Name,
			"trash": bson.M{"$ne": true},
		})

		if existing == nil {
			err := s.Add(ctx, &pageDto)
			if err != nil {
				return fmt.Errorf("failed to create static page '%s': %v", pageDto.Name, err)
			}
			successCount++
		}
	}

	if successCount > 0 {
		fmt.Printf("Successfully created %d static FAQ pages\n", successCount)
	}

	return nil
}

// =============================================================================
// FAQ GROUP SERVICE IMPLEMENTATION
// =============================================================================

func (s *faqGroupSvcs) GetOne(ctx context.Context, id string) (*models.FaqGroup, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *faqGroupSvcs) GetAll(ctx context.Context, filter filter.FaqGroupFilter) (models.FaqGroupPagination, error) {
	filterBody := filter.ToBsonFilter()

	totalCount, err := s.repo.Count(ctx, filterBody)
	if err != nil {
		return models.FaqGroupPagination{}, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	pipeline := []bson.M{
		{"$match": filterBody},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var groups []models.FaqGroup
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &groups)
	})
	if err != nil {
		return models.FaqGroupPagination{}, err
	}

	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	return models.FaqGroupPagination{
		FaqGroups: groups,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}, nil
}

func (s *faqGroupSvcs) GetWithQuestions(ctx context.Context, id string) (*models.FaqGroupWithQuestions, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	group, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	var questions []models.FaqQuestion

	// First, try to get active questions
	questionPipeline := []bson.M{
		{
			"$match": bson.M{
				"faqGroupId": _id,
				"trash":      bson.M{"$ne": true},
				"isActive":   true,
			},
		},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
	}

	// Debug: Log the pipeline and group ID
	fmt.Printf("DEBUG: GetWithQuestions - Group ID: %s\n", _id.Hex())
	fmt.Printf("DEBUG: GetWithQuestions - Pipeline: %+v\n", questionPipeline)

	err = s.questionRepo.Aggregate(ctx, questionPipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &questions)
	})
	if err != nil {
		fmt.Printf("DEBUG: GetWithQuestions - Error: %v\n", err)
		return nil, err
	}

	// Debug: Log the number of questions found
	fmt.Printf("DEBUG: GetWithQuestions - Found %d active questions\n", len(questions))

	// If no active questions found, try to get all questions (including inactive) for debugging
	if len(questions) == 0 {
		fmt.Printf("DEBUG: GetWithQuestions - No active questions found, checking all questions...\n")
		allQuestionsPipeline := []bson.M{
			{
				"$match": bson.M{
					"faqGroupId": _id,
					"trash":      bson.M{"$ne": true},
				},
			},
			{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		}

		var allQuestions []models.FaqQuestion
		err = s.questionRepo.Aggregate(ctx, allQuestionsPipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &allQuestions)
		})
		if err == nil {
			fmt.Printf("DEBUG: GetWithQuestions - Found %d total questions (including inactive)\n", len(allQuestions))
			if len(allQuestions) > 0 {
				fmt.Printf("DEBUG: GetWithQuestions - First question: %+v\n", allQuestions[0])
				// For debugging, return all questions regardless of isActive status
				questions = allQuestions
			}
		}
	} else if len(questions) > 0 {
		fmt.Printf("DEBUG: GetWithQuestions - First active question: %+v\n", questions[0])
	}

	// Initialize questions slice if nil to avoid null in JSON response
	if questions == nil {
		questions = []models.FaqQuestion{}
	}

	return &models.FaqGroupWithQuestions{
		FaqGroup:  *group,
		Questions: questions,
	}, nil
}

func (s *faqGroupSvcs) Add(ctx context.Context, data *models.FaqGroupDto) error {
	userId := getUserId(ctx)

	faqGroup := &models.FaqGroup{
		Id:          primitive.NewObjectID(),
		FaqPageId:   data.FaqPageId,
		Name:        data.Name,
		Description: data.Description,
		IsActive:    data.IsActive,
		SortOrder:   data.SortOrder,
		Trash:       false,
		CreatedAt:   time.Now(),
		CreatedBy:   userId,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userId,
	}

	return s.repo.Add(ctx, faqGroup)
}

func (s *faqGroupSvcs) AddMany(ctx context.Context, data []models.FaqGroupDto) error {
	if len(data) == 0 {
		return errors.New("empty data array")
	}

	userId := getUserId(ctx)
	var faqGroups []any

	for _, dto := range data {
		faqGroup := &models.FaqGroup{
			Id:          primitive.NewObjectID(),
			FaqPageId:   dto.FaqPageId,
			Name:        dto.Name,
			Description: dto.Description,
			IsActive:    dto.IsActive,
			SortOrder:   dto.SortOrder,
			Trash:       false,
			CreatedAt:   time.Now(),
			CreatedBy:   userId,
			UpdatedAt:   time.Now(),
			UpdatedBy:   userId,
		}
		faqGroups = append(faqGroups, faqGroup)
	}

	return s.repo.AddMany(ctx, faqGroups)
}

func (s *faqGroupSvcs) Update(ctx context.Context, id string, data *models.FaqGroupDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	faqGroup := &models.FaqGroup{
		Id:          _id,
		FaqPageId:   data.FaqPageId,
		Name:        data.Name,
		Description: data.Description,
		IsActive:    data.IsActive,
		SortOrder:   data.SortOrder,
		Trash:       false,
		CreatedAt:   existing.CreatedAt,
		CreatedBy:   existing.CreatedBy,
		UpdatedAt:   time.Now(),
		UpdatedBy:   userId,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": faqGroup}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqGroupSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	userId := getUserId(ctx)

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	for key, value := range updates {
		switch key {
		case "faqPageId", "name", "description", "isActive", "sortOrder":
			updateDoc[key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqGroupSvcs) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

// =============================================================================
// FAQ QUESTION SERVICE IMPLEMENTATION
// =============================================================================

func (s *faqQuestionSvcs) GetOne(ctx context.Context, id string) (*models.FaqQuestion, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *faqQuestionSvcs) GetAll(ctx context.Context, filter filter.FaqQuestionFilter) (models.FaqQuestionPagination, error) {
	filterBody := filter.ToBsonFilter()

	totalCount, err := s.repo.Count(ctx, filterBody)
	if err != nil {
		return models.FaqQuestionPagination{}, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	pipeline := []bson.M{
		{"$match": filterBody},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var questions []models.FaqQuestion
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &questions)
	})
	if err != nil {
		return models.FaqQuestionPagination{}, err
	}

	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	return models.FaqQuestionPagination{
		FaqQuestions: questions,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}, nil
}

func (s *faqQuestionSvcs) Add(ctx context.Context, data *models.FaqQuestionDto) error {
	userId := getUserId(ctx)

	faqQuestion := &models.FaqQuestion{
		Id:         primitive.NewObjectID(),
		FaqPageId:  data.FaqPageId,
		FaqGroupId: data.FaqGroupId,
		Question:   data.Question,
		Answer:     data.Answer,
		IsActive:   data.IsActive,
		SortOrder:  data.SortOrder,
		Trash:      false,
		CreatedAt:  time.Now(),
		CreatedBy:  userId,
		UpdatedAt:  time.Now(),
		UpdatedBy:  userId,
	}

	return s.repo.Add(ctx, faqQuestion)
}

func (s *faqQuestionSvcs) AddMany(ctx context.Context, data []models.FaqQuestionDto) error {
	if len(data) == 0 {
		return errors.New("empty data array")
	}

	userId := getUserId(ctx)
	var faqQuestions []any

	for _, dto := range data {
		faqQuestion := &models.FaqQuestion{
			Id:         primitive.NewObjectID(),
			FaqPageId:  dto.FaqPageId,
			FaqGroupId: dto.FaqGroupId,
			Question:   dto.Question,
			Answer:     dto.Answer,
			IsActive:   dto.IsActive,
			SortOrder:  dto.SortOrder,
			Trash:      false,
			CreatedAt:  time.Now(),
			CreatedBy:  userId,
			UpdatedAt:  time.Now(),
			UpdatedBy:  userId,
		}
		faqQuestions = append(faqQuestions, faqQuestion)
	}

	return s.repo.AddMany(ctx, faqQuestions)
}

func (s *faqQuestionSvcs) Update(ctx context.Context, id string, data *models.FaqQuestionDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	faqQuestion := &models.FaqQuestion{
		Id:         _id,
		FaqPageId:  data.FaqPageId,
		FaqGroupId: data.FaqGroupId,
		Question:   data.Question,
		Answer:     data.Answer,
		IsActive:   data.IsActive,
		SortOrder:  data.SortOrder,
		Trash:      false,
		CreatedAt:  existing.CreatedAt,
		CreatedBy:  existing.CreatedBy,
		UpdatedAt:  time.Now(),
		UpdatedBy:  userId,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": faqQuestion}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqQuestionSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	userId := getUserId(ctx)

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	for key, value := range updates {
		switch key {
		case "faqPageId", "faqGroupId", "question", "answer", "isActive", "sortOrder":
			updateDoc[key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *faqQuestionSvcs) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}
