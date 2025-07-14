package exhibition_management

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/exhibition-management/models"
	"larsa-tourism-microservices/pkg/services/exhibition-management/repo"
	interactionsModels "larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExhibitionSvcs defines the interface for exhibition services
type ExhibitionSvcs interface {
	Get(ctx context.Context, skip, limit int64) (*models.ExhibitionWithPagination, error)
	GetAll(ctx context.Context) ([]models.Exhibition, error)
	GetAuth(ctx context.Context, skip, limit int64) (*models.ExhibitionWithPagination, error)
	GetAllAuth(ctx context.Context) ([]models.Exhibition, error)
	GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.ExhibitionWithPagination, error)
	GetById(ctx context.Context, id string) (*models.Exhibition, error)
	Save(ctx context.Context, data *models.ExhibitionDto) (*models.Exhibition, error)
	Update(ctx context.Context, id string, data *models.ExhibitionDto) (*models.Exhibition, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Exhibition, error)
	Toggle(ctx context.Context, id string, isActive bool) (*models.Exhibition, error)
	Delete(ctx context.Context, id string) error
	GetRelatedExhibitions(ctx context.Context, id string) ([]models.Exhibition, error)

	// Ad management methods
	GetExhibitionAds(ctx context.Context, exhibitionId string) ([]models.Ad, error)
	UpdateAd(ctx context.Context, exhibitionId, adId string, data *models.AdDto) (*models.Ad, error)
	DeleteAd(ctx context.Context, exhibitionId, adId string) error
	ToggleAd(ctx context.Context, exhibitionId, adId string, isActive bool) (*models.Ad, error)

	// Validation methods for other services
	ValidateExhibitionExists(ctx context.Context, exhibitionId string) error

	// Debug methods - remove in production
	DebugCount(ctx context.Context) (map[string]interface{}, error)
	DebugRaw(ctx context.Context) ([]models.Exhibition, error)
}

// exhibitionSvcs implements the ExhibitionSvcs interface
type exhibitionSvcs struct {
	repo repo.ExhibitionRepo
}

// NewExhibitionSvcs creates a new instance of ExhibitionSvcs
func NewExhibitionSvcs(i *do.Injector) (ExhibitionSvcs, error) {
	return &exhibitionSvcs{
		repo: do.MustInvoke[repo.ExhibitionRepo](i),
	}, nil
}

// Get retrieves exhibitions with pagination (no filtering)
func (s *exhibitionSvcs) Get(ctx context.Context, skip, limit int64) (*models.ExhibitionWithPagination, error) {
	match := bson.M{"trash": false}

	pipeline := []bson.M{{"$match": match}}

	// Add fave lookup - check if user is authenticated
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		// User is authenticated, add lookup to check favorites
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"exhibitionId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$exhibitionId"}},
									{"$eq": []interface{}{"$type", "exhibition"}},
									{"$eq": []interface{}{"$userId", cfg.User.Id}},
									{"$eq": []interface{}{"$isFav", true}},
								},
							},
							"trash": bson.M{"$ne": true},
						},
					},
				},
				"as": "faveRecord",
			},
		})
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$faveRecord"},
						0,
					},
				},
			},
		})
		pipeline = append(pipeline, bson.M{
			"$project": bson.M{
				"faveRecord": 0,
			},
		})
	} else {
		// User not authenticated, set isFav to false
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": false,
			},
		})
	}

	count, err := s.repo.Count(ctx, match)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Exhibition
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.ExhibitionWithPagination{
		Exhibitions: result,
		Pagination:  pagination,
	}, nil
}

// GetAll retrieves all exhibitions without pagination
func (s *exhibitionSvcs) GetAll(ctx context.Context) ([]models.Exhibition, error) {
	match := bson.M{"trash": false}

	pipeline := []bson.M{{"$match": match}}

	// Add fave lookup - check if user is authenticated
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		// User is authenticated, add lookup to check favorites
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"exhibitionId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$exhibitionId"}},
									{"$eq": []interface{}{"$type", "exhibition"}},
									{"$eq": []interface{}{"$userId", cfg.User.Id}},
									{"$eq": []interface{}{"$isFav", true}},
								},
							},
							"trash": bson.M{"$ne": true},
						},
					},
				},
				"as": "faveRecord",
			},
		})
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$faveRecord"},
						0,
					},
				},
			},
		})
		pipeline = append(pipeline, bson.M{
			"$project": bson.M{
				"faveRecord": 0,
			},
		})
	} else {
		// User not authenticated, set isFav to false
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": false,
			},
		})
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})

	var result []models.Exhibition
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

