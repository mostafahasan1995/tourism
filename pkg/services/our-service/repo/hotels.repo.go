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
		size = int(totalCount) // return all if invalid
	}
	skip := int64((page - 1) * size)
	limit := int64(size)

	// Query options with pagination
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)

	cur, err := coll.Find(ctx, filterBody, findOptions)
	if err != nil {
		return models.HotelsPagination{}, err
	}

	var programs []models.Hotels
	if err := cur.All(ctx, &programs); err != nil {
		return models.HotelsPagination{}, err
	}

	// Prepare pagination result
	totalPages := float64(0)
	if size > 0 {
		totalPages = float64((totalCount + int64(size) - 1) / int64(size))
	}

	for _, p := range programs {
		p.CalculateAverageRating()
	}
	result := models.HotelsPagination{
		Hotels: programs,
		Pagination: common.Pagination{
			TotalPages: totalPages,
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
