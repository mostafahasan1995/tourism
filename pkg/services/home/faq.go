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
	"larsa-tourism-microservices/pkg/util"
	"math"
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

	// V2
	GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqPagePagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqPageWithStats, error)
}

// FAQ Group Service Interface
type FaqGroupSvcs interface {
	GetOne(ctx context.Context, id string) (*models.FaqGroup, error)
	GetAll(ctx context.Context, filter filter.FaqGroupFilter) (models.FaqGroupPagination, error)
	GetWithQuestions(ctx context.Context, id string) (*models.FaqGroupWithQuestions, error)
	SearchQuestions(ctx context.Context, pageId string, searchTerm string, page, size int) (models.FaqSearchResultPagination, error)
	Add(ctx context.Context, data *models.FaqGroupDto) error
	AddMany(ctx context.Context, data []models.FaqGroupDto) error
	Update(ctx context.Context, id string, data *models.FaqGroupDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// V2
	GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqGroupPagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqGroup, error)
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

	// V2
	GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqQuestionPagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqQuestion, error)
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
	pageRepo     repo.FaqPageRepo
}

type faqQuestionSvcs struct {
	repo      repo.FaqQuestionRepo
	groupRepo repo.FaqGroupRepo
	pageRepo  repo.FaqPageRepo
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
		pageRepo:     do.MustInvoke[repo.FaqPageRepo](i),
	}, nil
}

func NewFaqQuestionSvcs(i *do.Injector) (FaqQuestionSvcs, error) {
	return &faqQuestionSvcs{
		repo:      do.MustInvoke[repo.FaqQuestionRepo](i),
		groupRepo: do.MustInvoke[repo.FaqGroupRepo](i),
		pageRepo:  do.MustInvoke[repo.FaqPageRepo](i),
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

// Helper function to update parent FAQ page's updatedAt when child entities are modified
func (s *faqPageSvcs) updateParentPageTimestamp(ctx context.Context, pageId primitive.ObjectID) error {
	userId := getUserId(ctx)
	filter := bson.M{"_id": pageId}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": userId,
		},
	}
	_, err := s.repo.Patch(ctx, filter, update)
	return err
}

// Helper function for group service to update parent page
func (s *faqGroupSvcs) updateParentPageTimestamp(ctx context.Context, pageId primitive.ObjectID) error {
	userId := getUserId(ctx)
	filter := bson.M{"_id": pageId}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": userId,
		},
	}
	_, err := s.pageRepo.Patch(ctx, filter, update)
	return err
}

// Helper function for question service to update parent page
func (s *faqQuestionSvcs) updateParentPageTimestamp(ctx context.Context, pageId primitive.ObjectID) error {
	userId := getUserId(ctx)

	// Update the page timestamp
	filter := bson.M{"_id": pageId}
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": userId,
		},
	}
	_, err := s.pageRepo.Patch(ctx, filter, update)
	return err
}

