package home

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/services/home/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Review, error)
	GetAll(ctx context.Context, filter filter.ReviewsFilter) (models.ReviewPagination, error)
	GetStats(ctx context.Context) (*models.ReviewStats, error)
	GetProgramReviews(ctx context.Context, programId string, filter filter.ReviewsFilter) (models.ReviewPagination, error)
	GetProgramStats(ctx context.Context, programId string) (*models.ProgramReviewSummary, error)
	Add(ctx context.Context, data *models.ReviewDto) error
	AddMany(ctx context.Context, data []models.ReviewDto) error
	Update(ctx context.Context, id string, data *models.ReviewDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	AddReply(ctx context.Context, reviewId string, reply models.ReviewReply) error
	UpdateReplyStatus(ctx context.Context, id string, status string) error
}

type reviewsSvcs struct {
	repo repo.ReviewsRepo
}

func NewReviewsSvcs(i *do.Injector) (ReviewsSvcs, error) {
	return &reviewsSvcs{
		repo: do.MustInvoke[repo.ReviewsRepo](i),
	}, nil
}

func (s *reviewsSvcs) GetOne(ctx context.Context, id string) (*models.Review, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *reviewsSvcs) GetAll(ctx context.Context, filter filter.ReviewsFilter) (models.ReviewPagination, error) {
	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := s.repo.Count(ctx, filterBody)
	if err != nil {
		return models.ReviewPagination{}, err
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
		{"$sort": bson.M{"date": -1, "createdAt": -1}}, // Sort by review date, newest first
		{"$skip": skip},
		{"$limit": limit},
	}

	var reviews []models.Review
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &reviews)
	})
	if err != nil {
		return models.ReviewPagination{}, err
	}

	// Prepare pagination result
	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	result := models.ReviewPagination{
		Reviews: reviews,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (s *reviewsSvcs) GetStats(ctx context.Context) (*models.ReviewStats, error) {
	// Get total count of approved reviews
	totalCount, err := s.repo.Count(ctx, bson.M{
		"trash":  bson.M{"$ne": true},
		"status": "approved",
	})
	if err != nil {
		return nil, err
	}

	// Initialize stats with default values
	stats := &models.ReviewStats{
		TotalReviews:      totalCount,
		AverageRating:     0,
		RatingBreakdown:   make(map[int]int64),
		TopDestinations:   []string{},
		TopCountries:      []string{},
		ReviewsWithImages: 0,
	}

	// Initialize rating breakdown
	for i := 1; i <= 5; i++ {
		stats.RatingBreakdown[i] = 0
	}

	// If no reviews, return empty stats
	if totalCount == 0 {
		return stats, nil
	}

	// Calculate average rating, rating breakdown, and additional stats
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"trash":  bson.M{"$ne": true},
				"status": "approved",
			},
		},
		{
			"$group": bson.M{
				"_id":           nil,
				"averageRating": bson.M{"$avg": "$value"},
				"ratings":       bson.M{"$push": "$value"},
				"destinations":  bson.M{"$push": "$destination"},
				"countries":     bson.M{"$push": "$countries"},
				"reviewsWithImages": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$gt": []interface{}{bson.M{"$size": bson.M{"$ifNull": []interface{}{"$images", []interface{}{}}}}, 0}},
							1,
							0,
						},
					},
				},
			},
		},
	}

	var result []struct {
		AverageRating     float64    `bson:"averageRating"`
		Ratings           []float64  `bson:"ratings"`
		Destinations      []string   `bson:"destinations"`
		Countries         [][]string `bson:"countries"`
		ReviewsWithImages int64      `bson:"reviewsWithImages"`
	}

	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		stats.AverageRating = result[0].AverageRating
		stats.ReviewsWithImages = result[0].ReviewsWithImages

		// Calculate rating breakdown
		for _, rating := range result[0].Ratings {
			ratingInt := int(rating)
			if ratingInt >= 1 && ratingInt <= 5 {
				stats.RatingBreakdown[ratingInt]++
			}
		}

		// Process destinations
		destinationCount := make(map[string]int)
		for _, dest := range result[0].Destinations {
			if dest != "" {
				destinationCount[dest]++
			}
		}

		// Process countries (flatten the array of arrays)
		countryCount := make(map[string]int)
		for _, countryArray := range result[0].Countries {
			for _, country := range countryArray {
				if country != "" {
					countryCount[country]++
				}
			}
		}

		// Get top 5 destinations and countries
		stats.TopDestinations = getTopItems(destinationCount, 5)
		stats.TopCountries = getTopItems(countryCount, 5)
	}

	return stats, nil
}

