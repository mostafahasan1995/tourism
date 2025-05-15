package ourService

import (
	"context"
	"errors"
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

type PackageSvcs interface {
	Get(ctx context.Context, skip, limit int64, query string) (*models.PackageWithPagination, error)
	GetAll(ctx context.Context) ([]models.Package, error)
	Add(ctx context.Context, data *models.PackageDto) (*models.Package, error)
	Update(ctx context.Context, pkgId string, data *models.PackageDto) (*models.Package, error)
	Delete(ctx context.Context, pkgId string) error
}

type packagesvcs struct {
	repo repo.PackageRepo
}

func NewPackageSvcs(i *do.Injector) (PackageSvcs, error) {
	return &packagesvcs{
		repo: do.MustInvoke[repo.PackageRepo](i),
	}, nil
}

func (p *packagesvcs) GetAll(ctx context.Context) ([]models.Package, error) {

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	var result []models.Package
	err := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *packagesvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.PackageWithPagination, error) {

	match := bson.M{"trash": false}

	filters, err := filter.NewPackageFilter(query)
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

	var result []models.Package
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

	return &models.PackageWithPagination{
		Packages:   result,
		Pagination: pagination,
	}, nil
}

func (p *packagesvcs) Add(ctx context.Context, data *models.PackageDto) (*models.Package, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	pkg := &models.Package{
		Id:         primitive.NewObjectID(),
		PackageDto: *data,
		CreatedAt:  time.Now(),
		CreatedBy:  cfg.User.Id,
	}

	if err := p.repo.Add(ctx, pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}

func (p *packagesvcs) Update(ctx context.Context, pkgId string, data *models.PackageDto) (*models.Package, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(pkgId)
	if err != nil {
		return nil, err
	}

	pkg := &models.Package{
		Id:         _id,
		PackageDto: *data,
		UpdatedAt:  time.Now(),
		UpdatedBy:  cfg.User.Id,
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": pkg}

	updatedPackage, err := p.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedPackage, nil
}

func (p *packagesvcs) Delete(ctx context.Context, pkgId string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(pkgId)
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
