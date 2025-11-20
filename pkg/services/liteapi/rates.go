package liteapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"larsa-tourism-microservices/pkg/services/member"
	"reflect"
	"time"

	"larsa-tourism-microservices/pkg/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatesSvcs interface {
	GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error)
	GetMinRates(ctx context.Context, data any) (any, error)
	PreBook(ctx context.Context, data map[string]any) (any, error)
	Book(ctx context.Context, data map[string]any) (any, error)
	MyPrebooks(ctx context.Context) ([]models.UserPrebook, error)
	MyBookings(ctx context.Context) ([]models.UserBooking, error)
	CancelBooking(ctx context.Context, bookingId string) (any, error)
	//
	GetFullRatesStream(ctx context.Context, data any) (any, error)
	// New booking list methods
	GetBookingsByEmail(ctx context.Context, fromDate, toDate *time.Time) ([]models.Booking, error)
	GetBookingsAdmin(ctx context.Context, email *string, fromDate, toDate *time.Time) ([]models.Booking, error)
}

type ratessvcs struct {
	liteApiInitFunc       liteApiSdk.LiteApiInitFunc
	preBookRepo           repo.PreBookRepo
	bookingRepo           repo.BookingRepo
	cachedBookingListRepo repo.CachedBookingListRepo
	customersvcs          member.CustomerSvcs
}

func NewRatesSvcs(i *do.Injector) (RatesSvcs, error) {
	return &ratessvcs{
		liteApiInitFunc:       do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		preBookRepo:           do.MustInvoke[repo.PreBookRepo](i),
		bookingRepo:           do.MustInvoke[repo.BookingRepo](i),
		cachedBookingListRepo: do.MustInvoke[repo.CachedBookingListRepo](i),
		customersvcs:          do.MustInvoke[member.CustomerSvcs](i),
	}, nil
}

