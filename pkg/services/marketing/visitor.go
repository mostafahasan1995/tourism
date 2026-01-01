package marketing

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/services/marketing/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"sort"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitorSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Visitor, error)
	Get(ctx context.Context, skip, limit int64, filterQuery filter.VisitorFilter) (*models.VisitorPagination, error)
	GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.VisitorPagination, error)
	Add(ctx context.Context, data *models.VisitorDto) (*models.Visitor, error)
	Update(ctx context.Context, id string, data *models.VisitorDto) (*models.Visitor, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Visitor, error)
	Delete(ctx context.Context, id string) error

	AddTag(ctx context.Context, visitorId string, tagData *models.AddVisitorTagDto) (*models.Visitor, error)
	RemoveTag(ctx context.Context, visitorId string, tag string) (*models.Visitor, error)
	UpdateTags(ctx context.Context, visitorId string, tags *models.UpdateVisitorTagsDto) (*models.Visitor, error)

	AddInterest(ctx context.Context, visitorId string, interestData *models.AddVisitorInterestDto) (*models.Visitor, error)
	RemoveInterest(ctx context.Context, visitorId string, interest string) (*models.Visitor, error)
	UpdateInterests(ctx context.Context, visitorId string, interests *models.UpdateVisitorInterestsDto) (*models.Visitor, error)

	GetStats(ctx context.Context, hotelId *primitive.ObjectID) (*models.VisitorStats, error)

	BulkUpdate(ctx context.Context, bulkData *models.BulkUpdateVisitorsDto) error
	MarkAsVIP(ctx context.Context, visitorId string) (*models.Visitor, error)
	RemoveVIP(ctx context.Context, visitorId string) (*models.Visitor, error)
	ToggleActive(ctx context.Context, visitorId string) (*models.Visitor, error)

	// Visitor Activity methods
	AddActivity(ctx context.Context, activity *models.VisitorActivity) error
	GetVisitorActivities(ctx context.Context, visitorId string, skip, limit int64) (*models.VisitorActivityPagination, error)
}

type visitorSvcs struct {
	repo         repo.VisitorRepo
	activityRepo repo.VisitorActivityRepo
}

func NewVisitorSvcs(i *do.Injector) (VisitorSvcs, error) {
	return &visitorSvcs{
		repo:         do.MustInvoke[repo.VisitorRepo](i),
		activityRepo: do.MustInvoke[repo.VisitorActivityRepo](i),
	}, nil
}

func (s *visitorSvcs) GetOne(ctx context.Context, id string) (*models.Visitor, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *visitorSvcs) Get(ctx context.Context, skip, limit int64, filterQuery filter.VisitorFilter) (*models.VisitorPagination, error) {
	match := filterQuery.ToBsonFilter()

	countPipeline := []bson.M{{"$match": match}}
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"lastVisitDate": -1, "createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var result []models.Visitor
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

	return &models.VisitorPagination{
		Visitors:   result,
		Pagination: pagination,
	}, nil
}

