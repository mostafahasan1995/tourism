package repository

import (
	"context"
	"larsa-tourism-microservices/pkg/services/statistics/models"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
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
	db       *mongo.Client
	collName string
}

func NewManualStatisticsRepository(i *do.Injector) (ManualStatisticsRepository, error) {
	return &manualStatisticsRepository{
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "manualStatistics",
	}, nil
}

func (r *manualStatisticsRepository) GetManualStatistics(ctx context.Context) (*models.ManualStatistics, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	coll := r.db.Database(cfg.Db).Collection(r.collName)
	var stats models.ManualStatistics
	err = coll.FindOne(ctx, bson.M{}).Decode(&stats)
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
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	coll := r.db.Database(cfg.Db).Collection(r.collName)
	stats.UpdatedAt = time.Now()

	filter := bson.M{}
	update := bson.M{
		"$set": stats,
	}

	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *manualStatisticsRepository) CreateManualStatistics(ctx context.Context, stats *models.ManualStatistics) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	coll := r.db.Database(cfg.Db).Collection(r.collName)
	stats.ID = primitive.NewObjectID()
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	_, err = coll.InsertOne(ctx, stats)
	return err
}

func (r *manualStatisticsRepository) ToggleAutoCalculate(ctx context.Context, autoCalculate bool) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	coll := r.db.Database(cfg.Db).Collection(r.collName)
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"autoCalculate": autoCalculate,
			"updatedAt":     time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	return err
}
