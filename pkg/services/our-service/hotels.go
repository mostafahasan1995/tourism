package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"log"
	"math"
	"time"

	"github.com/goccy/go-json"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Hotels, error)
	GetAll(ctx context.Context, query any, page, perPage int64) (*models.HotelsPaginationRes, error)
	GetAllHotels(ctx context.Context, query any) ([]models.Hotels, error)
	GetAuth(ctx context.Context, query any, page, perPage int64) (*models.HotelsPaginationRes, error)
	GetAllAuth(ctx context.Context, query any) ([]models.Hotels, error)
	Add(ctx context.Context, data *models.HotelsDto) (*models.Hotels, error)
	Update(ctx context.Context, id string, data *models.HotelsDto) (*models.Hotels, error)
	UpdateIsFav(ctx context.Context, id string, isFav bool) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.HotelsPaginationRes, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Hotels, error)
}

type hotelsSvcs struct {
	repo repo.HotelsRepo
}

func NewHotelsSvcs(i *do.Injector) (HotelsSvcs, error) {
	return &hotelsSvcs{
		repo: do.MustInvoke[repo.HotelsRepo](i),
	}, nil
}

// Helper method to convert Hotels to HotelsRes with isFav populated
func (h *hotelsSvcs) convertToHotelsRes(ctx context.Context, hotels []models.Hotels) ([]models.HotelsRes, error) {
	if len(hotels) == 0 {
		return []models.HotelsRes{}, nil
	}

	// Convert to HotelsRes - isFav is now stored directly in the entity
	result := make([]models.HotelsRes, len(hotels))
	for i, hotel := range hotels {
		result[i] = models.HotelsRes{
			Hotels: hotel,
			IsFav:  hotel.IsFav, // Use the isFav field directly from the entity
		}
	}

	return result, nil
}

func (h *hotelsSvcs) GetOne(ctx context.Context, id string) (*models.Hotels, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return h.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (h *hotelsSvcs) GetAllHotels(ctx context.Context, query any) ([]models.Hotels, error) {
	// Log the incoming query
	if jsonBytes, err := json.Marshal(query); err == nil {
		log.Printf("Hotel service received query: %s", string(jsonBytes))
	}

	// Handle different types of query input
	var hotelFilter filter.HotelsFilter

	switch q := query.(type) {
	case filter.HotelsFilter:
		// Direct struct
		hotelFilter = q
	case *filter.HotelsFilter:
		// Pointer to struct
		if q != nil {
			hotelFilter = *q
		}
	case string:
		// JSON string
		if q != "" {
			if err := json.Unmarshal([]byte(q), &hotelFilter); err != nil {
				log.Printf("Error unmarshaling filter from string: %v", err)
				return nil, errors.New("invalid query format")
			}
		}
	case map[string]interface{}:
		// Map (from JSON)
		jsonBytes, err := json.Marshal(q)
		if err != nil {
			log.Printf("Error marshaling map to JSON: %v", err)
			return nil, errors.New("invalid query format")
		}
		if err := json.Unmarshal(jsonBytes, &hotelFilter); err != nil {
			log.Printf("Error unmarshaling filter from map: %v", err)
			return nil, errors.New("invalid query format")
		}
	default:
		// Try using the ParseFilters helper as fallback
		parsed, err := helpers.ParseFilters[filter.HotelsFilter](query)
		if err != nil {
			log.Printf("Error parsing filter with helper: %v", err)
			return nil, errors.New("invalid query: " + err.Error())
		}
		hotelFilter = *parsed
	}

	// Validate the price range
	if err := hotelFilter.PriceRange.Validate(); err != nil {
		log.Printf("Price range validation failed: %v", err)
		return nil, errors.New("invalid price range: " + err.Error())
	}

	// Special handling for empty arrays that should be nil
	if len(hotelFilter.HotelTypes) == 0 {
		hotelFilter.HotelTypes = nil
	}
	if len(hotelFilter.RoomAmenities) == 0 {
		hotelFilter.RoomAmenities = nil
	}
	if len(hotelFilter.NearbyAttractions) == 0 {
		hotelFilter.NearbyAttractions = nil
	}
	if len(hotelFilter.Locations) == 0 {
		hotelFilter.Locations = nil
	}

	// Log the filter being used
	if jsonBytes, err := json.Marshal(hotelFilter); err == nil {
		log.Printf("Using hotel filter: %s", string(jsonBytes))
	}

	// Call the repository with the parsed filter
	result, err := h.repo.GetAllHotels(ctx, hotelFilter)
	if err != nil {
		log.Printf("Repository error: %v", err)
		return nil, err
	}

	return result, nil
}

func (h *hotelsSvcs) GetAll(ctx context.Context, query any, page, perPage int64) (*models.HotelsPaginationRes, error) {
	// Log the incoming query
	if jsonBytes, err := json.Marshal(query); err == nil {
		log.Printf("Hotel service received query: %s", string(jsonBytes))
	}

	// Handle different types of query input
	var hotelFilter filter.HotelsFilter

	switch q := query.(type) {
	case filter.HotelsFilter:
		// Direct struct
		hotelFilter = q
	case *filter.HotelsFilter:
		// Pointer to struct
		if q != nil {
			hotelFilter = *q
		}
	case string:
		// JSON string
		if q != "" {
			if err := json.Unmarshal([]byte(q), &hotelFilter); err != nil {
				log.Printf("Error unmarshaling filter from string: %v", err)
				return nil, errors.New("invalid query format")
			}
		}
	case map[string]interface{}:
		// Map (from JSON)
		jsonBytes, err := json.Marshal(q)
		if err != nil {
			log.Printf("Error marshaling map to JSON: %v", err)
			return nil, errors.New("invalid query format")
		}
		if err := json.Unmarshal(jsonBytes, &hotelFilter); err != nil {
			log.Printf("Error unmarshaling filter from map: %v", err)
			return nil, errors.New("invalid query format")
		}
	default:
		// Try using the ParseFilters helper as fallback
		parsed, err := helpers.ParseFilters[filter.HotelsFilter](query)
		if err != nil {
			log.Printf("Error parsing filter with helper: %v", err)
			return nil, errors.New("invalid query: " + err.Error())
		}
		hotelFilter = *parsed
	}

	// Validate the price range
	if err := hotelFilter.PriceRange.Validate(); err != nil {
		log.Printf("Price range validation failed: %v", err)
		return nil, errors.New("invalid price range: " + err.Error())
	}

	// Build aggregation pipeline
	skip := (page - 1) * perPage
	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
	}

	// Add filter stages based on hotelFilter
	if filterStages := hotelFilter.BuildPipeline(bson.M{}); len(filterStages) > 0 {
		pipeline = append(pipeline, filterStages...)
	}

	// Add fave lookup - check if user is authenticated
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		// User is authenticated, add lookup to check favorites
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"hotelId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$hotelId"}},
									{"$eq": []interface{}{"$type", "hotel"}},
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

	// Count for pagination (exclude skip/limit from count pipeline)
	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := h.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	// Add sorting, skip, and limit
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": perPage})

	var hotels []models.Hotels
	errAg := h.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &hotels)
	})
	if errAg != nil {
		return nil, errAg
	}

	// Convert to HotelsRes with isFav populated from pipeline
	hotelsRes := make([]models.HotelsRes, len(hotels))
	for i, hotel := range hotels {
		hotelsRes[i] = models.HotelsRes{
			Hotels: hotel,
			IsFav:  hotel.IsFav,
		}
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(perPage))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    perPage,
		TotalCount: count,
	}

	result := &models.HotelsPaginationRes{
		Hotels:     hotelsRes,
		Pagination: pagination,
	}

	log.Printf("Found %d hotels", len(hotelsRes))
	return result, nil
}

