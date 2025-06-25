package repository

import (
	"context"
	"larsa-tourism-microservices/pkg/services/statistics/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ManualStatisticsRepository interface {
	GetManualStatistics(ctx context.Context) (*models.ManualStatistics, error)
	UpdateManualStatistics(ctx context.Context, stats *models.ManualStatistics) error
	CreateManualStatistics(ctx context.Context, stats *models.ManualStatistics) error
	ToggleAutoCalculate(ctx context.Context, autoCalculate bool) error
}

type manualStatisticsRepository struct {
	collection *mongo.Collection
}

func NewManualStatisticsRepository(db *mongo.Database) ManualStatisticsRepository {
	return &manualStatisticsRepository{
		collection: db.Collection("manual_statistics"),
	}
}

func (r *manualStatisticsRepository) GetManualStatistics(ctx context.Context) (*models.ManualStatistics, error) {
	var stats models.ManualStatistics
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default statistics with auto-calculate enabled
			return &models.ManualStatistics{
				AutoCalculate: true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}, nil
		}
		return nil, err
	}
	return &stats, nil
}

func (r *manualStatisticsRepository) UpdateManualStatistics(ctx context.Context, stats *models.ManualStatistics) error {
	stats.UpdatedAt = time.Now()

	filter := bson.M{}
	update := bson.M{
		"$set": stats,
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *manualStatisticsRepository) CreateManualStatistics(ctx context.Context, stats *models.ManualStatistics) error {
	stats.ID = primitive.NewObjectID()
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, stats)
	return err
}

func (r *manualStatisticsRepository) ToggleAutoCalculate(ctx context.Context, autoCalculate bool) error {
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"autoCalculate": autoCalculate,
			"updatedAt":     time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