// GetAuth retrieves exhibitions with pagination, requiring authentication
func (s *exhibitionSvcs) GetAuth(ctx context.Context, skip, limit int64) (*models.ExhibitionWithPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as Get but ensure user is authenticated
	return s.Get(ctx, skip, limit)
}

// GetAllAuth retrieves all exhibitions without pagination, requiring authentication
func (s *exhibitionSvcs) GetAllAuth(ctx context.Context) ([]models.Exhibition, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as GetAll but ensure user is authenticated
	return s.GetAll(ctx)
}

// GetV2 retrieves exhibitions with pagination and filters
func (s *exhibitionSvcs) GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.ExhibitionWithPagination, error) {
	// Start with base match condition
	match := bson.M{}

	// Convert filter conditions to MongoDB filter
	if filters != nil && len(filters.Columns) > 0 {
		filterBson, err := filters.ConvertToMongo()
		if err != nil {
			return nil, helpers.BadRequest("Invalid filter conditions: " + err.Error())
		}

		// Check if user is explicitly filtering on trash field
		hasTrashFilter := false
		for _, column := range filters.Columns {
			if column.Name == "trash" {
				hasTrashFilter = true
				break
			}
		}

		// If user filters have content
		if len(filterBson) > 0 {
			if hasTrashFilter {
				// User is explicitly filtering trash, so don't add our base condition
				match = filterBson
			} else {
				// User is not filtering trash, so add our base condition to exclude trash
				match = bson.M{
					"$and": []bson.M{
						{"trash": false},
						filterBson,
					},
				}
			}
		} else {
			// No valid filters, use base condition
			match = bson.M{"trash": false}
		}
	} else {
		// No filters provided, use base condition
		match = bson.M{"trash": false}
	}

	pipeline := []bson.M{{"$match": match}}

	// Add fave lookup - check if user is authenticated
	cfg, cfgErr := util.GetReqAppCfg(ctx)
	if cfgErr == nil && cfg.User != nil {
		// User is authenticated, add lookup to check favorites
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"exhibitionId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$exhibitionId"}},
									{"$eq": []interface{}{"$type", "exhibition"}},
									{"$eq": []interface{}{"$userId", cfg.User.Id}},
									{"$eq": []interface{}{"$isFav", true}},
								},
							},
							"trash": bson.M{"$ne": true},
						},
					},
				},
				"as": "faveRecord",
			},
		})
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$faveRecord"},
						0,
					},
				},
			},
		})
		pipeline = append(pipeline, bson.M{
			"$project": bson.M{
				"faveRecord": 0,
			},
		})
	} else {
		// User not authenticated, set isFav to false
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"isFav": false,
			},
		})
	}

	count, err := s.repo.Count(ctx, match)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Exhibition
	errAg := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.ExhibitionWithPagination{
		Exhibitions: result,
		Pagination:  pagination,
	}, nil
}

// GetById retrieves an exhibition by ID with full details
func (s *exhibitionSvcs) GetById(ctx context.Context, id string) (*models.Exhibition, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": objID, "trash": false}}}

	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeExhibition)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.Exhibition
	err = s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, helpers.NotFoundError("Exhibition not found")
	}

	return &result[0], nil
}

