package statistics

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/picklist"
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
	result, err := s.programSvcs.GetAll(ctx, "")
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}

func (s *statisticsSvcs) GetTripsCount(ctx context.Context) (int64, error) {
	approved := "approved"
	f := &filter.TravelReqFilters{Status: &approved}
	result, err := s.travelRequestSvcs.GetAll(ctx, f)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}

func (s *statisticsSvcs) GetCountriesCount(ctx context.Context) (int64, error) {
	destinations, err := s.destSvcs.GetAll(ctx, "")
	if err != nil {
		return 0, err
	}
	return int64(len(destinations)), nil
}

func (s *statisticsSvcs) GetHotelsCount(ctx context.Context) (int64, error) {
	emptyFilter := filter.HotelsFilter{}
	result, err := s.hotelsSvcs.GetAll(ctx, emptyFilter, 1, 1000)
	if err != nil {
		return 0, err
	}
	return int64(result.Pagination.TotalCount), nil
}

func (s *statisticsSvcs) GetHappyTravelersCount(ctx context.Context) (int64, error) {
	reviews, err := s.reviewSvcs.GetAll(ctx)
	if err != nil {
		return 0, err
	}
	seen := make(map[string]struct{})
	for _, r := range reviews {
		if r.Value >= 3 {
			seen[r.UserId] = struct{}{}
		}
	}
	return int64(len(seen)), nil
}
