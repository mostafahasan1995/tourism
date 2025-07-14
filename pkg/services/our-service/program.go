package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/services/picklist"
	pModels "larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	interactionsModels "larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var updatedByUserLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "users",
		"localField":   "updatedBy",
		"foreignField": "_id",
		"as":           "updatedByUser",
	}},
	{"$unwind": bson.M{
		"path":                       "$updatedByUser",
		"preserveNullAndEmptyArrays": true,
	}},
	{
		"$set": bson.M{
			"updatedByName": bson.M{
				"$concat": []interface{}{
					"$updatedByUser.firstName",
					" ",
					"$updatedByUser.lastName",
				},
			},
		},
	},
	{
		"$project": bson.M{
			"updatedByUser": 0,
		},
	},
}

var packageLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "tourismPackages",
		"localField":   "package",
		"foreignField": "_id",
		"as":           "packageObj",
	}},
	{"$unwind": bson.M{
		"path":                       "$packageObj",
		"preserveNullAndEmptyArrays": true,
	}},
	{
		"$set": bson.M{
			"packageName": "$packageObj.name",
		},
	},
	{
		"$project": bson.M{
			"packageObj": 0,
		},
	},
}

var customerLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "tourismCustomers",
		"localField":   "customerId",
		"foreignField": "_id",
		"as":           "customer",
	}},
	{"$unwind": bson.M{
		"path":                       "$customer",
		"preserveNullAndEmptyArrays": true,
	}},
	{
		"$set": bson.M{
			"customerName": "$customer.name",
		},
	},
	{
		"$project": bson.M{
			"customer": 0,
		},
	},
}

var durationLookup = bson.M{
	"$addFields": bson.M{
		"duration": bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$and": []interface{}{
						bson.M{"$ne": []interface{}{"$startDate", nil}},
						bson.M{"$ne": []interface{}{"$endDate", nil}},
					},
				},
				"then": bson.M{
					"$dateDiff": bson.M{
						"startDate": "$startDate",
						"endDate":   "$endDate",
						"unit":      "day",
					},
				},
				"else": 0,
			},
		},
	},
}

type ProgramSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ProgramRes, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error)
	GetAll(ctx context.Context, query string) ([]models.Program, error)
	GetAuth(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error)
	GetAllAuth(ctx context.Context, query string) ([]models.Program, error)
	Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error)
	Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error)
	UpdateIsFav(ctx context.Context, id string, isFav bool) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ProgramPagination, error)
}

type programsvcs struct {
	repo           repo.ProgramRepo
	travelreqsvcs  TravelRequestSvcs
	withtxn        *db.WithTxn
	activitiesSvcs picklist.ActivitiesSvcs
}

func NewProgramSvcs(i *do.Injector) (ProgramSvcs, error) {
	return &programsvcs{
		repo:           do.MustInvoke[repo.ProgramRepo](i),
		travelreqsvcs:  do.MustInvoke[TravelRequestSvcs](i),
		withtxn:        do.MustInvoke[*db.WithTxn](i),
		activitiesSvcs: do.MustInvoke[picklist.ActivitiesSvcs](i),
	}, nil
}

//

func (p *programsvcs) GetOne(ctx context.Context, id string) (*models.ProgramRes, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Build aggregation pipeline for detailed program information
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":   _id,
				"trash": false,
			},
		},
	}

	// Add lookups for comprehensive data
	pipeline = append(pipeline, customerLookup...)
	pipeline = append(pipeline, packageLookup...)
	pipeline = append(pipeline, updatedByUserLookup...)
	pipeline = append(pipeline, durationLookup)
	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeProgram)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}
	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	if len(result) == 0 {
		return nil, errors.New("program not found")
	}

	// isFav is now populated directly from the pipeline
	//result[0].IsFav = result[0].Program.IsFav

	return &result[0], nil
}

func (p *programsvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error) {
	match := bson.M{"trash": false}

	f, err := helpers.ParseFilters[filter.ProgramFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)

	// Add customer, package, and user lookups first
	pipeline = append(pipeline, customerLookup...)
	pipeline = append(pipeline, packageLookup...)
	pipeline = append(pipeline, updatedByUserLookup...)

	// Calculate duration in days between startDate and endDate
	pipeline = append(pipeline, durationLookup)

	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeProgram)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := p.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	// The isFav field should already be populated correctly from the pipeline
	// No manual assignment needed as it comes directly from the aggregation

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.ProgramPagination{
		Programs:   result,
		Pagination: pagination,
	}, nil
}

func (p *programsvcs) GetAll(ctx context.Context, query string) ([]models.Program, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.ProgramFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	// add favorite pipeline
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeProgram)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}
	var result []models.Program
	err = p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAuth retrieves programs with pagination, requiring authentication
func (p *programsvcs) GetAuth(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as Get but ensure user is authenticated
	return p.Get(ctx, skip, limit, query)
}

