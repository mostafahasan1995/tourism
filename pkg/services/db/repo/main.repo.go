package repo

import (
	"context"
	"larsa-tourism-microservices/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MainRepo[T any] interface {
	GetByFilter(ctx context.Context, filter bson.M) (*T, error)
	Aggregate(ctx context.Context, pipeline any, callback func(cur *mongo.Cursor) error) error
	Add(ctx context.Context, data *T) error
	AddMany(ctx context.Context, data []any) error
	Patch(ctx context.Context, filter, update bson.M, ops ...*options.FindOneAndUpdateOptions) (*T, error)
	BulkWrite(ctx context.Context, writeOps []mongo.WriteModel) (*mongo.BulkWriteResult, error)
}

type MainRepoImpl[T any] struct {
	Db       *mongo.Client
	CollName string
}

func (m *MainRepoImpl[T]) Add(ctx context.Context, data *T) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	if _, err := coll.InsertOne(ctx, data); err != nil {
		return err
	}
	return nil
}

func (m *MainRepoImpl[T]) AddMany(ctx context.Context, data []any) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	if _, err := coll.InsertMany(ctx, data); err != nil {
		return err
	}
	return nil
}

func (m *MainRepoImpl[T]) GetByFilter(ctx context.Context, filter bson.M) (*T, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	var result T
	if err := coll.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MainRepoImpl[T]) Patch(ctx context.Context, filter, update bson.M, ops ...*options.FindOneAndUpdateOptions) (*T, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	ops = append(ops, opts)

	res := coll.FindOneAndUpdate(ctx, filter, update, ops...)

	var result T
	if err := res.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (m *MainRepoImpl[T]) Aggregate(ctx context.Context, pipeline any, callback func(cur *mongo.Cursor) error) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	cur, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}

	if err := callback(cur); err != nil {
		return err
	}

	return nil
}

func (m *MainRepoImpl[T]) BulkWrite(ctx context.Context, writeOps []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	coll := m.Db.Database(cfg.Db).Collection(m.CollName)

	opts := options.BulkWrite().SetOrdered(true)
	result, errIns := coll.BulkWrite(ctx, writeOps, opts)

	if errIns != nil {
		return nil, errIns
	}

	return result, nil
}