// GetV2 retrieves visitors with pagination and advanced filters
func (s *visitorSvcs) GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.VisitorPagination, error) {
	// Start with base match condition
	match := bson.M{}

	// Convert filter conditions to MongoDB filter
	if filters != nil && len(filters.Columns) > 0 {
		filterBson, err := filters.ConvertToMongo()
		if err != nil {
			return nil, helpers.BadRequest("Invalid filter conditions: " + err.Error())
		}

		// Check if user is explicitly filtering on trash field
		hasTrashFilter := false
		for _, column := range filters.Columns {
			if column.Name == "trash" {
				hasTrashFilter = true
				break
			}
		}

		// If user filters have content
		if len(filterBson) > 0 {
			if hasTrashFilter {
				// User is explicitly filtering trash, so don't add our base condition
				match = filterBson
			} else {
				// User is not filtering trash, so add our base condition to exclude trash
				match = bson.M{
					"$and": []bson.M{
						{"trash": false},
						filterBson,
					},
				}
			}
		} else {
			// No valid filters, use base condition
			match = bson.M{"trash": false}
		}
	} else {
		// No filters provided, use base condition
		match = bson.M{"trash": false}
	}

	pipeline := []bson.M{{"$match": match}}

	count, err := s.repo.Count(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"lastVisitDate": -1, "createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Visitor
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

	return &models.VisitorPagination{
		Visitors:   result,
		Pagination: pagination,
	}, nil
}

func (s *visitorSvcs) Add(ctx context.Context, data *models.VisitorDto) (*models.Visitor, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	if data.VisitDate.IsZero() {
		data.VisitDate = time.Now()
	}

	if data.Interests == nil {
		data.Interests = []string{}
	}
	if data.Tags == nil {
		data.Tags = []string{}
	}

	visitor := &models.Visitor{
		VisitorDto:    *data,
		Id:            primitive.NewObjectID(),
		LastVisitDate: data.VisitDate,
		VisitCount:    1,
		TotalSpent:    0,
		IsActive:      true,
		IsVIP:         false,
		Trash:         false,
		CreatedAt:     time.Now(),
		CreatedBy:     cfg.User.Id,
		UpdatedAt:     time.Now(),
		UpdatedBy:     cfg.User.Id,
	}

	if err := s.repo.Add(ctx, visitor); err != nil {
		return nil, err
	}

	// Log the initial visit activity
	activity := &models.VisitorActivity{
		Id:           primitive.NewObjectID(),
		VisitorId:    visitor.Id,
		HotelId:      visitor.HotelId,
		ActivityType: "visit",
		Description:  "Initial visit recorded",
		CreatedAt:    time.Now(),
	}
	s.activityRepo.Add(ctx, activity)

	return visitor, nil
}

func (s *visitorSvcs) Update(ctx context.Context, id string, data *models.VisitorDto) (*models.Visitor, error) {
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

	visitor := &models.Visitor{
		VisitorDto:    *data,
		Id:            _id,
		LastVisitDate: existing.LastVisitDate,
		VisitCount:    existing.VisitCount,
		TotalSpent:    existing.TotalSpent,
		IsActive:      existing.IsActive,
		IsVIP:         existing.IsVIP,
		Trash:         false,
		CreatedAt:     existing.CreatedAt,
		CreatedBy:     existing.CreatedBy,
		UpdatedBy:     cfg.User.Id,
		UpdatedAt:     time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": visitor}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return visitor, nil
}

func (s *visitorSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Visitor, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}

	for key, value := range updates {
		updateDoc[key] = value
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, id)
}

func (s *visitorSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *visitorSvcs) AddTag(ctx context.Context, visitorId string, tagData *models.AddVisitorTagDto) (*models.Visitor, error) {
	_id, err := primitive.ObjectIDFromHex(visitorId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$addToSet": bson.M{"tags": tagData.Tag},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, visitorId)
}

func (s *visitorSvcs) RemoveTag(ctx context.Context, visitorId string, tag string) (*models.Visitor, error) {
	_id, err := primitive.ObjectIDFromHex(visitorId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$pull": bson.M{"tags": tag},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, visitorId)
}

func (s *visitorSvcs) UpdateTags(ctx context.Context, visitorId string, tags *models.UpdateVisitorTagsDto) (*models.Visitor, error) {
	return s.Patch(ctx, visitorId, map[string]interface{}{"tags": tags.Tags})
}

func (s *visitorSvcs) AddInterest(ctx context.Context, visitorId string, interestData *models.AddVisitorInterestDto) (*models.Visitor, error) {
	_id, err := primitive.ObjectIDFromHex(visitorId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$addToSet": bson.M{"interests": interestData.Interest},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, visitorId)
}

func (s *visitorSvcs) RemoveInterest(ctx context.Context, visitorId string, interest string) (*models.Visitor, error) {
	_id, err := primitive.ObjectIDFromHex(visitorId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$pull": bson.M{"interests": interest},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, visitorId)
}

func (s *visitorSvcs) UpdateInterests(ctx context.Context, visitorId string, interests *models.UpdateVisitorInterestsDto) (*models.Visitor, error) {
	return s.Patch(ctx, visitorId, map[string]interface{}{"interests": interests.Interests})
}

func (s *visitorSvcs) GetStats(ctx context.Context, hotelId *primitive.ObjectID) (*models.VisitorStats, error) {
	match := bson.M{"trash": bson.M{"$ne": true}}
	if hotelId != nil {
		match["hotelId"] = *hotelId
	}

	totalCount, err := s.repo.Count(ctx, []bson.M{{"$match": match}})
	if err != nil {
		return nil, err
	}

	stats := &models.VisitorStats{
		TotalVisitors:        totalCount,
		NationalityBreakdown: make(map[string]int64),
		SourceBreakdown:      make(map[string]int64),
		VisitTrends:          make(map[string]int64),
	}

	if totalCount == 0 {
		return stats, nil
	}

	// Aggregation for statistics
	pipeline := []bson.M{
		{"$match": match},
		{
			"$group": bson.M{
				"_id": nil,
				"activeCount": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$isActive", true}},
							1,
							0,
						},
					},
				},
				"vipCount": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$isVIP", true}},
							1,
							0,
						},
					},
				},
				"nationalities": bson.M{"$push": "$nationality"},
				"sources":       bson.M{"$push": "$source"},
				"interests":     bson.M{"$push": "$interests"},
				"tags":          bson.M{"$push": "$tags"},
				"avgSpent":      bson.M{"$avg": "$totalSpent"},
				"visitCounts":   bson.M{"$push": "$visitCount"},
			},
		},
	}

	var result []struct {
		ActiveCount   int64      `bson:"activeCount"`
		VIPCount      int64      `bson:"vipCount"`
		Nationalities []string   `bson:"nationalities"`
		Sources       []string   `bson:"sources"`
		Interests     [][]string `bson:"interests"`
		Tags          [][]string `bson:"tags"`
		AvgSpent      float64    `bson:"avgSpent"`
		VisitCounts   []int      `bson:"visitCounts"`
	}

	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		data := result[0]
		stats.ActiveVisitors = data.ActiveCount
		stats.VIPVisitors = data.VIPCount
		stats.AverageSpent = data.AvgSpent

		// Count nationalities
		for _, nationality := range data.Nationalities {
			if nationality != "" {
				stats.NationalityBreakdown[nationality]++
			}
		}

		// Count sources
		for _, source := range data.Sources {
			if source != "" {
				stats.SourceBreakdown[source]++
			}
		}

		// Get popular interests (flatten and count)
		interestCount := make(map[string]int)
		for _, interestList := range data.Interests {
			for _, interest := range interestList {
				if interest != "" {
					interestCount[interest]++
				}
			}
		}
		stats.PopularInterests = getTopStrings(interestCount, 10)

		// Get popular tags (flatten and count)
		tagCount := make(map[string]int)
		for _, tagList := range data.Tags {
			for _, tag := range tagList {
				if tag != "" {
					tagCount[tag]++
				}
			}
		}
		stats.PopularTags = getTopStrings(tagCount, 10)

		// Calculate return visitor rate
		returnVisitors := 0
		for _, visitCount := range data.VisitCounts {
			if visitCount > 1 {
				returnVisitors++
			}
		}
		if totalCount > 0 {
			stats.ReturnVisitorRate = float64(returnVisitors) / float64(totalCount) * 100
		}
	}

	return stats, nil
}