// GetAuth retrieves hotels with pagination, requiring authentication
func (h *hotelsSvcs) GetAuth(ctx context.Context, query any, page, perPage int64) (*models.HotelsPaginationRes, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as GetAll but ensure user is authenticated
	return h.GetAll(ctx, query, page, perPage)
}

// GetAllAuth retrieves all hotels without pagination, requiring authentication
func (h *hotelsSvcs) GetAllAuth(ctx context.Context, query any) ([]models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as GetAllHotels but ensure user is authenticated
	return h.GetAllHotels(ctx, query)
}

func (h *hotelsSvcs) Add(ctx context.Context, data *models.HotelsDto) (*models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	hotels := &models.Hotels{
		HotelsDto: *data,
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	if err := h.repo.Add(ctx, hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (h *hotelsSvcs) Update(ctx context.Context, id string, data *models.HotelsDto) (*models.Hotels, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return h.repo.Update(ctx, _id, data)
}

func (h *hotelsSvcs) UpdateIsFav(ctx context.Context, id string, isFav bool) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"isFav":     isFav,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = h.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (h *hotelsSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = h.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (h *hotelsSvcs) Count(ctx context.Context, filter any) (int64, error) {
	return h.repo.Count(ctx, filter)
}

// v2
func (h *hotelsSvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.HotelsPaginationRes, error) {
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

	// Add fave lookup - check if user is authenticated
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		// User is authenticated, add lookup to check favorites
		pipeline = append(pipeline, bson.M{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"hotelId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$hotelId"}},
									{"$eq": []interface{}{"$type", "hotel"}},
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

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := h.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Hotels
	errAg := h.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	for i := range result {
		result[i].CalculateAverageRating()
	}

	// Convert to HotelsRes with isFav populated from pipeline
	hotelsRes := make([]models.HotelsRes, len(result))
	for i, hotel := range result {
		hotelsRes[i] = models.HotelsRes{
			Hotels: hotel,
			IsFav:  hotel.IsFav,
		}
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.HotelsPaginationRes{
		Hotels:     hotelsRes,
		Pagination: pagination,
	}, nil
}

func (h *hotelsSvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Hotels, error) {
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

	var result []models.Hotels
	errAg := h.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	for i := range result {
		result[i].CalculateAverageRating()
	}

	return result, nil
}