// Helper function to get top items from a count map
func getTopItems(countMap map[string]int, limit int) []string {
	type item struct {
		name  string
		count int
	}

	items := make([]item, 0, len(countMap))
	for name, count := range countMap {
		items = append(items, item{name: name, count: count})
	}

	// Simple bubble sort for small datasets
	for i := 0; i < len(items)-1; i++ {
		for j := 0; j < len(items)-i-1; j++ {
			if items[j].count < items[j+1].count {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}

	result := make([]string, 0, limit)
	for i := 0; i < len(items) && i < limit; i++ {
		result = append(result, items[i].name)
	}

	return result
}

func (s *reviewsSvcs) GetProgramReviews(ctx context.Context, programId string, filter filter.ReviewsFilter) (models.ReviewPagination, error) {
	// Convert programId to ObjectID
	_programId, err := primitive.ObjectIDFromHex(programId)
	if err != nil {
		return models.ReviewPagination{}, helpers.InvalidObjectId()
	}

	// Set the program filter
	filter.ProgramId = _programId

	// Use the existing GetAll method with the program filter
	return s.GetAll(ctx, filter)
}

func (s *reviewsSvcs) GetProgramStats(ctx context.Context, programId string) (*models.ProgramReviewSummary, error) {
	// Convert programId to ObjectID
	_programId, err := primitive.ObjectIDFromHex(programId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get total count of approved reviews for this program
	totalCount, err := s.repo.Count(ctx, bson.M{
		"trash":     bson.M{"$ne": true},
		"status":    "approved",
		"programId": _programId,
	})
	if err != nil {
		return nil, err
	}

	// Initialize summary with default values
	summary := &models.ProgramReviewSummary{
		ProgramId:       _programId,
		TotalReviews:    totalCount,
		AverageRating:   0,
		RatingBreakdown: make(map[int]int64),
	}

	// Initialize rating breakdown
	for i := 1; i <= 5; i++ {
		summary.RatingBreakdown[i] = 0
	}

	// If no reviews for this program, return empty stats
	if totalCount == 0 {
		return summary, nil
	}

	// Calculate average rating and rating breakdown for this program
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"trash":     bson.M{"$ne": true},
				"status":    "approved",
				"programId": _programId,
			},
		},
		{
			"$group": bson.M{
				"_id":           nil,
				"averageRating": bson.M{"$avg": "$value"},
				"ratings":       bson.M{"$push": "$value"},
			},
		},
	}

	var result []struct {
		AverageRating float64   `bson:"averageRating"`
		Ratings       []float64 `bson:"ratings"`
	}

	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		summary.AverageRating = result[0].AverageRating

		// Calculate rating breakdown
		for _, rating := range result[0].Ratings {
			ratingInt := int(rating)
			if ratingInt >= 1 && ratingInt <= 5 {
				summary.RatingBreakdown[ratingInt]++
			}
		}
	}

	return summary, nil
}

func (s *reviewsSvcs) Add(ctx context.Context, data *models.ReviewDto) error {
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

	// Set default status if not provided
	if data.Status == "" {
		data.Status = "pending"
	}

	// Set date if not provided
	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	// Initialize replies if nil
	if data.Replies == nil {
		data.Replies = []models.ReviewReply{}
	}

	review := &models.Review{
		ReviewDto: *data,
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: userId,
		UpdatedAt: time.Now(),
		UpdatedBy: userId,
	}

	if err := s.repo.Add(ctx, review); err != nil {
		return err
	}

	return nil
}

func (s *reviewsSvcs) AddMany(ctx context.Context, data []models.ReviewDto) error {
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

	var reviewsArray []any
	for _, reviewDto := range data {
		// Set default status if not provided
		if reviewDto.Status == "" {
			reviewDto.Status = "pending"
		}

		// Set date if not provided
		if reviewDto.Date.IsZero() {
			reviewDto.Date = time.Now()
		}

		// Initialize replies if nil
		if reviewDto.Replies == nil {
			reviewDto.Replies = []models.ReviewReply{}
		}

		review := &models.Review{
			ReviewDto: reviewDto,
			Id:        primitive.NewObjectID(),
			Trash:     false,
			CreatedAt: time.Now(),
			CreatedBy: userId,
			UpdatedAt: time.Now(),
			UpdatedBy: userId,
		}
		reviewsArray = append(reviewsArray, review)
	}

	err = s.repo.AddMany(ctx, reviewsArray)
	if err != nil {
		return err
	}

	return nil
}

func (s *reviewsSvcs) Update(ctx context.Context, id string, data *models.ReviewDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Get existing review to preserve created fields
	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}

	// Set date if not provided
	if data.Date.IsZero() {
		data.Date = existing.Date // Preserve original review date
	}

	review := &models.Review{
		ReviewDto: *data,
		Id:        _id,
		Trash:     false,
		CreatedAt: existing.CreatedAt,
		CreatedBy: existing.CreatedBy,
		UpdatedBy: userId,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": review}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *reviewsSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Handle case where user is not authenticated
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}

	// Build the update document with only the fields that are being updated
	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

	// Add the specific fields to update
	for key, value := range updates {
		switch key {
		case "status", "username", "userId", "value", "text":
			updateDoc[key] = value
		case "userImg":
			updateDoc["userImg"] = value
		case "replies":
			updateDoc["replies"] = value
		}
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *reviewsSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Handle case where user is not authenticated
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}

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

	return nil
}

func (s *reviewsSvcs) AddReply(ctx context.Context, reviewId string, reply models.ReviewReply) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Validate reply text
	if reply.Text == "" {
		return helpers.BadRequest("Reply text is required")
	}

	// Check if review exists
	existingReview, err := s.GetOne(ctx, reviewId)
	if err != nil {
		return err
	}
	if existingReview == nil {
		return errors.New("review not found")
	}

	// Handle case where user is not authenticated
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}

	// Set reply metadata
	reply.Date = time.Now()

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$push": bson.M{"replies": reply},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": userId,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *reviewsSvcs) UpdateReplyStatus(ctx context.Context, id string, status string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	// Handle case where user is not authenticated
	var userId primitive.ObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	} else {
		userId = primitive.NilObjectID
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