// updateParentGroupAndPageTimestamp updates both the parent group and grandparent page timestamps
func (s *faqQuestionSvcs) updateParentGroupAndPageTimestamp(ctx context.Context, groupId *primitive.ObjectID, pageId primitive.ObjectID) error {
	userId := getUserId(ctx)
	now := time.Now()
	var actualPageId primitive.ObjectID

	// Update the parent group timestamp (only if groupId is not nil)
	if groupId != nil {
		groupFilter := bson.M{"_id": *groupId}
		groupUpdate := bson.M{
			"$set": bson.M{
				"updatedAt": now,
				"updatedBy": userId,
			},
		}
		_, err := s.groupRepo.Patch(ctx, groupFilter, groupUpdate)
		if err != nil {
			fmt.Printf("Warning: Failed to update parent group timestamp: %v\n", err)
		}

		// Get the correct page ID from the group
		groupFilter = bson.M{"_id": *groupId}
		group, err := s.groupRepo.GetByFilter(ctx, groupFilter)
		if err != nil {
			fmt.Printf("Warning: Failed to get group to find page ID: %v\n", err)
			return err
		}
		actualPageId = group.FaqPageId
	} else {
		// If no group, use the provided pageId (though it might be incorrect)
		actualPageId = pageId
	}

	// Update the grandparent page timestamp using the correct page ID
	pageFilter := bson.M{"_id": actualPageId}
	pageUpdate := bson.M{
		"$set": bson.M{
			"updatedAt": now,
			"updatedBy": userId,
		},
	}
	_, err := s.pageRepo.Patch(ctx, pageFilter, pageUpdate)
	if err != nil {
		fmt.Printf("Warning: Failed to update grandparent page timestamp: %v\n", err)
		return err
	}

	return nil
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

	// Convert to FaqPageWithStats and calculate dynamic counts
	var pagesWithStats []models.FaqPageWithStats
	for _, page := range pages {
		// Count groups for this page
		groupCount, err := s.groupRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
		})
		if err != nil {
			groupCount = 0 // Fallback to 0 if error
		}

		// Count questions for this page (count through groups + general questions)
		questionCount := int64(0)

		// Count questions in groups that belong to this page
		groupQuestionsPipeline := []bson.M{
			// First, find all groups for this page
			{
				"$match": bson.M{
					"faqPageId": page.Id,
					"trash":     bson.M{"$ne": true},
				},
			},
			// Then lookup questions for each group
			{
				"$lookup": bson.M{
					"from": "tourismFaqQuestions",
					"let":  bson.M{"groupId": "$_id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$and": []bson.M{
										{"$eq": []interface{}{"$faqGroupId", "$$groupId"}},
										{"$ne": []interface{}{"$trash", true}},
									},
								},
							},
						},
					},
					"as": "questions",
				},
			},
			// Count questions in each group
			{
				"$project": bson.M{
					"questionCount": bson.M{"$size": "$questions"},
				},
			},
			// Sum all question counts
			{
				"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": "$questionCount"},
				},
			},
		}

		var groupQuestionsResult []bson.M
		err = s.groupRepo.Aggregate(ctx, groupQuestionsPipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &groupQuestionsResult)
		})
		if err == nil && len(groupQuestionsResult) > 0 {
			if total, ok := groupQuestionsResult[0]["total"].(int32); ok {
				questionCount += int64(total)
			} else if total, ok := groupQuestionsResult[0]["total"].(int64); ok {
				questionCount += total
			}
		}

		// Count general questions (questions without group but directly assigned to page)
		generalQuestionsCount, err := s.questionRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
			"$or": []bson.M{
				{"faqGroupId": bson.M{"$exists": false}},
				{"faqGroupId": nil},
			},
		})
		if err == nil {
			questionCount += generalQuestionsCount
		}

		pagesWithStats = append(pagesWithStats, models.FaqPageWithStats{
			FaqPage:       page,
			GroupCount:    int(groupCount),
			QuestionCount: int(questionCount),
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

// V2
func (s *faqPageSvcs) GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqPagePagination, error) {
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

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := math.Ceil(float64(count) / float64(limit))

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// Get FAQ Pages
	var pages []models.FaqPage
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &pages)
	})
	if err != nil {
		return nil, err
	}

	// Convert to FaqPageWithStats and calculate counts
	var pagesWithStats []models.FaqPageWithStats
	for _, page := range pages {
		// Count Groups
		groupCount, err := s.groupRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
		})
		if err != nil {
			groupCount = 0
		}

		// Count Questions
		questionCount, err := s.questionRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
		})
		if err != nil {
			questionCount = 0
		}

		pagesWithStats = append(pagesWithStats, models.FaqPageWithStats{
			FaqPage:       page,
			GroupCount:    int(groupCount),
			QuestionCount: int(questionCount),
		})
	}

	return &models.FaqPagePagination{
		FaqPages: pagesWithStats,
		Pagination: common.Pagination{
			TotalPages: totalPages,
			PerPage:    limit,
			TotalCount: count,
		},
	}, nil
}

