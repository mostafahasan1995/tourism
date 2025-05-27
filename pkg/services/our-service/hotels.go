package ourservice

import (
	"context"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelsSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Hotels, error)
	GetAll(ctx context.Context, filter filter.HotelsFilter) (models.HotelsPagination, error)
	Add(ctx context.Context, data *models.HotelsDto) (any, error)
	Update(ctx context.Context, id string, data *models.HotelsDto) error
	Delete(ctx context.Context, id string) error
}

type hotelssvcs struct {
	repo repo.HotelsRepo
}

func NewHotelsSvcs(i *do.Injector) (HotelsSvcs, error) {
	return &hotelssvcs{
		repo: do.MustInvoke[repo.HotelsRepo](i),
	}, nil
}

func (l *hotelssvcs) GetOne(ctx context.Context, id string) (*models.Hotels, error) {
	return l.repo.GetOne(ctx, id)

}

func (l *hotelssvcs) GetAll(ctx context.Context, filter filter.HotelsFilter) (models.HotelsPagination, error) {

	data, err := l.repo.GetAll(ctx, filter)

	if err != nil {
		return models.HotelsPagination{}, err
	}

	return data, nil
}

func (l *hotelssvcs) Add(ctx context.Context, data *models.HotelsDto) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
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
		Id:        primitive.NewObjectID(),
		Trash:     false,
		CreatedAt: time.Now(),
		CreatedBy: cfg.User.Id,
		UpdatedAt: time.Now(),
		UpdatedBy: cfg.User.Id,
	}

	if err := l.repo.Add(ctx, hotels); err != nil {
		return nil, err
	}

	return hotels, nil

}

func (a *hotelssvcs) Update(ctx context.Context, id string, data *models.HotelsDto) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InvalidObjectId()
	}
	return a.repo.Update(ctx, _id, data)
}

func (a *hotelssvcs) Delete(ctx context.Context, id string) error {

	return a.repo.Delete(ctx, id)
}
