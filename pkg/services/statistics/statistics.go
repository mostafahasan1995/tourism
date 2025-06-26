package statistics

import (
	"context"
	"larsa-tourism-microservices/pkg/services/interactions"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/statistics/models"
	"larsa-tourism-microservices/pkg/services/statistics/repository"

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
	// Manual statistics methods
	GetManualStatistics(ctx context.Context) (*models.ManualStatistics, error)
	UpdateManualStatistics(ctx context.Context, req *models.ManualStatisticsRequest) (*models.ManualStatistics, error)
	ToggleAutoCalculate(ctx context.Context, autoCalculate bool) (*models.ManualStatistics, error)
}

type statisticsSvcs struct {
	programSvcs       ourservice.ProgramSvcs
	travelRequestSvcs ourservice.TravelRequestSvcs
	hotelsSvcs        ourservice.HotelsSvcs
	packageSvcs       ourservice.PackageSvcs
	destSvcs          picklist.DestinationSvcs
	reviewSvcs        interactions.ReviewsSvcs
	manualStatsRepo   repository.ManualStatisticsRepository
}

func NewStatisticsSvcs(
	programSvcs ourservice.ProgramSvcs,
	travelRequestSvcs ourservice.TravelRequestSvcs,
	hotelsSvcs ourservice.HotelsSvcs,
	packageSvcs ourservice.PackageSvcs,
	destSvcs picklist.DestinationSvcs,
	reviewSvcs interactions.ReviewsSvcs,
	manualStatsRepo repository.ManualStatisticsRepository,
) StatisticsSvcs {
	return &statisticsSvcs{
		programSvcs:       programSvcs,
		travelRequestSvcs: travelRequestSvcs,
		hotelsSvcs:        hotelsSvcs,
		packageSvcs:       packageSvcs,
		destSvcs:          destSvcs,
		reviewSvcs:        reviewSvcs,
		manualStatsRepo:   manualStatsRepo,
	}
}

// GetStatistics concurrently fetches all metrics
func (s *statisticsSvcs) GetStatistics(ctx context.Context) (*StatisticsResponse, error) {
	// Check if auto-calculate is enabled
	manualStats, err := s.manualStatsRepo.GetManualStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// If auto-calculate is disabled, return manual statistics
	if !manualStats.AutoCalculate {
		return &StatisticsResponse{
			ProgramsCount:      manualStats.ProgramsCount,
			TripsCount:         manualStats.TripsCount,
			CountriesCount:     manualStats.CountriesCount,
			HotelsCount:        manualStats.HotelsCount,
			HappyTravelerCount: manualStats.HappyTravelerCount,
		}, nil
	}

	// Auto-calculate is enabled, fetch real data
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

	reviews, err := s.reviewSvcs.GetAllApproved(ctx)
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

// GetManualStatistics retrieves the current manual statistics configuration
func (s *statisticsSvcs) GetManualStatistics(ctx context.Context) (*models.ManualStatistics, error) {
	return s.manualStatsRepo.GetManualStatistics(ctx)
}

// UpdateManualStatistics updates manual statistics values
func (s *statisticsSvcs) UpdateManualStatistics(ctx context.Context, req *models.ManualStatisticsRequest) (*models.ManualStatistics, error) {
	// Get current manual statistics
	currentStats, err := s.manualStatsRepo.GetManualStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.ProgramsCount != nil {
		currentStats.ProgramsCount = *req.ProgramsCount
	}
	if req.TripsCount != nil {
		currentStats.TripsCount = *req.TripsCount
	}
	if req.CountriesCount != nil {
		currentStats.CountriesCount = *req.CountriesCount
	}
	if req.HotelsCount != nil {
		currentStats.HotelsCount = *req.HotelsCount
	}
	if req.HappyTravelerCount != nil {
		currentStats.HappyTravelerCount = *req.HappyTravelerCount
	}
	if req.AutoCalculate != nil {
		currentStats.AutoCalculate = *req.AutoCalculate
	}

	// Update in database
	err = s.manualStatsRepo.UpdateManualStatistics(ctx, currentStats)
	if err != nil {
		return nil, err
	}

	return currentStats, nil
}

// ToggleAutoCalculate toggles the auto-calculate mode
func (s *statisticsSvcs) ToggleAutoCalculate(ctx context.Context, autoCalculate bool) (*models.ManualStatistics, error) {
	err := s.manualStatsRepo.ToggleAutoCalculate(ctx, autoCalculate)
	if err != nil {
		return nil, err
	}

	return s.manualStatsRepo.GetManualStatistics(ctx)
}
