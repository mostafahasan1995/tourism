package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
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
	Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error)
	Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error)
	Delete(ctx context.Context, id string) error
}

type programsvcs struct {
	repo          repo.ProgramRepo
	travelreqsvcs TravelRequestSvcs
	withtxn       *db.WithTxn
}

func NewProgramSvcs(i *do.Injector) (ProgramSvcs, error) {
	return &programsvcs{
		repo:          do.MustInvoke[repo.ProgramRepo](i),
		travelreqsvcs: do.MustInvoke[TravelRequestSvcs](i),
		withtxn:       do.MustInvoke[*db.WithTxn](i),
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

// add general or custom program - update related travel request
func (p *programsvcs) Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error) {
	result, err := p.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}

		program := &models.Program{
			Id:         primitive.NewObjectID(),
			ProgramDto: *data,
			CreatedAt:  time.Now(),
			CreatedBy:  cfg.User.Id,
		}

		if program.TravelReqId != primitive.NilObjectID {

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

func (p *programsvcs) Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error) {
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

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": program}

	updatedProgram, err := p.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedProgram, nil
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
