package picklist

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	interactionsModels "larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/picklist/filter"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/services/picklist/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationSvcs interface {
	GetOne(ctx context.Context, id string) (*models.DestinationRes, error)
	Get(ctx context.Context, skip, limit int64, query any) (*models.DestinationPaginationRes, error)
	GetAll(ctx context.Context, query any) ([]models.Destination, error)
	Add(ctx context.Context, data *models.DestinationDto) (*models.Destination, error)
	AddManyNameOnly(ctx context.Context, data []string) (countriesToSave []string, err error)
	Update(ctx context.Context, id string, data *models.DestinationDto) (*models.Destination, error)
	UpdateIsFav(ctx context.Context, id string, isFav bool) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	GetDestinationByCountry(ctx context.Context, data *models.DestinationCountry) (*models.Destination, error)
	//v2
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.DestinationRes, error)
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.DestinationPaginationRes, error)
}

type destinationSvcs struct {
	repo repo.DestinationRepo
}

func NewDestinationSvcs(i *do.Injector) (DestinationSvcs, error) {
	return &destinationSvcs{
		repo: do.MustInvoke[repo.DestinationRepo](i),
	}, nil
}

func (d *destinationSvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.DestinationPaginationRes, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.DestinationFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := d.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeDestination)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.DestinationRes
	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.DestinationPaginationRes{
		Destinations: result,
		Pagination:   pagination,
	}, nil
}

func (d *destinationSvcs) GetOne(ctx context.Context, id string) (*models.DestinationRes, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": _id, "trash": false}}}

	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeDestination)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.DestinationRes
	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return &result[0], nil
}

func (d *destinationSvcs) GetAll(ctx context.Context, query any) ([]models.Destination, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.DestinationFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	var result []models.Destination
	errAg := d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (d *destinationSvcs) Add(ctx context.Context, data *models.DestinationDto) (*models.Destination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	destination := &models.Destination{
		Id:             primitive.NewObjectID(),
		DestinationDto: *data,
		CreatedAt:      time.Now(),
		CreatedBy:      cfg.User.Id,
	}

	if err := d.repo.Add(ctx, destination); err != nil {
		return nil, err
	}

	return destination, nil
}

func (d *destinationSvcs) AddManyNameOnly(ctx context.Context, countries []string) (countriesToSave []string, err error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	if len(countries) == 0 {
		return nil, errors.New("no data provided")
	}
	// The destinations to be added to the database
	createdDests := make([]any, 0, len(countries))
	// the countries provided updated with the case found in the db
	countriesToSave = make([]string, 0, len(countries))

	for _, item := range countries {
		// Chcecking if the destination exists regardless of the case
		filter := bson.M{"name": bson.M{"$regex": "^" + item + "$", "$options": "i"}, "trash": false}
		dest, err := d.repo.GetByFilter(ctx, filter)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return countriesToSave, err
		}
		if dest != nil {
			countriesToSave = append(countriesToSave, dest.Name)
			continue
		}
		destination := &models.Destination{
			Id: primitive.NewObjectID(),
			DestinationDto: models.DestinationDto{
				Name: item,
			},
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
		}
		countriesToSave = append(countriesToSave, item)
		createdDests = append(createdDests, destination)
	}

	if len(createdDests) == 0 {
		return countriesToSave, nil
	}

	err = d.repo.AddMany(ctx, createdDests)
	if err != nil {
		return countriesToSave, err
	}
	return countriesToSave, nil
}

func (d *destinationSvcs) Update(ctx context.Context, id string, data *models.DestinationDto) (*models.Destination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	destination := &models.Destination{
		Id:             _id,
		DestinationDto: *data,
		UpdatedAt:      time.Now(),
		UpdatedBy:      cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": destination}

	updatedDestination, err := d.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedDestination, nil
}

func (d *destinationSvcs) UpdateIsFav(ctx context.Context, id string, isFav bool) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"isFav": isFav, "updatedAt": time.Now(), "updatedBy": cfg.User.Id}}

	_, err = d.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *destinationSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = d.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *destinationSvcs) Count(ctx context.Context, filter any) (int64, error) {
	return d.repo.Count(ctx, filter)
}

func (d *destinationSvcs) GetDestinationByCountry(ctx context.Context, data *models.DestinationCountry) (*models.Destination, error) {
	dest, err := d.repo.GetByFilter(ctx, bson.M{"name": data.Country, "trash": false})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//add new destination
			dest = &models.Destination{
				Id: primitive.NewObjectID(),
				DestinationDto: models.DestinationDto{
					Name: data.Country,
				},
				CreatedAt: time.Now(),
			}

			if err := d.repo.Add(ctx, dest); err != nil {
				return nil, errors.New("failed to add new destination")
			}
			return dest, nil

		} else {
			return nil, err
		}
	}

	return dest, nil
}

// v2
func (d *destinationSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.DestinationRes, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	fmt.Println("cfg", cfg.User)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeDestination)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.DestinationRes
	err = d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *destinationSvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.DestinationPaginationRes, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := d.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeDestination)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.DestinationRes
	err = d.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	// Convert to DestinationRes with isFav populated
	// destinationsRes, err := d.convertToDestinationRes(ctx, result)
	// if err != nil {
	// 	return nil, err
	// }

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.DestinationPaginationRes{
		Destinations: result,
		Pagination:   pagination,
	}, nil
}

// Helper method to convert Destination to DestinationRes with isFav populated
// func (d *destinationSvcs) convertToDestinationRes(ctx context.Context, destinations []models.Destination) ([]models.DestinationRes, error) {
// 	if len(destinations) == 0 {
// 		return []models.DestinationRes{}, nil
// 	}

// 	// Convert to DestinationRes - isFav is now stored directly in the entity
// 	result := make([]models.DestinationRes, len(destinations))
// 	for i, dest := range destinations {
// 		result[i] = models.DestinationRes{
// 			Destination: dest,
// 			IsFav:       dest.IsFav, // Use the isFav field directly from the entity
// 		}
// 	}

// 	return result, nil
// }
