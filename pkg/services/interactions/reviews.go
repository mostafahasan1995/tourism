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
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Review, error)
	Get(ctx context.Context, skip, limit int64, query any) (*models.ReviewPagination, error)
	GetAllApproved(ctx context.Context) ([]models.Review, error)
	GetAllWithPagination(ctx context.Context, skip, limit int64, query any) (*models.ReviewPagination, error)
	GetAllWithoutPagination(ctx context.Context, query any) ([]models.Review, error)
	GetStats(ctx context.Context) (*models.ReviewStats, error)
	Add(ctx context.Context, data *models.ReviewDto) (*models.Review, error)
	AddFromDashboard(ctx context.Context, data *models.ReviewDto) (*models.Review, error)
	Update(ctx context.Context, id string, data *models.ReviewDto) (*models.Review, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Review, error)
	Delete(ctx context.Context, id string) error
	AddReply(ctx context.Context, reviewId string, reply models.ReviewReply) error
	UpdateReviewStatus(ctx context.Context, reviewId string, status string) (*models.Review, error)
	ApproveReview(ctx context.Context, reviewId string) (*models.Review, error)
	RejectReview(ctx context.Context, reviewId string) (*models.Review, error)
	Count(ctx context.Context, filter any) (int64, error)
}

type reviewsSvcs struct {
	repo repo.ReviewsRepo
}

func NewReviewsSvcs(i *do.Injector) (ReviewsSvcs, error) {
	return &reviewsSvcs{
		repo: do.MustInvoke[repo.ReviewsRepo](i),
	}, nil
}

