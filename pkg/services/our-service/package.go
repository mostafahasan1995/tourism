package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
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
	Get(ctx context.Context, skip, limit int64, query any) (*models.PackageWithPagination, error)
	GetAll(ctx context.Context, query any) ([]models.Package, error)
	Add(ctx context.Context, data *models.PackageDto) (*models.Package, error)
	Update(ctx context.Context, pkgId string, data *models.PackageDto) (*models.Package, error)
	Delete(ctx context.Context, pkgId string) error
	SeedDefaults(ctx context.Context) error
}

type packagesvcs struct {
	repo repo.PackageRepo
}

func NewPackageSvcs(i *do.Injector) (PackageSvcs, error) {
	return &packagesvcs{
		repo: do.MustInvoke[repo.PackageRepo](i),
	}, nil
}

type packageSeed struct {
	IDHex        string
	Name         map[string]string
	Status       enums.PackageStatus
	AllowGeneral bool
	AllowCustom  bool
}

var defaultPackageSeeds = []packageSeed{
	{
		IDHex:        "66f900000000000000000001",
		Name:         map[string]string{"en": "Delegation Travel"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
	{
		IDHex:        "66f900000000000000000002",
		Name:         map[string]string{"en": "Luxury Travel"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
	{
		IDHex:        "66f900000000000000000003",
		Name:         map[string]string{"en": "Family Travel"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
	{
		IDHex:        "66f900000000000000000004",
		Name:         map[string]string{"en": "Honeymoon Package"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
	{
		IDHex:        "66f900000000000000000005",
		Name:         map[string]string{"en": "Business Man Travel"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
	{
		IDHex:        "66f900000000000000000006",
		Name:         map[string]string{"en": "Religious Tourism"},
		Status:       enums.PackageStatusActive,
		AllowGeneral: false,
		AllowCustom:  true,
	},
}

func (p *packagesvcs) ensureDefaultPackages(ctx context.Context) error {
	if len(defaultPackageSeeds) == 0 {
		return nil
	}

	now := time.Now()
	operations := make([]mongo.WriteModel, 0, len(defaultPackageSeeds))

	for _, seed := range defaultPackageSeeds {
		id, err := primitive.ObjectIDFromHex(seed.IDHex)
		if err != nil {
			return err
		}

		update := bson.M{
			"$set": bson.M{
				"name":                seed.Name,
				"status":              seed.Status,
				"isSystemPkg":         true,
				"allowGeneralProgram": seed.AllowGeneral,
				"allowCustomProgram":  seed.AllowCustom,
				"trash":               false,
			},
			"$setOnInsert": bson.M{
				"thumbnail": []any{},
				"createdAt": now,
				"createdBy": primitive.NilObjectID,
				"updatedAt": now,
				"updatedBy": primitive.NilObjectID,
			},
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, op)
	}

	_, err := p.repo.BulkWrite(ctx, operations)
	return err
}

func (p *packagesvcs) SeedDefaults(ctx context.Context) error {
	return p.ensureDefaultPackages(ctx)
}

func (p *packagesvcs) GetAll(ctx context.Context, query any) ([]models.Package, error) {
	match := bson.M{"trash": false}
	f, err := helpers.ParseFilters[filter.PackageFilter](query)
	if err != nil {
		return nil, err
	}

	pipeline := f.BuildPipeline(match)

	var result []models.Package
	errAgg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if errAgg != nil {
		return nil, errAgg
	}

	return result, nil
}

func (p *packagesvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.PackageWithPagination, error) {
	match := bson.M{"trash": false}

	f, err := helpers.ParseFilters[filter.PackageFilter](query)
	if err != nil {
		return nil, err
	}

	pipeline := f.BuildPipeline(match)

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

	existingPkg, err := p.repo.GetByFilter(ctx, bson.M{"_id": _id})
	if err != nil {
		return err
	}

	if existingPkg.IsSystemPkg {
		return errors.New("cannot delete system package")
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