// GetAllAuth retrieves all programs without pagination, requiring authentication
func (p *programsvcs) GetAllAuth(ctx context.Context, query string) ([]models.Program, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as GetAll but ensure user is authenticated
	return p.GetAll(ctx, query)
}

// add general or custom program - update related travel request
func (p *programsvcs) Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error) {
	result, err := p.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}
		// Check if this is a general program and has daily itinerary
		if data.ProgramType == "general" && data.GeneralType != nil {
			// Loop through daily itinerary
			for i, day := range data.GeneralType.DailyItinerary {
				// Check if NewActions exists and has elements
				if len(day.NewActions) > 0 {
					for _, actionName := range day.NewActions {
						newActivity := &pModels.ActivitiesDto{
							Name:        actionName,
							Description: nil, // You can customize this
							// Add other fields as needed
						}
						createdActivity, err := p.activitiesSvcs.Add(ctx, newActivity)
						if err != nil {
							return nil, errors.New("failed to create new activity: " + err.Error())
						}
						data.GeneralType.DailyItinerary[i].Actions = append(data.GeneralType.DailyItinerary[i].Actions, createdActivity.Id)
					}
				}
			}
		}

		program := &models.Program{
			Id:         primitive.NewObjectID(),
			ProgramDto: *data,
			CreatedAt:  time.Now(),
			CreatedBy:  cfg.User.Id,
			UpdatedAt:  time.Now(),
			UpdatedBy:  cfg.User.Id,
		}

		if program.TravelReqId != primitive.NilObjectID {
			if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
				return nil, errors.New("error updating travel request, check if it is already assigned to a program")
			}
		}

		if err := p.repo.Add(ctx, program); err != nil {
			return nil, err
		}

		return program, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Program), nil
}

func (p *programsvcs) AssignProgramToTravelRequest(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"_id": program.TravelReqId,
		"$or": []bson.M{
			{"program": primitive.NilObjectID},
			{"program": bson.M{"$exists": false}},
		},
		"trash": false,
	}
	update := bson.M{"$set": bson.M{
		"program":     program.Id,
		"package":     program.Package,
		"status":      enums.TravelReqStatusWaiting,
		"revisionNum": 1,
	}}

	if _, err := p.travelreqsvcs.Patch(ctx, filter, update); err != nil {
		return errors.New("error updating travel request, check if it is already assigned to a program")
	}

	return nil
}

func (p *programsvcs) UnassignProgramFromTravelRequest(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"program": program.Id,
	}
	update := bson.M{"$set": bson.M{
		"program":     primitive.NilObjectID,
		"package":     primitive.NilObjectID,
		"revisionNum": 0,
	}}

	_, err := p.travelreqsvcs.Patch(ctx, filter, update)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	return nil
}

func (p *programsvcs) UpdateTravelRequestRevisionNum(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"_id":     program.TravelReqId,
		"program": program.Id,
	}
	update := bson.M{
		"$inc": bson.M{"revisionNum": 1},
		"$set": bson.M{"status": enums.TravelReqStatusWaiting},
	}

	if _, err := p.travelreqsvcs.Patch(ctx, filter, update); err != nil {
		return errors.New("error updating travel request")
	}

	return nil
}

func (p *programsvcs) Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error) {
	result, err := p.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}

		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		program := &models.Program{
			Id:         _id,
			ProgramDto: *data,
			UpdatedAt:  time.Now(),
			UpdatedBy:  cfg.User.Id,
		}

		if program.TravelReqId == primitive.NilObjectID {
			if err := p.UnassignProgramFromTravelRequest(ctx, program); err != nil {
				return nil, err
			}
		} else {
			currentProgram, err := p.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
			if err != nil {
				return nil, errors.New("error fetching program")
			}

			if currentProgram.TravelReqId == primitive.NilObjectID {
				if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
					return nil, err
				}
			} else {
				if currentProgram.TravelReqId == program.TravelReqId {
					if err := p.UpdateTravelRequestRevisionNum(ctx, program); err != nil {
						return nil, err
					}
				} else {

					if err := p.UnassignProgramFromTravelRequest(ctx, program); err != nil {
						return nil, err
					}

					if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
						return nil, err
					}
				}
			}

		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": program}

		updatedProgram, err := p.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedProgram, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Program), nil
}

func (p *programsvcs) UpdateIsFav(ctx context.Context, id string, isFav bool) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"isFav": isFav}}

	if _, err := p.repo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p *programsvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	if _, err := p.repo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p *programsvcs) Count(ctx context.Context, filter any) (int64, error) {
	return p.repo.Count(ctx, filter)
}

// v2
func (p *programsvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ProgramPagination, error) {

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

	count, err := p.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// Add favorite status using helper function
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		// User is authenticated - add favorite lookup
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeProgram)...)
	} else {
		// User not authenticated - set default favorite status
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.ProgramPagination{
		Programs:   result,
		Pagination: pagination,
	}, nil

}
