package statistics

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/picklist"

	"go.mongodb.org/mongo-driver/bson"
)

type StatisticsResponse struct {
	ProgramsCount      int64 `json:"programsCount"`
	TripsCount         int64 `json:"tripsCount"`
	CountriesCount     int64 `json:"countriesCount"`
	HotelsCount        int64 `json:"hotelsCount"`
	HappyTravelerCount int64 `json:"happyTravelerCount"`
}

// StatisticsSvcs defines available statistics methods
type StatisticsSvcs interface {
	GetStatistics(ctx context.Context) (*StatisticsResponse, error)
	GetProgramsCount(ctx context.Context) (int64, error)
	GetTripsCount(ctx context.Context) (int64, error)
	GetCountriesCount(ctx context.Context) (int64, error)
	GetHotelsCount(ctx context.Context) (int64, error)
	GetHappyTravelersCount(ctx context.Context) (int64, error)
}

type statisticsSvcs struct {
	programSvcs       ourservice.ProgramSvcs
	travelRequestSvcs ourservice.TravelRequestSvcs
	hotelsSvcs        ourservice.HotelsSvcs
	packageSvcs       ourservice.PackageSvcs
	destSvcs          picklist.DestinationSvcs
	reviewSvcs        interactions.ReviewsSvcs
}

func NewStatisticsSvcs(
	programSvcs ourservice.ProgramSvcs,
	travelRequestSvcs ourservice.TravelRequestSvcs,
	hotelsSvcs ourservice.HotelsSvcs,
	packageSvcs ourservice.PackageSvcs,
	destSvcs picklist.DestinationSvcs,
	reviewSvcs interactions.ReviewsSvcs,
) StatisticsSvcs {
	return &statisticsSvcs{
		programSvcs:       programSvcs,
		travelRequestSvcs: travelRequestSvcs,
		hotelsSvcs:        hotelsSvcs,
		packageSvcs:       packageSvcs,
		destSvcs:          destSvcs,
		reviewSvcs:        reviewSvcs,
	}
}

// GetStatistics concurrently fetches all metrics
func (s *statisticsSvcs) GetStatistics(ctx context.Context) (*StatisticsResponse, error) {
	programsCh := make(chan int64)
	tripsCh := make(chan int64)
	countriesCh := make(chan int64)
	hotelsCh := make(chan int64)
	travelersCh := make(chan int64)
	errCh := make(chan error, 5)

	go func() {
		count, err := s.GetProgramsCount(ctx)
		programsCh <- count
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		count, err := s.GetTripsCount(ctx)
		tripsCh <- count
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		count, err := s.GetCountriesCount(ctx)
		countriesCh <- count
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		count, err := s.GetHotelsCount(ctx)
		hotelsCh <- count
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		count, err := s.GetHappyTravelersCount(ctx)
		travelersCh <- count
		if err != nil {
			errCh <- err
		}
	}()

	programsCount := <-programsCh
	tripsCount := <-tripsCh
	countriesCount := <-countriesCh
	hotelsCount := <-hotelsCh
	travelersCount := <-travelersCh

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return &StatisticsResponse{
		ProgramsCount:      programsCount,
		TripsCount:         tripsCount,
		CountriesCount:     countriesCount,
		HotelsCount:        hotelsCount,
		HappyTravelerCount: travelersCount,
	}, nil
}

func (s *statisticsSvcs) GetProgramsCount(ctx context.Context) (int64, error) {
	filter := bson.M{"trash": bson.M{"$ne": true}}
	return s.programSvcs.Count(ctx, filter)
}

func (s *statisticsSvcs) GetTripsCount(ctx context.Context) (int64, error) {
	approved := "approved"
	f := &filter.TravelReqFilters{Status: &approved}
	pipeline := f.BuildPipeline(bson.M{})
	return s.travelRequestSvcs.Count(ctx, pipeline)
}

func (s *statisticsSvcs) GetCountriesCount(ctx context.Context) (int64, error) {
	filter := bson.M{"trash": bson.M{"$ne": true}}
	return s.destSvcs.Count(ctx, filter)
}

func (s *statisticsSvcs) GetHotelsCount(ctx context.Context) (int64, error) {

	filter := bson.M{"trash": false}
	return s.hotelsSvcs.Count(ctx, filter)
}

func (s *statisticsSvcs) GetHappyTravelersCount(ctx context.Context) (int64, error) {

	reviews, err := s.reviewSvcs.GetAll(ctx)
	if err != nil {
		return 0, err
	}

	uniqueUsers := make(map[string]bool)
	for _, review := range reviews {

		if !review.Trash && review.Value >= 3 && review.Status == "approved" {
			if review.UserId == "000000000000000000000000" {

				if review.Email != "" {
					uniqueUsers[review.Email] = true
				}
			} else {

				if review.UserId != "" {
					uniqueUsers[review.UserId] = true
				}
			}
		}
	}

	return int64(len(uniqueUsers)), nil
}