func getTopStrings(countMap map[string]int, limit int) []string {
	type item struct {
		name  string
		count int
	}

	items := make([]item, 0, len(countMap))
	for name, count := range countMap {
		items = append(items, item{name: name, count: count})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})

	result := make([]string, 0, limit)
	for i := 0; i < len(items) && i < limit; i++ {
		result = append(result, items[i].name)
	}

	return result
}

func (s *visitorSvcs) BulkUpdate(ctx context.Context, bulkData *models.BulkUpdateVisitorsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}

	if bulkData.Tags != nil {
		updateDoc["tags"] = *bulkData.Tags
	}
	if bulkData.Interests != nil {
		updateDoc["interests"] = *bulkData.Interests
	}
	if bulkData.IsVIP != nil {
		updateDoc["isVIP"] = *bulkData.IsVIP
	}
	if bulkData.IsActive != nil {
		updateDoc["isActive"] = *bulkData.IsActive
	}

	filter := bson.M{"_id": bson.M{"$in": bulkData.VisitorIds}}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *visitorSvcs) MarkAsVIP(ctx context.Context, visitorId string) (*models.Visitor, error) {
	return s.Patch(ctx, visitorId, map[string]interface{}{"isVIP": true})
}

func (s *visitorSvcs) RemoveVIP(ctx context.Context, visitorId string) (*models.Visitor, error) {
	return s.Patch(ctx, visitorId, map[string]interface{}{"isVIP": false})
}

func (s *visitorSvcs) ToggleActive(ctx context.Context, visitorId string) (*models.Visitor, error) {
	visitor, err := s.GetOne(ctx, visitorId)
	if err != nil {
		return nil, err
	}

	return s.Patch(ctx, visitorId, map[string]interface{}{"isActive": !visitor.IsActive})
}

func (s *visitorSvcs) AddActivity(ctx context.Context, activity *models.VisitorActivity) error {
	if activity.Id.IsZero() {
		activity.Id = primitive.NewObjectID()
	}
	if activity.CreatedAt.IsZero() {
		activity.CreatedAt = time.Now()
	}

	return s.activityRepo.Add(ctx, activity)
}

func (s *visitorSvcs) GetVisitorActivities(ctx context.Context, visitorId string, skip, limit int64) (*models.VisitorActivityPagination, error) {
	_visitorId, err := primitive.ObjectIDFromHex(visitorId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	match := bson.M{"visitorId": _visitorId}

	countPipeline := []bson.M{{"$match": match}}
	count, err := s.activityRepo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var result []models.VisitorActivity
	errAg := s.activityRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.VisitorActivityPagination{
		Activities: result,
		Pagination: pagination,
	}, nil
}