func (s *faqPageSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqPageWithStats, error) {
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

	var result []models.FaqPage
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	// Convert to FaqPageWithStats and calculate counts
	var pagesWithStats []models.FaqPageWithStats
	for _, page := range result {
		// Count Groups
		groupCount, err := s.groupRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
		})
		if err != nil {
			groupCount = 0
		}

		// Count Questions
		questionCount, err := s.questionRepo.Count(ctx, bson.M{
			"faqPageId": page.Id,
			"trash":     bson.M{"$ne": true},
		})
		if err != nil {
			questionCount = 0
		}

		pagesWithStats = append(pagesWithStats, models.FaqPageWithStats{
			FaqPage:       page,
			GroupCount:    int(groupCount),
			QuestionCount: int(questionCount),
		})
	}
	return pagesWithStats, nil
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

	var pipeline []bson.M
	var countPipeline []bson.M

	// If we have a search term for questions/answers, use a more complex pipeline
	if filter.HasQuestionSearch() {
		searchTerm := filter.GetQuestionSearchTerm()

		// Base match conditions
		baseMatch := filter.ToBsonFilter()

		// Main pipeline with lookup and search
		pipeline = []bson.M{
			{"$match": baseMatch},
			{
				"$lookup": bson.M{
					"from": "tourismFaqQuestions",
					"let":  bson.M{"groupId": "$_id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$and": []bson.M{
										{"$eq": []interface{}{"$faqGroupId", "$$groupId"}},
										{"$ne": []interface{}{"$trash", true}},
									},
								},
								"$or": []bson.M{
									{"question": bson.M{"$regex": searchTerm, "$options": "i"}},
									{"answer": bson.M{"$regex": searchTerm, "$options": "i"}},
								},
							},
						},
					},
					"as": "matchingQuestions",
				},
			},
			{
				"$match": bson.M{
					"matchingQuestions": bson.M{"$ne": []interface{}{}},
				},
			},
			{
				"$project": bson.M{
					"matchingQuestions": 0, // Remove the temporary field
				},
			},
			{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
			{"$skip": skip},
			{"$limit": limit},
		}

		// Count pipeline (same logic but without skip/limit)
		countPipeline = []bson.M{
			{"$match": baseMatch},
			{
				"$lookup": bson.M{
					"from": "tourismFaqQuestions",
					"let":  bson.M{"groupId": "$_id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$and": []bson.M{
										{"$eq": []interface{}{"$faqGroupId", "$$groupId"}},
										{"$ne": []interface{}{"$trash", true}},
									},
								},
								"$or": []bson.M{
									{"question": bson.M{"$regex": searchTerm, "$options": "i"}},
									{"answer": bson.M{"$regex": searchTerm, "$options": "i"}},
								},
							},
						},
					},
					"as": "matchingQuestions",
				},
			},
			{
				"$match": bson.M{
					"matchingQuestions": bson.M{"$ne": []interface{}{}},
				},
			},
			{"$count": "total"},
		}
	} else {
		// Simple pipeline when no search is needed
		filterBody := filter.ToBsonFilter()

		pipeline = []bson.M{
			{"$match": filterBody},
			{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
			{"$skip": skip},
			{"$limit": limit},
		}

		countPipeline = []bson.M{
			{"$match": filterBody},
			{"$count": "total"},
		}
	}

	// Get total count
	var totalCount int64
	if filter.HasQuestionSearch() {
		var countResult []bson.M
		err := s.repo.Aggregate(ctx, countPipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &countResult)
		})
		if err != nil {
			return models.FaqGroupPagination{}, err
		}

		if len(countResult) > 0 {
			if count, ok := countResult[0]["total"].(int32); ok {
				totalCount = int64(count)
			} else if count, ok := countResult[0]["total"].(int64); ok {
				totalCount = count
			}
		}
	} else {
		// Use simple count for non-search queries
		filterBody := filter.ToBsonFilter()
		var err error
		totalCount, err = s.repo.Count(ctx, filterBody)
		if err != nil {
			return models.FaqGroupPagination{}, err
		}
	}

	// Get the groups
	var groups []models.FaqGroup
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