// Save creates a new exhibition with ads
func (s *exhibitionSvcs) Save(ctx context.Context, data *models.ExhibitionDto) (*models.Exhibition, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	// Set auto status if not provided
	data.SetAutoStatus()

	// Validate related exhibitions exist
	if len(data.RelatedExhibitions) > 0 {
		if err := s.validateRelatedExhibitions(ctx, data.RelatedExhibitions); err != nil {
			return nil, err
		}
	}

	// Create ads from the dto
	var ads []models.Ad
	for _, adDto := range data.Ads {
		// Calculate price if not set
		if adDto.Price == 0 {
			adDto.Price = models.CalculateAdPrice(adDto.PackageType, adDto.Duration)
		}

		// Set activation and expiration dates if not provided
		if adDto.ActivationDate.IsZero() {
			adDto.ActivationDate = time.Now()
		}
		if adDto.ExpirationDate.IsZero() {
			adDto.ExpirationDate = adDto.ActivationDate.AddDate(0, 0, adDto.Duration)
		}

		ad := models.Ad{
			Id:             primitive.NewObjectID(),
			Title:          adDto.Title,
			Description:    adDto.Description,
			Tags:           adDto.Tags,
			Images:         adDto.Images,
			WebsiteUrl:     adDto.WebsiteUrl,
			PackageType:    adDto.PackageType,
			Duration:       adDto.Duration,
			PaymentMethod:  adDto.PaymentMethod,
			PaymentStatus:  adDto.PaymentStatus,
			ActivationDate: adDto.ActivationDate,
			ExpirationDate: adDto.ExpirationDate,
			Price:          adDto.Price,
			IsActive:       adDto.IsActive,
			CreatedAt:      time.Now(),
			CreatedBy:      cfg.User.Id,
			UpdatedAt:      time.Now(),
			UpdatedBy:      cfg.User.Id,
		}
		ads = append(ads, ad)
	}

	// Remove ads from ExhibitionDto before creating the exhibition
	exhibitionDto := *data
	exhibitionDto.Ads = nil

	exhibition := &models.Exhibition{
		Id:            primitive.NewObjectID(),
		ExhibitionDto: exhibitionDto,
		Ads:           ads,
		Trash:         false,
		CreatedAt:     time.Now(),
		CreatedBy:     cfg.User.Id,
		UpdatedAt:     time.Now(),
		UpdatedBy:     cfg.User.Id,
	}

	if err := s.repo.Add(ctx, exhibition); err != nil {
		return nil, err
	}

	return s.GetById(ctx, exhibition.Id.Hex())
}

// Update updates an existing exhibition
func (s *exhibitionSvcs) Update(ctx context.Context, id string, data *models.ExhibitionDto) (*models.Exhibition, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.BadRequest("Invalid exhibition ID")
	}

	// Set auto status if not provided
	data.SetAutoStatus()

	// Check if record exists
	existing, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate related exhibitions exist
	if len(data.RelatedExhibitions) > 0 {
		if err := s.validateRelatedExhibitions(ctx, data.RelatedExhibitions); err != nil {
			return nil, err
		}
	}

	// Create ads from the dto if provided, otherwise preserve existing
	var ads []models.Ad
	if len(data.Ads) > 0 {
		for _, adDto := range data.Ads {
			// Calculate price if not set
			if adDto.Price == 0 {
				adDto.Price = models.CalculateAdPrice(adDto.PackageType, adDto.Duration)
			}

			// Set activation and expiration dates if not provided
			if adDto.ActivationDate.IsZero() {
				adDto.ActivationDate = time.Now()
			}
			if adDto.ExpirationDate.IsZero() {
				adDto.ExpirationDate = adDto.ActivationDate.AddDate(0, 0, adDto.Duration)
			}

			ad := models.Ad{
				Id:             primitive.NewObjectID(),
				Title:          adDto.Title,
				Description:    adDto.Description,
				Tags:           adDto.Tags,
				Images:         adDto.Images,
				WebsiteUrl:     adDto.WebsiteUrl,
				PackageType:    adDto.PackageType,
				Duration:       adDto.Duration,
				PaymentMethod:  adDto.PaymentMethod,
				PaymentStatus:  adDto.PaymentStatus,
				ActivationDate: adDto.ActivationDate,
				ExpirationDate: adDto.ExpirationDate,
				Price:          adDto.Price,
				IsActive:       adDto.IsActive,
				CreatedAt:      time.Now(),
				CreatedBy:      cfg.User.Id,
				UpdatedAt:      time.Now(),
				UpdatedBy:      cfg.User.Id,
			}
			ads = append(ads, ad)
		}
	} else {
		ads = existing.Ads // Preserve existing ads if none provided
	}

	// Remove ads from ExhibitionDto before updating
	exhibitionDto := *data
	exhibitionDto.Ads = nil

	// Update the record
	exhibition := &models.Exhibition{
		Id:            objID,
		ExhibitionDto: exhibitionDto,
		Ads:           ads,
		Trash:         false,
		CreatedAt:     existing.CreatedAt,
		CreatedBy:     existing.CreatedBy,
		UpdatedAt:     time.Now(),
		UpdatedBy:     cfg.User.Id,
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": exhibition}

	if _, err := s.repo.Patch(ctx, filter, update); err != nil {
		return nil, err
	}

	return s.GetById(ctx, id)
}

