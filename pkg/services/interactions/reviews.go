package interactions

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/interactions/filter"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/interactions/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Review, error)
	GetAll(ctx context.Context, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error)
	GetStats(ctx context.Context) (*models.ReviewStats, error)
	GetEntityReviews(ctx context.Context, entityType string, refId string, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error)
	GetEntityStats(ctx context.Context, entityType string, refId string) (*models.EntityReviewSummary, error)
	Add(ctx context.Context, data *models.ReviewDto) (*models.Review, error)
	Update(ctx context.Context, id string, data *models.ReviewDto) error
	Patch(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	AddReply(ctx context.Context, reviewId string, reply models.ReviewReply) error
	UpdateReviewStatus(ctx context.Context, reviewId string, status string) error
	ApproveReview(ctx context.Context, reviewId string) error
	RejectReview(ctx context.Context, reviewId string) error
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
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *reviewsSvcs) GetAll(ctx context.Context, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error) {
	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	return s.repo.GetAllByFilter(ctx, filter, page, perPage)
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

func (s *reviewsSvcs) GetEntityReviews(ctx context.Context, entityType string, refId string, filter filter.ReviewsFilter, page, perPage int) (models.ReviewPagination, error) {
	// Convert refId to ObjectID
	_refId, err := primitive.ObjectIDFromHex(refId)
	if err != nil {
		return models.ReviewPagination{}, helpers.InvalidObjectId()
	}

	// Set the entity filter
	filter.Type = entityType
	filter.Ref = _refId

	// Set default pagination values
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	return s.repo.GetAllByFilter(ctx, filter, page, perPage)
}

func (s *reviewsSvcs) GetEntityStats(ctx context.Context, entityType string, refId string) (*models.EntityReviewSummary, error) {
	// Convert refId to ObjectID
	_refId, err := primitive.ObjectIDFromHex(refId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	// Get total count of approved reviews for this entity
	totalCount, err := s.repo.Count(ctx, bson.M{
		"trash":  bson.M{"$ne": true},
		"status": "approved",
		"type":   entityType,
		"ref":    _refId,
	})
	if err != nil {
		return nil, err
	}

	// Initialize summary with default values
	summary := &models.EntityReviewSummary{
		Type:            entityType,
		RefId:           _refId,
		TotalReviews:    totalCount,
		AverageRating:   0,
		RatingBreakdown: make(map[int]int64),
	}

	// Initialize rating breakdown
	for i := 1; i <= 5; i++ {
		summary.RatingBreakdown[i] = 0
	}

	// If no reviews for this entity, return empty stats
	if totalCount == 0 {
		return summary, nil
	}

	// Calculate average rating and rating breakdown for this entity
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"trash":  bson.M{"$ne": true},
				"status": "approved",
				"type":   entityType,
				"ref":    _refId,
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

func (s *reviewsSvcs) Add(ctx context.Context, data *models.ReviewDto) (*models.Review, error) {
	// Handle user information based on authentication status
	var createdBy primitive.ObjectID

	// Check if user is authenticated by getting user from context
	user, _ := ctx.Value(util.ReqUser).(*types.User)
	if user == nil {
		// Unauthenticated user - require personal info fields for some review types
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return nil, helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
		createdBy = primitive.NilObjectID
		data.UserId = "000000000000000000000000" // Default ObjectID for unauthenticated users
	} else {
		// Authenticated user - use their information
		createdBy = user.Id
		data.UserId = user.Id.Hex()
		// For authenticated users, clear personal info fields (privacy)
		data.FirstName = ""
		data.LastName = ""
		data.Email = ""
	}

	// Set default status if not provided
	if data.Status == "" {
		data.Status = "pending"
	}

	// Set date if not provided
	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	// Initialize replies and images if nil
	if data.Replies == nil {
		data.Replies = []models.ReviewReply{}
	}
	if data.Images == nil {
		data.Images = []models.ReviewImage{}
	}

	review := &models.Review{
		ReviewDto: *data,
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
		UpdatedAt: time.Now(),
		UpdatedBy: createdBy,
	}

	if err := s.repo.Add(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
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
		case "status", "username", "userId", "value", "description", "adviceForTravelers":
			updateDoc[key] = value
		case "userImg", "images", "replies":
			updateDoc[key] = value
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

	// First, ensure replies field is an array (initialize if null)
	filter := bson.M{"_id": _id, "replies": bson.M{"$type": "null"}}
	initUpdate := bson.M{"$set": bson.M{"replies": []models.ReviewReply{}}}

	// This will only update documents where replies is null
	s.repo.Patch(ctx, filter, initUpdate)

	// Now add the reply using $push (replies is guaranteed to be an array)
	filter = bson.M{"_id": _id}
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

func (s *reviewsSvcs) UpdateReviewStatus(ctx context.Context, reviewId string, status string) error {
	return s.repo.UpdateReviewStatus(ctx, reviewId, status)
}

func (s *reviewsSvcs) ApproveReview(ctx context.Context, reviewId string) error {
	return s.UpdateReviewStatus(ctx, reviewId, "approved")
}

func (s *reviewsSvcs) RejectReview(ctx context.Context, reviewId string) error {
	return s.UpdateReviewStatus(ctx, reviewId, "rejected")
}
