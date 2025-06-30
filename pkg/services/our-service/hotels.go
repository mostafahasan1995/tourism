package ourservice

import (
	"context"
	"errors"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/util"
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Hotels, error)
	GetAll(ctx context.Context, query any, page, perPage int64) (*models.HotelsPagination, error)
	GetAllHotels(ctx context.Context, query any) ([]models.Hotels, error)
	Add(ctx context.Context, data *models.HotelsDto) (*models.Hotels, error)
	Update(ctx context.Context, id string, data *models.HotelsDto) (*models.Hotels, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
}

type hotelsSvcs struct {
	repo repo.HotelsRepo
}

func NewHotelsSvcs(i *do.Injector) (HotelsSvcs, error) {
	return &hotelsSvcs{
		repo: do.MustInvoke[repo.HotelsRepo](i),
	}, nil
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

func (h *hotelsSvcs) GetAll(ctx context.Context, query any, page, perPage int64) (*models.HotelsPagination, error) {
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
	result, err := h.repo.GetAll(ctx, hotelFilter, page, perPage)
	if err != nil {
		log.Printf("Repository error: %v", err)
		return nil, err
	}

	log.Printf("Found %d hotels", len(result.Hotels))
	return &result, nil
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