func (r *ratessvcs) GetHotelsRates(ctx context.Context, data any) (*models.RatesList, error) {
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetFullRates(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.RatesList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ratessvcs) GetMinRates(ctx context.Context, data any) (any, error) {
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetMinRates(data)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ratessvcs) PreBook(ctx context.Context, data map[string]any) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.customersvcs.GetOne(ctx, cfg.User.Id.Hex())
	if err != nil {
		return nil, errors.New("error get customer")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.PreBook(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.PreBookData
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	userPrebook := &models.UserPrebook{
		PreBook:    result.Data,
		GusetLevel: result.GuestLevel,
		UserId:     cfg.User.Id,
		Status:     "pending",
		CreatedAt:  time.Now(),
		Trash:      false,
	}

	go util.WithRetry(func() error {
		return r.preBookRepo.Add(context.WithoutCancel(ctx), userPrebook)
	}, 3)

	//in case the liteApi request success we must send the response
	return result, nil
}

func (r *ratessvcs) Book(ctx context.Context, data map[string]any) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.customersvcs.GetOne(ctx, cfg.User.Id.Hex())
	if err != nil {
		return nil, errors.New("error get customer")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.Book(data)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.BookingData
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	userBooking := &models.UserBooking{
		Booking:    result.Data,
		GusetLevel: result.GuestLevel,
		UserId:     cfg.User.Id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	go util.WithRetry(func() error {
		return r.bookingRepo.Add(context.WithoutCancel(ctx), userBooking)
	}, 3)

	return result, nil
}

func (r *ratessvcs) MyPrebooks(ctx context.Context) ([]models.UserPrebook, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": cfg.User.Id, "trash": false}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var result []models.UserPrebook
	err = r.preBookRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ratessvcs) MyBookings(ctx context.Context) ([]models.UserBooking, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": cfg.User.Id}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var result []models.UserBooking
	err = r.bookingRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (r *ratessvcs) CancelBooking(ctx context.Context, bookingId string) (any, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.bookingRepo.GetByFilter(ctx, bson.M{"userId": cfg.User.Id, "bookingId": bookingId})
	if err != nil {
		return nil, errors.New("booking not found or not belongs to the user")
	}

	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.CancelBooking(bookingId)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go util.WithRetry(func() error {
		filter := bson.M{"userId": cfg.User.Id, "bookingId": bookingId}
		update := bson.M{"$set": bson.M{"status": "CANCELLED", "updatedAt": time.Now()}}
		_, err := r.bookingRepo.Patch(context.WithoutCancel(ctx), filter, update)
		return err
	}, 3)

	return result, nil
}

func (r *ratessvcs) GetFullRatesStream(ctx context.Context, data any) (any, error) {
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := liteApiSdk.GetFullRatesStream(data)
	if err != nil {
		return nil, err
	}

	var allRates []map[string]any

	for event := range ch {
		if event.Error != nil {
			return nil, event.Error
		}
		if event.Done {
			break
		}
		fmt.Println(string(event.Data))

		var result map[string]any
		if err := json.Unmarshal(event.Data, &result); err != nil {
			return nil, err
		}

		allRates = append(allRates, result)

	}

	fmt.Println(allRates[0])

	return nil, nil

}

// getEmailFromUser extracts email from user context using reflection
func getEmailFromUser(ctx context.Context) (string, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return "", err
	}

	if cfg.User == nil {
		return "", errors.New("user not authenticated")
	}

	// Use reflection to get Email field from common.User
	userDataValue := reflect.ValueOf(cfg.User.UserData)
	if userDataValue.Kind() == reflect.Ptr {
		userDataValue = userDataValue.Elem()
	}

	if userDataValue.Kind() == reflect.Struct {
		emailField := userDataValue.FieldByName("Email")
		if emailField.IsValid() && emailField.Kind() == reflect.String {
			return emailField.String(), nil
		}
	}

	return "", errors.New("email not found in user data")
}

// filterBookingsByDate filters bookings by date range
func filterBookingsByDate(bookings []models.Booking, fromDate, toDate time.Time) []models.Booking {
	var filtered []models.Booking
	for _, booking := range bookings {
		// Parse booking createdAt date
		bookingDate, err := time.Parse(time.RFC3339, booking.CreatedAt)
		if err != nil {
			// Try alternative format
			bookingDate, err = time.Parse("2006-01-02T15:04:05Z07:00", booking.CreatedAt)
			if err != nil {
				continue // Skip if date parsing fails
			}
		}

		// Check if booking date is within range
		if (bookingDate.After(fromDate) || bookingDate.Equal(fromDate)) &&
			(bookingDate.Before(toDate) || bookingDate.Equal(toDate)) {
			filtered = append(filtered, booking)
		}
	}
	return filtered
}

// GetBookingsByEmail gets bookings for authenticated user by email with date range
// Default date range: last year to now
func (r *ratessvcs) GetBookingsByEmail(ctx context.Context, fromDate, toDate *time.Time) ([]models.Booking, error) {
	email, err := getEmailFromUser(ctx)
	if err != nil {
		return nil, err
	}

	// Set default date range: last year to now
	now := time.Now()
	if fromDate == nil {
		oneYearAgo := now.AddDate(-1, 0, 0)
		fromDate = &oneYearAgo
	}
	if toDate == nil {
		toDate = &now
	}

	// Check cache first
	filter := bson.M{
		"email":     email,
		"fromDate":  *fromDate,
		"toDate":    *toDate,
		"isAdmin":   false,
		"expiresAt": bson.M{"$gt": time.Now()},
	}

	cached, err := r.cachedBookingListRepo.GetByFilter(ctx, filter)
	if err == nil && cached != nil {
		// Return cached bookings
		return cached.Bookings, nil
	}

	// Fetch from LiteAPI
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetBookingsByEmail(email)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	// Parse response - LiteAPI returns bookings array
	var bookingsData struct {
		Data []models.Booking `json:"data"`
	}
	if err := json.Unmarshal(resp.Data, &bookingsData); err != nil {
		// Try alternative format - might be direct array
		var directBookings []models.Booking
		if err2 := json.Unmarshal(resp.Data, &directBookings); err2 != nil {
			return nil, fmt.Errorf("failed to parse bookings response: %v", err)
		}
		bookingsData.Data = directBookings
	}

	// Filter by date range
	filteredBookings := filterBookingsByDate(bookingsData.Data, *fromDate, *toDate)

	// Cache the results
	cachedList := &models.CachedBookingList{
		Id:        primitive.NewObjectID(),
		Email:     email,
		FromDate:  *fromDate,
		ToDate:    *toDate,
		Bookings:  filteredBookings,
		IsAdmin:   false,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // Cache for 1 hour
	}

	go util.WithRetry(func() error {
		return r.cachedBookingListRepo.Add(context.WithoutCancel(ctx), cachedList)
	}, 3)

	return filteredBookings, nil
}

// GetBookingsAdmin gets all bookings for admin with optional email filter
// Default date range: year start to now
// Note: Admin authorization should be handled by middleware
func (r *ratessvcs) GetBookingsAdmin(ctx context.Context, email *string, fromDate, toDate *time.Time) ([]models.Booking, error) {
	// Set default date range: year start to now
	now := time.Now()
	if fromDate == nil {
		yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		fromDate = &yearStart
	}
	if toDate == nil {
		toDate = &now
	}

	// Check cache
	filter := bson.M{
		"fromDate":  *fromDate,
		"toDate":    *toDate,
		"isAdmin":   true,
		"expiresAt": bson.M{"$gt": time.Now()},
	}
	if email != nil {
		filter["email"] = *email
	}

	cached, err := r.cachedBookingListRepo.GetByFilter(ctx, filter)
	if err == nil && cached != nil {
		return cached.Bookings, nil
	}

	// For admin, we need to fetch all bookings
	// Since LiteAPI doesn't have a direct "get all bookings" endpoint,
	// we'll use GetBookingsList with empty or use a different approach
	// For now, if email is provided, use it; otherwise return error (admin needs to specify email or we need to implement aggregation)
	if email == nil {
		return nil, errors.New("email filter is required for admin bookings query")
	}

	// Fetch from LiteAPI
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetBookingsByEmail(*email)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	// Parse response
	var bookingsData struct {
		Data []models.Booking `json:"data"`
	}
	if err := json.Unmarshal(resp.Data, &bookingsData); err != nil {
		var directBookings []models.Booking
		if err2 := json.Unmarshal(resp.Data, &directBookings); err2 != nil {
			return nil, fmt.Errorf("failed to parse bookings response: %v", err)
		}
		bookingsData.Data = directBookings
	}

	// Filter by date range
	filteredBookings := filterBookingsByDate(bookingsData.Data, *fromDate, *toDate)

	// Cache the results
	cachedList := &models.CachedBookingList{
		Id:        primitive.NewObjectID(),
		Email:     *email,
		FromDate:  *fromDate,
		ToDate:    *toDate,
		Bookings:  filteredBookings,
		IsAdmin:   true,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // Cache for 1 hour
	}

	go util.WithRetry(func() error {
		return r.cachedBookingListRepo.Add(context.WithoutCancel(ctx), cachedList)
	}, 3)

	return filteredBookings, nil
}
