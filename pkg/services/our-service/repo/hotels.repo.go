package repo

import (
	"context"
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"

	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelsRepo interface {
	dbrepo.MainRepo[models.Hotels]
	GetOne(ctx context.Context, id string) (*models.Hotels, error)
	GetAll(ctx context.Context, filter filter.HotelsFilter) (models.HotelsPagination, error)
	Update(ctx context.Context, id primitive.ObjectID, data *models.HotelsDto) error
	Delete(ctx context.Context, id string) error
	// Hotel Reviews
	AddReview(ctx context.Context, review *models.HotelReview) error
	GetHotelReviews(ctx context.Context, hotelId string, page, size int) (models.HotelReviewPagination, error)
	GetReview(ctx context.Context, reviewId string) (*models.HotelReview, error)
	// Debug method
	GetAllReviews(ctx context.Context) ([]models.HotelReview, error)
	// Review management
	UpdateReviewStatus(ctx context.Context, reviewId string, status string) error
}

type hotelsrepo struct {
	dbrepo.MainRepoImpl[models.Hotels]
	db       *mongo.Client
	collName string
}

func NewHotelsRepo(i *do.Injector) (HotelsRepo, error) {
	return &hotelsrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Hotels]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismHotels",
		},
		db:       do.MustInvoke[*mongo.Client](i),
		collName: "tourismHotels",
	}, nil
}

func (l *hotelsrepo) GetOne(ctx context.Context, id string) (*models.Hotels, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll := l.db.Database(cfg.Db).Collection(l.collName)

	var data models.Hotels
	if err := coll.FindOne(ctx, bson.M{"_id": _id, "trash": false}).Decode(&data); err != nil {
		return nil, err
	}
	data.CalculateAverageRating()
	return &data, nil

}

func (l *hotelsrepo) GetAll(ctx context.Context, filter filter.HotelsFilter) (models.HotelsPagination, error) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.HotelsPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filterBody := filter.ToBsonFilter()

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filterBody)
	if err != nil {
		return models.HotelsPagination{}, err
	}

	// Pagination defaults and limits
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.Size
	if size <= 0 {
		size = 10 // Default page size instead of returning all
	}
	if size > 100 {
		size = 100 // Maximum page size limit
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Query options with pagination and sorting
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.HotelsPagination{}, err
	}

	var hotels []models.Hotels
	if err := cur.All(ctx, &hotels); err != nil {
		return models.HotelsPagination{}, err
	}

	// Calculate average rating for each hotel
	for i := range hotels {
		hotels[i].CalculateAverageRating()
	}

	// Prepare pagination result
	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	result := models.HotelsPagination{
		Hotels: hotels,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

func (l *hotelsrepo) Update(ctx context.Context, id primitive.ObjectID, data *models.HotelsDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	preHotels, err := l.GetOne(ctx, id.Hex())
	if err != nil {
		return err
	}

	hotels := &models.Hotels{
		HotelsDto: models.HotelsDto{
			Name:                          data.Name,
			HotelType:                     data.HotelType,
			Location:                      data.Location,
			CheckInAndCheckOut:            data.CheckInAndCheckOut,
			Price:                         data.Price,
			IsDisplayInPerfectStay:        data.IsDisplayInPerfectStay,
			Ratings:                       data.Ratings,
			Image:                         data.Image,
			RoomAmenities:                 data.RoomAmenities,
			DistanceFromCityCenter:        data.DistanceFromCityCenter,
			NearbyAttractions:             data.NearbyAttractions,
			ImagesGallery:                 data.ImagesGallery,
			OverviewPage:                  data.OverviewPage,
			RoomsAndSuitesPage:            data.RoomsAndSuitesPage,
			AmenitiesAndFacilitiesPage:    data.AmenitiesAndFacilitiesPage,
			LocationNearbyAttractionsPage: data.LocationNearbyAttractionsPage,
			ReviewsAndRatingsPage:         data.ReviewsAndRatingsPage,
			BookingAndPoliciesPage:        data.BookingAndPoliciesPage,
			PositionOnMap:                 data.PositionOnMap,
			Contacts:                      data.Contacts,
			CloseReservations:             data.CloseReservations,
			RatingObjects:                 data.RatingObjects,
			OfferAndDiscount:              data.OfferAndDiscount,
			PoliciesPage:                  data.PoliciesPage,
			AmenitiesAndFacilitiesPageV2:  data.AmenitiesAndFacilitiesPageV2,
		},

		Id:        id,
		Trash:     false,
		CreatedAt: preHotels.CreatedAt,
		CreatedBy: preHotels.CreatedBy,
		UpdatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": hotels}

	upsert := false
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updatedHotels models.Hotels
	if err := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedHotels); err != nil {
		return err
	}

	return nil
}

func (l *hotelsrepo) Delete(ctx context.Context, id string) error {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection(l.collName)

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": primitive.NilObjectID,
	}}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	result := coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&opt,
	)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// AddReview adds a new hotel review