// Patch updates an existing exhibition with a map of fields to update
func (s *exhibitionSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.Exhibition, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.BadRequest("Invalid exhibition ID")
	}

	// Add audit fields to the update
	updateDoc := bson.M{}
	for key, value := range updates {
		updateDoc[key] = value
	}
	updateDoc["updatedAt"] = time.Now()
	updateDoc["updatedBy"] = cfg.User.Id

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateDoc}

	updatedExhibition, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updatedExhibition == nil {
		return nil, helpers.NotFoundError("Exhibition not found")
	}

	return updatedExhibition, nil
}

// Toggle enables or disables an exhibition
func (s *exhibitionSvcs) Toggle(ctx context.Context, id string, isActive bool) (*models.Exhibition, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.BadRequest("Invalid exhibition ID")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"isActive":  isActive,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	if _, err := s.repo.Patch(ctx, filter, update); err != nil {
		return nil, err
	}

	return s.GetById(ctx, id)
}

// Delete marks an exhibition as trash
func (s *exhibitionSvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

// GetRelatedExhibitions retrieves related exhibitions for a given exhibition
func (s *exhibitionSvcs) GetRelatedExhibitions(ctx context.Context, id string) ([]models.Exhibition, error) {
	exhibition, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(exhibition.RelatedExhibitions) == 0 {
		return []models.Exhibition{}, nil
	}

	relatedExhibitions, err := s.getExhibitionsByIds(ctx, exhibition.RelatedExhibitions)
	if err != nil {
		return nil, err
	}

	return relatedExhibitions, nil
}

// Ad Management Methods

// GetExhibitionAds retrieves all ads for an exhibition
func (s *exhibitionSvcs) GetExhibitionAds(ctx context.Context, exhibitionId string) ([]models.Ad, error) {
	exhibition, err := s.GetById(ctx, exhibitionId)
	if err != nil {
		return nil, err
	}

	return exhibition.Ads, nil
}

// UpdateAd updates an existing ad
func (s *exhibitionSvcs) UpdateAd(ctx context.Context, exhibitionId, adId string, data *models.AdDto) (*models.Ad, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	exhibitionObjId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return nil, helpers.BadRequest("Invalid exhibition ID")
	}

	adObjId, err := primitive.ObjectIDFromHex(adId)
	if err != nil {
		return nil, helpers.BadRequest("Invalid ad ID")
	}

	// Update specific ad in the ads array
	filter := bson.M{
		"_id":     exhibitionObjId,
		"ads._id": adObjId,
	}

	update := bson.M{
		"$set": bson.M{
			"ads.$.title":          data.Title,
			"ads.$.description":    data.Description,
			"ads.$.tags":           data.Tags,
			"ads.$.images":         data.Images,
			"ads.$.websiteUrl":     data.WebsiteUrl,
			"ads.$.packageType":    data.PackageType,
			"ads.$.duration":       data.Duration,
			"ads.$.paymentMethod":  data.PaymentMethod,
			"ads.$.paymentStatus":  data.PaymentStatus,
			"ads.$.activationDate": data.ActivationDate,
			"ads.$.expirationDate": data.ExpirationDate,
			"ads.$.price":          data.Price,
			"ads.$.isActive":       data.IsActive,
			"ads.$.updatedAt":      time.Now(),
			"ads.$.updatedBy":      cfg.User.Id,
			"updatedAt":            time.Now(),
			"updatedBy":            cfg.User.Id,
		},
	}

	updatedExhibition, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updatedExhibition == nil {
		return nil, helpers.NotFoundError("Ad not found")
	}

	// Find and return the updated ad
	for _, ad := range updatedExhibition.Ads {
		if ad.Id == adObjId {
			return &ad, nil
		}
	}

	return nil, helpers.NotFoundError("Ad not found")
}