func (s *faqGroupSvcs) SearchQuestions(ctx context.Context, pageId string, searchTerm string, page, size int) (models.FaqSearchResultPagination, error) {
	_pageId, err := primitive.ObjectIDFromHex(pageId)
	if err != nil {
		return models.FaqSearchResultPagination{}, err
	}

	// Set pagination defaults
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// First, find all groups that belong to this page
	groupFilter := bson.M{
		"faqPageId": _pageId,
		"trash":     bson.M{"$ne": true},
	}

	// Get all group IDs for this page using Aggregate
	groupPipeline := []bson.M{
		{"$match": groupFilter},
		{"$project": bson.M{"_id": 1}}, // Only get the IDs
	}

	var groups []models.FaqGroup
	err = s.repo.Aggregate(ctx, groupPipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &groups)
	})
	if err != nil {
		return models.FaqSearchResultPagination{}, err
	}

	// Extract group IDs
	var groupIds []primitive.ObjectID
	for _, group := range groups {
		groupIds = append(groupIds, group.Id)
	}

	// Build search filter for questions
	// Include questions that either:
	// 1. Belong to groups within this page, OR
	// 2. Belong directly to this page (faqPageId) with no group (faqGroupId is null)
	searchFilter := bson.M{
		"trash": bson.M{"$ne": true},
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"faqPageId": _pageId, "faqGroupId": nil}, // Direct page questions (no group)
					{"faqGroupId": bson.M{"$in": groupIds}},   // Questions in groups belonging to this page
				},
			},
			{
				"$or": []bson.M{
					{"question": bson.M{"$regex": searchTerm, "$options": "i"}},
					{"answer": bson.M{"$regex": searchTerm, "$options": "i"}},
				},
			},
		},
	}

	// Count total matching questions
	totalCount, err := s.questionRepo.Count(ctx, searchFilter)
	if err != nil {
		return models.FaqSearchResultPagination{}, err
	}

	// Get matching questions with pagination
	pipeline := []bson.M{
		{"$match": searchFilter},
		{"$sort": bson.M{"sortOrder": 1, "createdAt": 1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var questions []models.FaqQuestion
	err = s.questionRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &questions)
	})
	if err != nil {
		return models.FaqSearchResultPagination{}, err
	}

	// Convert to search results
	var searchResults []models.FaqSearchResult
	for _, q := range questions {
		groupId := ""
		if q.FaqGroupId != nil {
			groupId = q.FaqGroupId.Hex()
		}

		searchResults = append(searchResults, models.FaqSearchResult{
			Id:        q.Id.Hex(),
			PageId:    q.FaqPageId.Hex(),
			GroupId:   groupId,
			Question:  q.Question,
			Answer:    q.Answer,
			IsActive:  q.IsActive,
			SortOrder: q.SortOrder,
			CreatedAt: q.CreatedAt,
			UpdatedAt: q.UpdatedAt,
		})
	}

	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	return models.FaqSearchResultPagination{
		FaqSearchResults: searchResults,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
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

	err := s.repo.Add(ctx, faqGroup)
	if err != nil {
		return err
	}

	// Update parent page timestamp
	err = s.updateParentPageTimestamp(ctx, faqGroup.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent page timestamp: %v\n", err)
	}

	return nil
}

func (s *faqGroupSvcs) AddMany(ctx context.Context, data []models.FaqGroupDto) error {
	if len(data) == 0 {
		return errors.New("empty data array")
	}

	userId := getUserId(ctx)
	var faqGroups []any
	pageIds := make(map[primitive.ObjectID]bool)

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
		// Track unique page IDs
		pageIds[dto.FaqPageId] = true
	}

	err := s.repo.AddMany(ctx, faqGroups)
	if err != nil {
		return err
	}

	// Update parent page timestamps for all affected pages
	for pageId := range pageIds {
		err = s.updateParentPageTimestamp(ctx, pageId)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to update parent page timestamp for page %v: %v\n", pageId, err)
		}
	}

	return nil
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
	if err != nil {
		return err
	}

	// Update parent page timestamp
	err = s.updateParentPageTimestamp(ctx, faqGroup.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent page timestamp: %v\n", err)
	}

	return nil
}

