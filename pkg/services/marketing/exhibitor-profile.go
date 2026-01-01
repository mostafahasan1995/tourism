package marketing

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/services/marketing/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExhibitorProfileSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ExhibitorProfile, error)
	GetByHotelId(ctx context.Context, hotelId string) (*models.ExhibitorProfile, error)
	Get(ctx context.Context, skip, limit int64, filterQuery filter.ExhibitorProfileFilter) (*models.ExhibitorProfilePagination, error)
	GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.ExhibitorProfilePagination, error)
	Add(ctx context.Context, data *models.ExhibitorProfileDto) (*models.ExhibitorProfile, error)
	AddExhibitorRequest(ctx context.Context, data *models.ExhibitorRequestDto) (*models.ExhibitorProfile, error)
	Update(ctx context.Context, id string, data *models.ExhibitorProfileDto) (*models.ExhibitorProfile, error)
	Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.ExhibitorProfile, error)
	Delete(ctx context.Context, id string) error

	GetExhibitorRequests(ctx context.Context, skip, limit int64, filterQuery filter.ExhibitorRequestFilter) (*models.ExhibitorRequestPagination, error)
	UpdateExhibitorRequestStatus(ctx context.Context, id string, status string) (*models.ExhibitorProfile, error)

	AddDynamicSection(ctx context.Context, profileId string, section *models.AddDynamicSectionDto) (*models.ExhibitorProfile, error)
	UpdateDynamicSection(ctx context.Context, profileId string, sectionId string, updates *models.UpdateDynamicSectionDto) (*models.ExhibitorProfile, error)
	DeleteDynamicSection(ctx context.Context, profileId string, sectionId string) (*models.ExhibitorProfile, error)
	ReorderDynamicSections(ctx context.Context, profileId string, sectionOrders map[string]int) (*models.ExhibitorProfile, error)

	UpdateFacility(ctx context.Context, profileId string, facilityUpdate *models.FacilityUpdateDto) (*models.ExhibitorProfile, error)
	UpdateHeroSection(ctx context.Context, profileId string, heroSection *models.HeroSection) (*models.ExhibitorProfile, error)

	Publish(ctx context.Context, profileId string) (*models.ExhibitorProfile, error)
	Unpublish(ctx context.Context, profileId string) (*models.ExhibitorProfile, error)
	ToggleActive(ctx context.Context, profileId string) (*models.ExhibitorProfile, error)
}

type exhibitorProfileSvcs struct {
	repo repo.ExhibitorProfileRepo
}

func NewExhibitorProfileSvcs(i *do.Injector) (ExhibitorProfileSvcs, error) {
	return &exhibitorProfileSvcs{
		repo: do.MustInvoke[repo.ExhibitorProfileRepo](i),
	}, nil
}

