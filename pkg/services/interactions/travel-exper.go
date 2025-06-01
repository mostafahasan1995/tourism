package interactions

import (
	"context"
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

type TravelExperSvcs interface {
	GetTravelerStory(ctx context.Context, storyId string) (*models.TravelerStory, error)
	GetAllTravelerStories(ctx context.Context) ([]models.TravelerStory, error)
	GetTravelerStories(ctx context.Context, skip, limit int64) (*models.TravelerStoryWithPagination, error)
	AddTravelerStory(ctx context.Context, data *models.TravelerStoryDto) (*models.TravelerStory, error)
	SetTravelerStoryStatus(ctx context.Context, storyId string, data *models.TravelerStoryStatusDto) (*models.TravelerStory, error)
	SendFeedback(ctx context.Context, storyId string, data *models.TravelerStoryFeedback) (*models.TravelerStory, error)
	UpdateTravelerStory(ctx context.Context, storyId string, data *models.TravelerStoryDto) (*models.TravelerStory, error)
	DeleteTravlerStory(ctx context.Context, storyId string) error
	RestoreTravlerStory(ctx context.Context, storyId string) error
	// client
	GetClientStory(ctx context.Context, storyId string) (*models.ClientStory, error)
	GetClientStories(ctx context.Context, skip, limit int64) (*models.ClientStoryWithPagination, error)
	AddClientStory(ctx context.Context, data *models.ClientStoryDto) (*models.ClientStory, error)
	UpdateClientStory(ctx context.Context, storyId string, data *models.ClientStoryDto) (*models.ClientStory, error)
	DeleteClientStory(ctx context.Context, storyId string) error
}

type travelexpersvcs struct {
	travelerStoryRepo repo.TravelerStoryRepo
	clientStoryRepo   repo.ClientStoryRepo
}

func NewTravelExperSvcs(i *do.Injector) (TravelExperSvcs, error) {
	return &travelexpersvcs{
		travelerStoryRepo: do.MustInvoke[repo.TravelerStoryRepo](i),
		clientStoryRepo:   do.MustInvoke[repo.ClientStoryRepo](i),
	}, nil
}

// traveler
func (t *travelexpersvcs) GetTravelerStory(ctx context.Context, storyId string) (*models.TravelerStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}

	return t.travelerStoryRepo.GetByFilter(ctx, filter)
}

func (t *travelexpersvcs) GetAllTravelerStories(ctx context.Context) ([]models.TravelerStory, error) {

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	var result []models.TravelerStory
	err := t.travelerStoryRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

func (t *travelexpersvcs) GetTravelerStories(ctx context.Context, skip, limit int64) (*models.TravelerStoryWithPagination, error) {
	match := bson.M{"trash": false}

	count, err := t.travelerStoryRepo.Count(ctx, match)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}
	var result []models.TravelerStory
	errAg := t.travelerStoryRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
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

	return &models.TravelerStoryWithPagination{
		TravelerStories: result,
		Pagination:      pagination,
	}, nil
}

func (t *travelexpersvcs) AddTravelerStory(ctx context.Context, data *models.TravelerStoryDto) (*models.TravelerStory, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	story := &models.TravelerStory{
		Id:               primitive.NewObjectID(),
		TravelerStoryDto: *data,
		CreatedAt:        time.Now(),
		CreatedBy:        cfg.User.Id,
	}

	story.Status = "pending"

	if err := t.travelerStoryRepo.Add(ctx, story); err != nil {
		return nil, err
	}

	return story, nil

}

func (t *travelexpersvcs) SetTravelerStoryStatus(ctx context.Context, storyId string, data *models.TravelerStoryStatusDto) (*models.TravelerStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"status":       data.Status,
			"rejectReason": data.Reason,
			"updatedAt":    time.Now(),
			"updatedBy":    cfg.User.Id,
		},
	}

	updatedStory, err := t.travelerStoryRepo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedStory, nil

}

func (t *travelexpersvcs) SendFeedback(ctx context.Context, storyId string, data *models.TravelerStoryFeedback) (*models.TravelerStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": bson.M{
			"feedback":  data.Feedback,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	updatedStory, err := t.travelerStoryRepo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	//todo : send feedback to the traveler by email

	return updatedStory, nil
}

func (t *travelexpersvcs) UpdateTravelerStory(ctx context.Context, storyId string, data *models.TravelerStoryDto) (*models.TravelerStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	story := &models.TravelerStory{
		Id:               _id,
		TravelerStoryDto: *data,
		UpdatedAt:        time.Now(),
		UpdatedBy:        cfg.User.Id,
	}

	story.Status = "edited"

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$set": story,
	}

	updatedStory, err := t.travelerStoryRepo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedStory, nil
}

func (t *travelexpersvcs) DeleteTravlerStory(ctx context.Context, storyId string) error {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"trash": true}}

	if _, err := t.travelerStoryRepo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (t *travelexpersvcs) RestoreTravlerStory(ctx context.Context, storyId string) error {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"trash": false}}

	if _, err := t.travelerStoryRepo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// client

func (t *travelexpersvcs) GetClientStory(ctx context.Context, storyId string) (*models.ClientStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id, "trash": false}

	return t.clientStoryRepo.GetByFilter(ctx, filter)
}

func (t *travelexpersvcs) GetClientStories(ctx context.Context, skip, limit int64) (*models.ClientStoryWithPagination, error) {
	match := bson.M{"trash": false}

	count, err := t.clientStoryRepo.Count(ctx, match)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}
	var result []models.ClientStory
	errAg := t.clientStoryRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
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

	return &models.ClientStoryWithPagination{
		ClientStories: result,
		Pagination:    pagination,
	}, nil
}

func (t *travelexpersvcs) AddClientStory(ctx context.Context, data *models.ClientStoryDto) (*models.ClientStory, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	story := &models.ClientStory{
		Id:             primitive.NewObjectID(),
		ClientStoryDto: *data,
		CreatedAt:      time.Now(),
		CreatedBy:      cfg.User.Id,
	}

	if err := t.clientStoryRepo.Add(ctx, story); err != nil {
		return nil, err
	}

	return story, nil

}

func (t *travelexpersvcs) UpdateClientStory(ctx context.Context, storyId string, data *models.ClientStoryDto) (*models.ClientStory, error) {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return nil, err
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	story := &models.ClientStory{
		Id:             _id,
		ClientStoryDto: *data,
		UpdatedAt:      time.Now(),
		UpdatedBy:      cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": story}

	updatedStory, err := t.clientStoryRepo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedStory, nil
}

func (t *travelexpersvcs) DeleteClientStory(ctx context.Context, storyId string) error {
	_id, err := primitive.ObjectIDFromHex(storyId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"trash": true}}

	if _, err := t.clientStoryRepo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