func (s *faqGroupSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing group to access pageId
	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	// Track the page ID for updating parent timestamp
	pageId := existing.FaqPageId

	for key, value := range updates {
		switch key {
		case "faqPageId":
			updateDoc[key] = value
			// Update pageId if it's being changed
			if newPageId, ok := value.(primitive.ObjectID); ok {
				pageId = newPageId
			}
		case "name", "description", "isActive", "sortOrder":
			updateDoc[key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	// Update parent page timestamp
	err = s.updateParentPageTimestamp(ctx, pageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent page timestamp: %v\n", err)
	}

	return nil
}

func (s *faqGroupSvcs) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Get existing group to access pageId
	existing, err := s.GetOne(ctx, id)
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
	if err != nil {
		return err
	}

	// Update parent page timestamp
	err = s.updateParentPageTimestamp(ctx, existing.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent page timestamp: %v\n", err)
	}

	return nil
}

// V2
func (s *faqGroupSvcs) GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqGroupPagination, error) {
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
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.FaqGroup
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}
	pg := common.Pagination{
		TotalPages: math.Ceil(float64(count) / float64(limit)),
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.FaqGroupPagination{
		FaqGroups:  result,
		Pagination: pg,
	}, nil
}

func (s *faqGroupSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqGroup, error) {
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

	var result []models.FaqGroup
	err = s.repo.Aggregate(ctx, pipline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
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

	err := s.repo.Add(ctx, faqQuestion)
	if err != nil {
		return err
	}

	// Update parent group and page timestamps
	err = s.updateParentGroupAndPageTimestamp(ctx, faqQuestion.FaqGroupId, faqQuestion.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent timestamps: %v\n", err)
	}

	return nil
}

func (s *faqQuestionSvcs) AddMany(ctx context.Context, data []models.FaqQuestionDto) error {
	if len(data) == 0 {
		return errors.New("empty data array")
	}

	userId := getUserId(ctx)
	var faqQuestions []any
	pageIds := make(map[primitive.ObjectID]bool)

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
		// Track unique page IDs
		pageIds[dto.FaqPageId] = true
	}

	err := s.repo.AddMany(ctx, faqQuestions)
	if err != nil {
		return err
	}

	// Update parent page timestamps for all affected pages
	for pageId := range pageIds {
		err = s.updateParentPageTimestamp(ctx, pageId)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to update parent page timestamp for page %v: %v\n", pageId, err)
		}
	}

	return nil
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
	if err != nil {
		return err
	}

	// Update parent group and page timestamps
	err = s.updateParentGroupAndPageTimestamp(ctx, faqQuestion.FaqGroupId, faqQuestion.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent timestamps: %v\n", err)
	}

	return nil
}

func (s *faqQuestionSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing question to access pageId
	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	userId := getUserId(ctx)

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	// Track the group and page IDs for updating parent timestamps
	groupId := existing.FaqGroupId
	pageId := existing.FaqPageId

	for key, value := range updates {
		switch key {
		case "faqPageId":
			updateDoc[key] = value
			// Update pageId if it's being changed
			if newPageId, ok := value.(primitive.ObjectID); ok {
				pageId = newPageId
			}
		case "faqGroupId":
			updateDoc[key] = value
			// Update groupId if it's being changed
			if newGroupId, ok := value.(primitive.ObjectID); ok {
				groupId = &newGroupId
			} else if value == nil {
				groupId = nil
			}
		case "question", "answer", "isActive", "sortOrder":
			updateDoc[key] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	// Update parent group and page timestamps
	err = s.updateParentGroupAndPageTimestamp(ctx, groupId, pageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent group and page timestamps: %v\n", err)
	}

	return nil
}

func (s *faqQuestionSvcs) Delete(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Get existing question to access pageId
	existing, err := s.GetOne(ctx, id)
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
	if err != nil {
		return err
	}

	// Update parent page timestamp
	err = s.updateParentPageTimestamp(ctx, existing.FaqPageId)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to update parent page timestamp: %v\n", err)
	}

	return nil
}

func (s *faqQuestionSvcs) GetV2(ctx context.Context, skip int64, limit int64, query *query.Conditions) (*models.FaqQuestionPagination, error) {
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
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.FaqQuestion
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}
	pg := common.Pagination{
		TotalPages: math.Ceil(float64(count) / float64(limit)),
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.FaqQuestionPagination{
		FaqQuestions: result,
		Pagination:   pg,
	}, nil
}

func (s *faqQuestionSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.FaqQuestion, error) {
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

	var result []models.FaqQuestion
	err = s.repo.Aggregate(ctx, pipline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
