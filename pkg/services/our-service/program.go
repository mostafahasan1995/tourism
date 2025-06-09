package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/db"
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

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProgramSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Program, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error)
	GetAll(ctx context.Context, query string) ([]models.Program, error)
	Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error)
	Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error)
	Delete(ctx context.Context, id string) error
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

func (p *programsvcs) GetOne(ctx context.Context, id string) (*models.Program, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return p.repo.GetByFilter(ctx, bson.M{"_id": _id})
}

func (p *programsvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error) {
	match := bson.M{"trash": false}

	filters, err := filter.NewProgramFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := p.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Program
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

func (p *programsvcs) GetAll(ctx context.Context, query string) ([]models.Program, error) {
	match := bson.M{"trash": false}

	filters, err := filter.NewProgramFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	var result []models.Program
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
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
							Description: "", // You can customize this
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
