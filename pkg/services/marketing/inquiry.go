package marketing

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/services/marketing/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InquirySvcs interface {
	GetOne(ctx context.Context, id string) (*models.Inquiry, error)
	Get(ctx context.Context, skip, limit int64, filterQuery filter.InquiryFilter) (*models.InquiryPagination, error)
	Add(ctx context.Context, data *models.InquiryDto) (*models.Inquiry, error)
	Update(ctx context.Context, id string, data *models.InquiryDto) (*models.Inquiry, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Inquiry, error)
	Delete(ctx context.Context, id string) error

	AddReply(ctx context.Context, inquiryId string, reply *models.AddInquiryReplyDto) (*models.Inquiry, error)
	UpdateStatus(ctx context.Context, inquiryId string, status *models.UpdateInquiryStatusDto) (*models.Inquiry, error)
	MarkAsRead(ctx context.Context, inquiryId string, markRead *models.MarkAsReadDto) (*models.Inquiry, error)

	GetStats(ctx context.Context, hotelId *primitive.ObjectID) (*models.InquiryStats, error)
	GetUnreadCount(ctx context.Context, hotelId *primitive.ObjectID) (int64, error)

	//BulkUpdateStatus(ctx context.Context, inquiryIds []string, status models.InquiryStatus) error
	//BulkMarkAsRead(ctx context.Context, inquiryIds []string, isRead bool) error
}

type inquirySvcs struct {
	repo repo.InquiryRepo
}

func NewInquirySvcs(i *do.Injector) (InquirySvcs, error) {
	return &inquirySvcs{
		repo: do.MustInvoke[repo.InquiryRepo](i),
	}, nil
}

func (s *inquirySvcs) GetOne(ctx context.Context, id string) (*models.Inquiry, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *inquirySvcs) Get(ctx context.Context, skip, limit int64, filterQuery filter.InquiryFilter) (*models.InquiryPagination, error) {
	match := filterQuery.ToBsonFilter()

	countPipeline := []bson.M{{"$match": match}}
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var result []models.Inquiry
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

	return &models.InquiryPagination{
		Inquiries:  result,
		Pagination: pagination,
	}, nil
}

func (s *inquirySvcs) Add(ctx context.Context, data *models.InquiryDto) (*models.Inquiry, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var createdBy primitive.ObjectID
	if cfg.User != nil {
		createdBy = cfg.User.Id
	} else {
		createdBy = primitive.NilObjectID
	}

	if data.Priority == "" {
		data.Priority = "medium"
	}

	if data.Source == "" {
		data.Source = "chat"
	}

	inquiry := &models.Inquiry{
		InquiryDto: *data,
		Id:         primitive.NewObjectID(),
		Status:     models.InquiryStatusReceived,
		Replies:    []models.InquiryReply{},
		IsRead:     false,
		Trash:      false,
		CreatedAt:  time.Now(),
		CreatedBy:  createdBy,
		UpdatedAt:  time.Now(),
		UpdatedBy:  createdBy,
	}

	if err := s.repo.Add(ctx, inquiry); err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *inquirySvcs) Update(ctx context.Context, id string, data *models.InquiryDto) (*models.Inquiry, error) {
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

	inquiry := &models.Inquiry{
		InquiryDto: *data,
		Id:         _id,
		Status:     existing.Status,
		Replies:    existing.Replies,
		IsRead:     existing.IsRead,
		ReadAt:     existing.ReadAt,
		ReadBy:     existing.ReadBy,
		Trash:      false,
		CreatedAt:  existing.CreatedAt,
		CreatedBy:  existing.CreatedBy,
		UpdatedBy:  cfg.User.Id,
		UpdatedAt:  time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": inquiry}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *inquirySvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Inquiry, error) {
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

func (s *inquirySvcs) Delete(ctx context.Context, id string) error {
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

func (s *inquirySvcs) AddReply(ctx context.Context, inquiryId string, replyData *models.AddInquiryReplyDto) (*models.Inquiry, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(inquiryId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	inquiry, err := s.GetOne(ctx, inquiryId)
	if err != nil {
		return nil, err
	}

	reply := models.InquiryReply{
		Id:              primitive.NewObjectID(),
		VisitorName:     inquiry.VisitorName,
		UserId:          &cfg.User.Id,
		ResponseMessage: replyData.ResponseMessage,
		Attachments:     replyData.Attachments,
		FollowUpNotes:   replyData.FollowUpNotes,
		Status:          replyData.Status,
		CreatedAt:       time.Now(),
		CreatedBy:       cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$push": bson.M{"replies": reply},
		"$set": bson.M{
			"status":    replyData.Status,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, inquiryId)
}

func (s *inquirySvcs) UpdateStatus(ctx context.Context, inquiryId string, statusUpdate *models.UpdateInquiryStatusDto) (*models.Inquiry, error) {
	return s.Patch(ctx, inquiryId, map[string]interface{}{"status": statusUpdate.Status})
}

func (s *inquirySvcs) MarkAsRead(ctx context.Context, inquiryId string, markRead *models.MarkAsReadDto) (*models.Inquiry, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"isRead": markRead.IsRead,
	}

	if markRead.IsRead {
		now := time.Now()
		updates["readAt"] = &now
		updates["readBy"] = &cfg.User.Id
	} else {
		updates["readAt"] = nil
		updates["readBy"] = nil
	}

	return s.Patch(ctx, inquiryId, updates)
}

func (s *inquirySvcs) GetStats(ctx context.Context, hotelId *primitive.ObjectID) (*models.InquiryStats, error) {
	match := bson.M{"trash": bson.M{"$ne": true}}
	if hotelId != nil {
		match["hotelId"] = *hotelId
	}

	totalCount, err := s.repo.Count(ctx, []bson.M{{"$match": match}})
	if err != nil {
		return nil, err
	}

	stats := &models.InquiryStats{
		TotalInquiries:    totalCount,
		StatusBreakdown:   make(map[string]int64),
		PriorityBreakdown: make(map[string]int64),
		SourceBreakdown:   make(map[string]int64),
		WeeklyTrend:       make([]int64, 7),
	}

	if totalCount == 0 {
		return stats, nil
	}

	// Aggregation for various statistics
	pipeline := []bson.M{
		{"$match": match},
		{
			"$group": bson.M{
				"_id": nil,
				"statusBreakdown": bson.M{
					"$push": "$status",
				},
				"priorityBreakdown": bson.M{
					"$push": "$priority",
				},
				"sourceBreakdown": bson.M{
					"$push": "$source",
				},
				"unreadCount": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$isRead", false}},
							1,
							0,
						},
					},
				},
				"todayCount": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$gte": []interface{}{"$createdAt", time.Now().AddDate(0, 0, -1)}},
							1,
							0,
						},
					},
				},
			},
		},
	}

	var result []struct {
		StatusBreakdown   []string `bson:"statusBreakdown"`
		PriorityBreakdown []string `bson:"priorityBreakdown"`
		SourceBreakdown   []string `bson:"sourceBreakdown"`
		UnreadCount       int64    `bson:"unreadCount"`
		TodayCount        int64    `bson:"todayCount"`
	}

	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		stats.UnreadCount = result[0].UnreadCount
		stats.TodayInquiries = result[0].TodayCount

		// Count occurrences for breakdowns
		for _, status := range result[0].StatusBreakdown {
			stats.StatusBreakdown[status]++
		}
		for _, priority := range result[0].PriorityBreakdown {
			stats.PriorityBreakdown[priority]++
		}
		for _, source := range result[0].SourceBreakdown {
			stats.SourceBreakdown[source]++
		}
	}

	// Get weekly trend (last 7 days)
	for i := 0; i < 7; i++ {
		dayStart := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		dayEnd := dayStart.Add(24 * time.Hour)

		dayMatch := bson.M{
			"trash": bson.M{"$ne": true},
			"createdAt": bson.M{
				"$gte": dayStart,
				"$lt":  dayEnd,
			},
		}
		if hotelId != nil {
			dayMatch["hotelId"] = *hotelId
		}

		dayCount, err := s.repo.Count(ctx, []bson.M{{"$match": dayMatch}})
		if err == nil {
			stats.WeeklyTrend[6-i] = dayCount
		}
	}

	return stats, nil
}

