package repo

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/services/db/models"
	"larsa-tourism-microservices/pkg/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SortingRepo interface {
	GetAndUpdateSourceSeq(ctx context.Context, name string) (int64, error)
	Reorder(ctx context.Context, id, collName string, newSeq int64) error
}

type sortingrepo struct {
	db       *mongo.Client
	collName string
}

func NewSortingRepo(i *do.Injector) (SortingRepo, error) {
	return &sortingrepo{
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismSorting",
	}, nil
}

func (s *sortingrepo) GetAndUpdateSourceSeq(ctx context.Context, name string) (int64, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return 0, err
	}

	coll := s.db.Database(cfg.Db).Collection(s.collName)

	filter := bson.M{"name": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	upsert := true
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	)

	var data models.Counter
	if err := result.Decode(&data); err != nil {
		return 0, err
	}

	return data.Seq, nil
}

func (s *sortingrepo) Reorder(ctx context.Context, id, collName string, newSeq int64) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := s.db.Database(cfg.Db).Collection(collName)

	var item struct {
		Seq int64 `bson:"seq"`
	}

	if err := coll.FindOne(ctx, bson.M{
		"_id":   _id,
		"trash": bson.M{"$ne": true},
	}).Decode(&item); err != nil {
		return errors.New("error get related item")
	}

	var filter, update bson.M

	if newSeq > item.Seq {

		filter = bson.M{"$and": bson.A{
			bson.M{"seq": bson.M{"$gt": item.Seq}},
			bson.M{"seq": bson.M{"$lte": newSeq}},
		}}

		update = bson.M{"$inc": bson.M{"seq": -1}}

	} else if newSeq < item.Seq {
		filter = bson.M{"$and": bson.A{
			bson.M{"seq": bson.M{"$gte": newSeq}},
			bson.M{"seq": bson.M{"$lt": item.Seq}},
		}}

		update = bson.M{"$inc": bson.M{"seq": 1}}

	} else {
		return errors.New("error sorting data")
	}

	_, errUp := coll.UpdateMany(ctx, filter, update)
	if errUp != nil {
		return errUp
	}

	_, errUpd := coll.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": bson.M{"seq": newSeq}})
	if errUpd != nil {
		return errUpd
	}

	return nil

}