func (s *reviewsSvcs) GetAllApproved(ctx context.Context) ([]models.Review, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false, "status": "approved"}},
		{"$sort": bson.M{"date": -1, "createdAt": -1}},
	}

	var result []models.Review
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *reviewsSvcs) GetAllWithoutPagination(ctx context.Context, query any) ([]models.Review, error) {
	match := bson.M{}

	filters, err := helpers.ParseFilters[filter.ReviewsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	// Don't filter by status - get all reviews regardless of status
	pipeline := filters.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"date": -1, "createdAt": -1}})

	var result []models.Review
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *reviewsSvcs) GetAllWithPagination(ctx context.Context, skip, limit int64, query any) (*models.ReviewPagination, error) {
	match := bson.M{}

	filters, err := helpers.ParseFilters[filter.ReviewsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	// Don't filter by status - get all reviews regardless of status
	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"date": -1, "createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Review
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.ReviewPagination{
		Reviews:    result,
		Pagination: pagination,
	}, nil
}

func (s *reviewsSvcs) GetOne(ctx context.Context, id string) (*models.Review, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *reviewsSvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.ReviewPagination, error) {
	match := bson.M{}

	filters, err := helpers.ParseFilters[filter.ReviewsFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"date": -1, "createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Review
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.ReviewPagination{
		Reviews:    result,
		Pagination: pagination,
	}, nil
}

func (s *reviewsSvcs) GetStats(ctx context.Context) (*models.ReviewStats, error) {
	match := bson.M{
		"trash":  bson.M{"$ne": true},
		"status": "approved",
	}

	countPipeline := []bson.M{{"$match": match}}
	totalCount, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	stats := &models.ReviewStats{
		TotalReviews:      totalCount,
		AverageRating:     0,
		RatingBreakdown:   make(map[int]int64),
		TopDestinations:   []string{},
		TopCountries:      []string{},
		ReviewsWithImages: 0,
	}

	for i := 1; i <= 5; i++ {
		stats.RatingBreakdown[i] = 0
	}

	if totalCount == 0 {
		return stats, nil
	}

	pipeline := []bson.M{
		{"$match": match},
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

		for _, rating := range result[0].Ratings {
			ratingInt := int(rating)
			if ratingInt >= 1 && ratingInt <= 5 {
				stats.RatingBreakdown[ratingInt]++
			}
		}

		destinationCount := make(map[string]int)
		for _, dest := range result[0].Destinations {
			if dest != "" {
				destinationCount[dest]++
			}
		}

		countryCount := make(map[string]int)
		for _, countryArray := range result[0].Countries {
			for _, country := range countryArray {
				if country != "" {
					countryCount[country]++
				}
			}
		}

		stats.TopDestinations = getTopItems(destinationCount, 5)
		stats.TopCountries = getTopItems(countryCount, 5)
	}

	return stats, nil
}

func getTopItems(countMap map[string]int, limit int) []string {
	type item struct {
		name  string
		count int
	}

	items := make([]item, 0, len(countMap))
	for name, count := range countMap {
		items = append(items, item{name: name, count: count})
	}

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

func (s *reviewsSvcs) Add(ctx context.Context, data *models.ReviewDto) (*models.Review, error) {
	var createdBy primitive.ObjectID

	user, _ := ctx.Value(util.ReqUser).(*types.User)
	if user == nil {
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return nil, helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
		createdBy = primitive.NilObjectID
		data.UserId = "000000000000000000000000"
		data.UserImg = nil
	} else {
		createdBy = user.Id
		data.UserId = user.Id.Hex()
		data.FirstName = ""
		data.LastName = ""
		data.Email = ""

		if user.UserData.Picture != "" {
			data.UserImg = &types.FileField{
				Path: user.UserData.Picture,
			}
		} else {
			data.UserImg = nil
		}
	}

	if data.Status == "" {
		data.Status = "pending"
	}

	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	if data.Replies == nil {
		data.Replies = []models.ReviewReply{}
	}
	if data.Images == nil {
		data.Images = []types.FileField{}
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

func (s *reviewsSvcs) AddFromDashboard(ctx context.Context, data *models.ReviewDto) (*models.Review, error) {
	var createdBy primitive.ObjectID

	user, _ := ctx.Value(util.ReqUser).(*types.User)
	if user == nil {
		return nil, helpers.BadRequest("authentication required for dashboard reviews")
	}

	// If userId is provided in the body, use it, otherwise use the authenticated user's ID
	if data.UserId == "" {
		data.UserId = user.Id.Hex()
	}

	// Get user ID as ObjectID
	userObjectId, err := primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		return nil, helpers.BadRequest("invalid userId format")
	}
	createdBy = userObjectId

	// Validate countries for destination type
	if data.Type == "destination" && (data.Countries == nil || len(data.Countries) == 0) {
		return nil, helpers.BadRequest("countries are required for destination reviews")
	}

	if data.Status == "" {
		data.Status = "pending"
	}

	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	if data.Replies == nil {
		data.Replies = []models.ReviewReply{}
	}
	if data.Images == nil {
		data.Images = []types.FileField{}
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

func (s *reviewsSvcs) Update(ctx context.Context, id string, data *models.ReviewDto) (*models.Review, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	userId := cfg.User.Id

	if data.Date.IsZero() {
		data.Date = existing.Date
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
		return nil, err
	}

	return review, nil
}

func (s *reviewsSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Review, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	userId := cfg.User.Id

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}

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
		return nil, err
	}

	return s.GetOne(ctx, id)
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

	userId := cfg.User.Id

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

	if reply.Text == "" {
		return helpers.BadRequest("Reply text is required")
	}

	existingReview, err := s.GetOne(ctx, reviewId)
	if err != nil {
		return err
	}
	if existingReview == nil {
		return errors.New("review not found")
	}

	userId := cfg.User.Id

	reply.Date = time.Now()

	filter := bson.M{"_id": _id, "replies": bson.M{"$type": "null"}}
	initUpdate := bson.M{"$set": bson.M{"replies": []models.ReviewReply{}}}

	s.repo.Patch(ctx, filter, initUpdate)

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

func (s *reviewsSvcs) UpdateReviewStatus(ctx context.Context, reviewId string, status string) (*models.Review, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	userId := cfg.User.Id

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now(),
		"updatedBy": userId,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, reviewId)
}

func (s *reviewsSvcs) ApproveReview(ctx context.Context, reviewId string) (*models.Review, error) {
	return s.UpdateReviewStatus(ctx, reviewId, "approved")
}

func (s *reviewsSvcs) RejectReview(ctx context.Context, reviewId string) (*models.Review, error) {
	return s.UpdateReviewStatus(ctx, reviewId, "rejected")
}

func (s *reviewsSvcs) Count(ctx context.Context, filter any) (int64, error) {
	return s.repo.Count(ctx, filter)
}