func (s *inquirySvcs) GetUnreadCount(ctx context.Context, hotelId *primitive.ObjectID) (int64, error) {
	match := bson.M{
		"trash":  bson.M{"$ne": true},
		"isRead": false,
	}
	if hotelId != nil {
		match["hotelId"] = *hotelId
	}

	return s.repo.Count(ctx, []bson.M{{"$match": match}})
}

// func (s *inquirySvcs) BulkUpdateStatus(ctx context.Context, inquiryIds []string, status models.InquiryStatus) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	objectIds := make([]primitive.ObjectID, len(inquiryIds))
// 	for i, id := range inquiryIds {
// 		objId, err := primitive.ObjectIDFromHex(id)
// 		if err != nil {
// 			return helpers.InvalidObjectId()
// 		}
// 		objectIds[i] = objId
// 	}

// 	filter := bson.M{"_id": bson.M{"$in": objectIds}}
// 	update := bson.M{"$set": bson.M{
// 		"status":    status,
// 		"updatedAt": time.Now(),
// 		"updatedBy": cfg.User.Id,
// 	}}

// 	_, err = s.repo.Patch(ctx, filter, update)
// 	return err
// }

// func (s *inquirySvcs) BulkMarkAsRead(ctx context.Context, inquiryIds []string, isRead bool) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	objectIds := make([]primitive.ObjectID, len(inquiryIds))
// 	for i, id := range inquiryIds {
// 		objId, err := primitive.ObjectIDFromHex(id)
// 		if err != nil {
// 			return helpers.InvalidObjectId()
// 		}
// 		objectIds[i] = objId
// 	}

// 	updateFields := bson.M{
// 		"isRead":    isRead,
// 		"updatedAt": time.Now(),
// 		"updatedBy": cfg.User.Id,
// 	}

// 	if isRead {
// 		now := time.Now()
// 		updateFields["readAt"] = &now
// 		updateFields["readBy"] = &cfg.User.Id
// 	} else {
// 		updateFields["readAt"] = nil
// 		updateFields["readBy"] = nil
// 	}

// 	filter := bson.M{"_id": bson.M{"$in": objectIds}}
// 	update := bson.M{"$set": updateFields}

// 	_, err = s.repo.Patch(ctx, filter, update)
// 	return err
// }