// DeleteAd removes an ad from an exhibition
func (s *exhibitionSvcs) DeleteAd(ctx context.Context, exhibitionId, adId string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	exhibitionObjId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	adObjId, err := primitive.ObjectIDFromHex(adId)
	if err != nil {
		return helpers.BadRequest("Invalid ad ID")
	}

	// First check if exhibition exists
	_, err = s.GetById(ctx, exhibitionId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": exhibitionObjId}
	update := bson.M{
		"$pull": bson.M{"ads": bson.M{"_id": adObjId}},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// ToggleAd enables or disables an ad
func (s *exhibitionSvcs) ToggleAd(ctx context.Context, exhibitionId, adId string, isActive bool) (*models.Ad, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	exhibitionObjId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return nil, helpers.BadRequest("Invalid exhibition ID")
	}

	adObjId, err := primitive.ObjectIDFromHex(adId)
	if err != nil {
		return nil, helpers.BadRequest("Invalid ad ID")
	}

	filter := bson.M{
		"_id":     exhibitionObjId,
		"ads._id": adObjId,
	}

	update := bson.M{
		"$set": bson.M{
			"ads.$.isActive":  isActive,
			"ads.$.updatedAt": time.Now(),
			"ads.$.updatedBy": cfg.User.Id,
			"updatedAt":       time.Now(),
			"updatedBy":       cfg.User.Id,
		},
	}

	updatedExhibition, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updatedExhibition == nil {
		return nil, helpers.NotFoundError("Ad not found")
	}

	// Find and return the updated ad
	for _, ad := range updatedExhibition.Ads {
		if ad.Id == adObjId {
			return &ad, nil
		}
	}

	return nil, helpers.NotFoundError("Ad not found")
}

// Helper methods

// validateRelatedExhibitions checks if all related exhibition IDs exist
func (s *exhibitionSvcs) validateRelatedExhibitions(ctx context.Context, relatedExhibitionIds []primitive.ObjectID) error {
	for _, id := range relatedExhibitionIds {
		filter := bson.M{"_id": id, "trash": false}
		_, err := s.repo.GetByFilter(ctx, filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return helpers.BadRequest("One or more related exhibitions not found")
			}
			return err
		}
	}
	return nil
}

// getExhibitionsByIds retrieves multiple exhibitions by their IDs
func (s *exhibitionSvcs) getExhibitionsByIds(ctx context.Context, ids []primitive.ObjectID) ([]models.Exhibition, error) {
	filter := bson.M{
		"_id":   bson.M{"$in": ids},
		"trash": false,
	}

	var exhibitions []models.Exhibition
	err := s.repo.Aggregate(ctx, []bson.M{{"$match": filter}}, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &exhibitions)
	})

	return exhibitions, err
}

// ValidateExhibitionExists checks if an exhibition exists and is not trashed
func (s *exhibitionSvcs) ValidateExhibitionExists(ctx context.Context, exhibitionId string) error {
	objID, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	filter := bson.M{"_id": objID, "trash": false}
	_, err = s.repo.GetByFilter(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return helpers.NotFoundError("Exhibition not found or is trashed")
		}
		return err
	}

	return nil
}

// Debug methods - remove in production
func (s *exhibitionSvcs) DebugCount(ctx context.Context) (map[string]interface{}, error) {
	totalCount, err := s.repo.Count(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	activeCount, err := s.repo.Count(ctx, bson.M{"trash": false})
	if err != nil {
		return nil, err
	}

	trashedCount, err := s.repo.Count(ctx, bson.M{"trash": true})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":   totalCount,
		"active":  activeCount,
		"trashed": trashedCount,
	}, nil
}

func (s *exhibitionSvcs) DebugRaw(ctx context.Context) ([]models.Exhibition, error) {
	pipeline := []bson.M{
		{"$match": bson.M{}}, // Get all records including trashed
		{"$limit": 10},       // Limit to first 10 records
	}

	var result []models.Exhibition
	err := s.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	return result, err
}

// GetAvailableFilterFields returns the available fields for filtering
func GetAvailableFilterFields() map[string]string {
	return map[string]string{
		// Exhibition fields
		"_id":                "ObjectID of the exhibition",
		"title":              "Exhibition title (string)",
		"description":        "Exhibition description (string)",
		"locationType":       "Type of location (string)",
		"locationAddress":    "Address of the location (string)",
		"startDate":          "Start date (time)",
		"endDate":            "End date (time)",
		"images":             "Array of image file fields",
		"relatedExhibitions": "Array of related exhibition ObjectIDs",
		"tags":               "Array of tags (strings)",
		"isActive":           "Active status (boolean)",
		"ads":                "Array of ads",
		"trash":              "Trash status (boolean)",
		"createdAt":          "Creation date (time)",
		"createdBy":          "Creator ObjectID",
		"updatedAt":          "Update date (time)",
		"updatedBy":          "Updater ObjectID",
	}
}