func (l *hotelsrepo) AddReview(ctx context.Context, review *models.HotelReview) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection("hotelReviews")

	review.Id = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	review.Trash = false

	if review.Status == "" {
		review.Status = "pending"
	}

	if review.Date.IsZero() {
		review.Date = time.Now()
	}

	_, err = coll.InsertOne(ctx, review)
	return err
}

// GetHotelReviews retrieves reviews for a specific hotel with pagination
func (l *hotelsrepo) GetHotelReviews(ctx context.Context, hotelId string, page, size int) (models.HotelReviewPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return models.HotelReviewPagination{}, err
	}

	hotelObjectId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return models.HotelReviewPagination{}, err
	}

	coll := l.db.Database(cfg.Db).Collection("hotelReviews")

	// Create filter for hotel reviews
	filter := bson.M{
		"hotelId": hotelObjectId,
		"trash":   bson.M{"$ne": true},
	}

	// Count total documents matching the filter
	totalCount, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return models.HotelReviewPagination{}, err
	}

	// Pagination defaults and limits
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10 // Default page size
	}
	if size > 100 {
		size = 100 // Maximum page size limit
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Query options with pagination and sorting
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	cur, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return models.HotelReviewPagination{}, err
	}

	var reviews []models.HotelReview
	if err := cur.All(ctx, &reviews); err != nil {
		return models.HotelReviewPagination{}, err
	}

	// Prepare pagination result
	totalPages := int64(0)
	if size > 0 {
		totalPages = (totalCount + int64(size) - 1) / int64(size)
	}

	result := models.HotelReviewPagination{
		Reviews: reviews,
		Pagination: common.Pagination{
			TotalPages: float64(totalPages),
			PerPage:    int64(size),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

// GetReview retrieves a specific review by ID
func (l *hotelsrepo) GetReview(ctx context.Context, reviewId string) (*models.HotelReview, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	reviewObjectId, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return nil, err
	}

	coll := l.db.Database(cfg.Db).Collection("hotelReviews")

	var review models.HotelReview
	if err := coll.FindOne(ctx, bson.M{"_id": reviewObjectId, "trash": false}).Decode(&review); err != nil {
		return nil, err
	}

	return &review, nil
}

// GetAllReviews retrieves all reviews in the database
func (l *hotelsrepo) GetAllReviews(ctx context.Context) ([]models.HotelReview, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	coll := l.db.Database(cfg.Db).Collection("hotelReviews")

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var reviews []models.HotelReview
	if err := cur.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

// UpdateReviewStatus updates the status of a specific review
func (l *hotelsrepo) UpdateReviewStatus(ctx context.Context, reviewId string, status string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	reviewObjectId, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return err
	}

	coll := l.db.Database(cfg.Db).Collection("hotelReviews")

	filter := bson.M{"_id": reviewObjectId}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err = coll.UpdateOne(ctx, filter, update)
	return err
}