func (s *exhibitorProfileSvcs) GetOne(ctx context.Context, id string) (*models.ExhibitorProfile, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (s *exhibitorProfileSvcs) GetByHotelId(ctx context.Context, hotelId string) (*models.ExhibitorProfile, error) {
	_hotelId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	return s.repo.GetByFilter(ctx, bson.M{"hotelId": _hotelId, "trash": false})
}

func (s *exhibitorProfileSvcs) Get(ctx context.Context, skip, limit int64, filterQuery filter.ExhibitorProfileFilter) (*models.ExhibitorProfilePagination, error) {
	match := filterQuery.ToBsonFilter()

	countPipeline := []bson.M{{"$match": match}}
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var result []models.ExhibitorProfile
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

	return &models.ExhibitorProfilePagination{
		Profiles:   result,
		Pagination: pagination,
	}, nil
}

// GetV2 retrieves exhibitor profiles with pagination and advanced filters
func (s *exhibitorProfileSvcs) GetV2(ctx context.Context, skip, limit int64, filters *query.Conditions) (*models.ExhibitorProfilePagination, error) {
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

	count, err := s.repo.Count(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.ExhibitorProfile
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

	return &models.ExhibitorProfilePagination{
		Profiles:   result,
		Pagination: pagination,
	}, nil
}

func (s *exhibitorProfileSvcs) Add(ctx context.Context, data *models.ExhibitorProfileDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	existing, _ := s.GetByHotelId(ctx, data.HotelId.Hex())
	if existing != nil {
		return nil, helpers.BadRequest("Exhibitor profile already exists for this hotel")
	}

	profile := &models.ExhibitorProfile{
		ExhibitorProfileDto: *data,
		Id:                  primitive.NewObjectID(),
		Trash:               false,
		CreatedAt:           time.Now(),
		CreatedBy:           cfg.User.Id,
		UpdatedAt:           time.Now(),
		UpdatedBy:           cfg.User.Id,
	}

	if err := s.repo.Add(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *exhibitorProfileSvcs) AddExhibitorRequest(ctx context.Context, data *models.ExhibitorRequestDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	existing, _ := s.GetByHotelId(ctx, data.HotelId.Hex())
	if existing != nil {
		return nil, helpers.BadRequest("Exhibitor profile already exists for this hotel")
	}

	// Create contact info map
	contactInfo := map[string]interface{}{
		"website":  data.HotelWebsite,
		"phone":    data.Phone,
		"email":    data.Email,
		"location": data.Location,
	}

	// Initialize a full profile with provided data
	profileDto := &models.ExhibitorProfileDto{
		ExhibitionId: data.ExhibitionId,
		HotelId:      data.HotelId,
		HeroSection: models.HeroSection{

			Rating: 0,

			Overview: data.Overview,
			Images:   []types.FileField{}, // Empty array for images to be added later
			Logo:     nil,                 // Null logo to be added later
		},
		FacilitiesSection: models.GetDefaultFacilityOptions(),
		DynamicSections:   []models.DynamicSection{},
		ContactInfo:       contactInfo,
		IsActive:          false,
		IsPublished:       false,
	}

	profile := &models.ExhibitorProfile{
		ExhibitorProfileDto: *profileDto,
		Id:                  primitive.NewObjectID(),
		Trash:               false,
		CreatedAt:           time.Now(),
		CreatedBy:           cfg.User.Id,
		UpdatedAt:           time.Now(),
		UpdatedBy:           cfg.User.Id,
		Status:              "Pending", // Set initial status as Pending
	}

	if err := s.repo.Add(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *exhibitorProfileSvcs) Update(ctx context.Context, id string, data *models.ExhibitorProfileDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	existing, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &models.ExhibitorProfile{
		ExhibitorProfileDto: *data,
		Id:                  _id,
		Trash:               false,
		CreatedAt:           existing.CreatedAt,
		CreatedBy:           existing.CreatedBy,
		UpdatedBy:           cfg.User.Id,
		UpdatedAt:           time.Now(),
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": profile}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *exhibitorProfileSvcs) Patch(ctx context.Context, id string, updates map[string]interface{}) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	updateDoc := bson.M{
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}

	for key, value := range updates {
		updateDoc[key] = value
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateDoc}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, id)
}

func (s *exhibitorProfileSvcs) Delete(ctx context.Context, id string) error {
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

	_, err = s.repo.Patch(ctx, filter, update)
	return err
}

func (s *exhibitorProfileSvcs) AddDynamicSection(ctx context.Context, profileId string, section *models.AddDynamicSectionDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.GetOne(ctx, profileId)
	if err != nil {
		return nil, err
	}

	newSection := models.DynamicSection{
		Id:       primitive.NewObjectID(),
		Title:    section.Title,
		Overview: section.Overview,
		Images:   section.Images,
		Order:    len(profile.DynamicSections) + 1,
	}

	_id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": _id}
	update := bson.M{
		"$push": bson.M{"dynamicSections": newSection},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) UpdateDynamicSection(ctx context.Context, profileId string, sectionId string, updates *models.UpdateDynamicSectionDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	_sectionId, err := primitive.ObjectIDFromHex(sectionId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	updateFields := bson.M{}
	if updates.Title != nil {
		updateFields["dynamicSections.$.title"] = *updates.Title
	}
	if updates.Overview != nil {
		updateFields["dynamicSections.$.overview"] = *updates.Overview
	}
	if updates.Images != nil {
		updateFields["dynamicSections.$.images"] = *updates.Images
	}
	if updates.Order != nil {
		updateFields["dynamicSections.$.order"] = *updates.Order
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = cfg.User.Id

	filter := bson.M{"_id": _id, "dynamicSections._id": _sectionId}
	update := bson.M{"$set": updateFields}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) DeleteDynamicSection(ctx context.Context, profileId string, sectionId string) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	_sectionId, err := primitive.ObjectIDFromHex(sectionId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{
		"$pull": bson.M{"dynamicSections": bson.M{"_id": _sectionId}},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) ReorderDynamicSections(ctx context.Context, profileId string, sectionOrders map[string]int) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.GetOne(ctx, profileId)
	if err != nil {
		return nil, err
	}

	for i, section := range profile.DynamicSections {
		if order, exists := sectionOrders[section.Id.Hex()]; exists {
			profile.DynamicSections[i].Order = order
		}
	}

	_id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"dynamicSections": profile.DynamicSections,
		"updatedAt":       time.Now(),
		"updatedBy":       cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) UpdateFacility(ctx context.Context, profileId string, facilityUpdate *models.FacilityUpdateDto) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.GetOne(ctx, profileId)
	if err != nil {
		return nil, err
	}

	// Update the facility in memory
	var facilityList *[]models.FacilityOption
	switch facilityUpdate.Category {
	case "restaurants":
		facilityList = &profile.FacilitiesSection.Restaurants
	case "poolsBeaches":
		facilityList = &profile.FacilitiesSection.PoolsBeaches
	case "spaGym":
		facilityList = &profile.FacilitiesSection.SpaGym
	case "hotelServices":
		facilityList = &profile.FacilitiesSection.HotelServices
	case "business":
		facilityList = &profile.FacilitiesSection.Business
	case "kidsFacilities":
		facilityList = &profile.FacilitiesSection.KidsFacilities
	case "recreational":
		facilityList = &profile.FacilitiesSection.Recreational
	default:
		return nil, helpers.BadRequest("Invalid facility category")
	}

	// Find and update the facility
	found := false
	for i, facility := range *facilityList {
		if facility.Id == facilityUpdate.FacilityId {
			(*facilityList)[i].IsSelected = facilityUpdate.IsSelected
			found = true
			break
		}
	}

	if !found {
		return nil, helpers.BadRequest("Facility not found")
	}

	_id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"facilitiesSection": profile.FacilitiesSection,
		"updatedAt":         time.Now(),
		"updatedBy":         cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) UpdateHeroSection(ctx context.Context, profileId string, heroSection *models.HeroSection) (*models.ExhibitorProfile, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		return nil, helpers.InvalidObjectId()
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"heroSection": heroSection,
		"updatedAt":   time.Now(),
		"updatedBy":   cfg.User.Id,
	}}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, profileId)
}

func (s *exhibitorProfileSvcs) Publish(ctx context.Context, profileId string) (*models.ExhibitorProfile, error) {
	return s.Patch(ctx, profileId, map[string]interface{}{"isPublished": true})
}

func (s *exhibitorProfileSvcs) Unpublish(ctx context.Context, profileId string) (*models.ExhibitorProfile, error) {
	return s.Patch(ctx, profileId, map[string]interface{}{"isPublished": false})
}

func (s *exhibitorProfileSvcs) ToggleActive(ctx context.Context, profileId string) (*models.ExhibitorProfile, error) {
	profile, err := s.GetOne(ctx, profileId)
	if err != nil {
		return nil, err
	}

	return s.Patch(ctx, profileId, map[string]interface{}{"isActive": !profile.IsActive})
}

func (s *exhibitorProfileSvcs) GetExhibitorRequests(ctx context.Context, skip, limit int64, filterQuery filter.ExhibitorRequestFilter) (*models.ExhibitorRequestPagination, error) {
	match := filterQuery.ToBsonFilter()

	// Add status condition if not already specified
	if _, ok := match["status"]; !ok {
		match["status"] = bson.M{"$exists": true}
	}

	countPipeline := []bson.M{{"$match": match}}
	count, err := s.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}

	var result []models.ExhibitorProfile
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

	return &models.ExhibitorRequestPagination{
		Profiles:   result,
		Pagination: pagination,
	}, nil
}

func (s *exhibitorProfileSvcs) UpdateExhibitorRequestStatus(ctx context.Context, id string, status string) (*models.ExhibitorProfile, error) {
	// Validate status
	if status != "Pending" && status != "Replied" && status != "Closed" {
		return nil, helpers.BadRequest("Invalid status value. Must be 'Pending', 'Replied', or 'Closed'")
	}

	return s.Patch(ctx, id, map[string]interface{}{"status": status})
}
