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
	"strconv"
	"time"

	"larsa-tourism-microservices/pkg/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
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
	GetBookingsAdmin2(ctx context.Context, email *string, fromDate, toDate *time.Time) ([]models.Booking, error)
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

// convertBookingsArray converts raw booking data and handles type conversions (e.g., string to int for guestId)
func convertBookingsArray(rawBookings []interface{}) []models.Booking {
	var bookings []models.Booking
	for _, rawBooking := range rawBookings {
		bookingMap, ok := rawBooking.(map[string]interface{})
		if !ok {
			continue
		}

		// Convert guestId from string to int if needed
		if guestIdVal, exists := bookingMap["guestId"]; exists {
			if guestIdStr, ok := guestIdVal.(string); ok {
				if guestIdInt, err := strconv.Atoi(guestIdStr); err == nil {
					bookingMap["guestId"] = guestIdInt
				}
			}
		}

		// Convert supplierId from string to int if needed
		if supplierIdVal, exists := bookingMap["supplierId"]; exists {
			if supplierIdStr, ok := supplierIdVal.(string); ok {
				if supplierIdInt, err := strconv.Atoi(supplierIdStr); err == nil {
					bookingMap["supplierId"] = supplierIdInt
				}
			}
		}

		// Convert userId from string to int if needed
		if userIdVal, exists := bookingMap["userId"]; exists {
			if userIdStr, ok := userIdVal.(string); ok {
				if userIdInt, err := strconv.Atoi(userIdStr); err == nil {
					bookingMap["userId"] = userIdInt
				}
			}
		}

		// Marshal back to JSON and unmarshal into Booking struct
		bookingJSON, err := json.Marshal(bookingMap)
		if err != nil {
			continue
		}

		var booking models.Booking
		if err := json.Unmarshal(bookingJSON, &booking); err != nil {
			continue
		}

		bookings = append(bookings, booking)
	}
	return bookings
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

// GetBookingsByEmail gets bookings for authenticated user by email with date range from local database
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

	// Query local database for bookings by email
	filter := bson.M{
		"email": email,
	}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var userBookings []models.UserBooking
	err = r.bookingRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &userBookings)
	})
	if err != nil {
		return nil, err
	}

	// Convert UserBooking to Booking and filter by date
	var bookings []models.Booking
	for _, userBooking := range userBookings {
		// Parse booking createdAt date from the Booking's CreatedAt field (string)
		var bookingDate time.Time

		if userBooking.Booking.CreatedAt != "" {
			var err error
			bookingDate, err = time.Parse(time.RFC3339, userBooking.Booking.CreatedAt)
			if err != nil {
				// Try alternative format
				bookingDate, err = time.Parse("2006-01-02T15:04:05Z07:00", userBooking.Booking.CreatedAt)
				if err != nil {
					// Try simple date format
					bookingDate, err = time.Parse("2006-01-02", userBooking.Booking.CreatedAt)
					if err != nil {
						// Fallback to UserBooking's CreatedAt (time.Time)
						bookingDate = userBooking.CreatedAt
					}
				}
			}
		} else {
			// Fallback to UserBooking's CreatedAt if Booking.CreatedAt is empty
			bookingDate = userBooking.CreatedAt
		}

		// Check if booking date is within range
		if (bookingDate.After(*fromDate) || bookingDate.Equal(*fromDate)) &&
			(bookingDate.Before(*toDate) || bookingDate.Equal(*toDate)) {
			bookings = append(bookings, userBooking.Booking)
		}
	}

	return bookings, nil
}

// GetBookingsAdmin2 gets bookings from LiteAPI by date range and filters by email in backend
// Default date range: year start to now
// Note: Admin authorization should be handled by middleware
func (r *ratessvcs) GetBookingsAdmin2(ctx context.Context, email *string, fromDate, toDate *time.Time) ([]models.Booking, error) {
	// Set default date range: year start to now
	now := time.Now()
	if fromDate == nil {
		yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		fromDate = &yearStart
	}
	if toDate == nil {
		toDate = &now
	}

	// Format dates for LiteAPI (YYYY-MM-DD format)
	startDateStr := fromDate.Format("2006-01-02")
	endDateStr := toDate.Format("2006-01-02")

	// Fetch from LiteAPI with date range
	liteApiSdk, err := r.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetBookingsByDateRange(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	// Parse response - LiteAPI returns: {"count": 25, "data": [...]}
	// First unmarshal into flexible structure to handle type conversions
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(resp.Data, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to parse bookings response: %v", err)
	}

	// Extract the data array from the response
	var bookingsArray []interface{}
	if dataField, exists := rawResponse["data"]; exists {
		if dataField == nil {
			// If data is null, return empty array
			return []models.Booking{}, nil
		}
		if arr, ok := dataField.([]interface{}); ok {
			bookingsArray = arr
		} else {
			// If data is not an array, return empty
			return []models.Booking{}, nil
		}
	} else {
		// If no "data" field, return empty
		return []models.Booking{}, nil
	}

	// Convert bookings array with type conversions
	convertedBookings := convertBookingsArray(bookingsArray)

	// Filter by email if provided
	var filteredBookings []models.Booking
	if email != nil && *email != "" {
		for _, booking := range convertedBookings {
			if booking.Email == *email {
				filteredBookings = append(filteredBookings, booking)
			}
		}
	} else {
		// Return all bookings if no email filter
		filteredBookings = convertedBookings
	}

	return filteredBookings, nil
}

// GetBookingsAdmin gets all bookings for admin with optional email filter from local database
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

	// Build filter for local database query
	filter := bson.M{}
	if email != nil {
		filter["email"] = *email
	}

	pipeline := []bson.M{
		{"$match": filter},
	}

	var userBookings []models.UserBooking
	err := r.bookingRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &userBookings)
	})
	if err != nil {
		return nil, err
	}

	// Convert UserBooking to Booking and filter by date
	var bookings []models.Booking
	for _, userBooking := range userBookings {
		// Parse booking createdAt date from the Booking's CreatedAt field (string)
		var bookingDate time.Time

		if userBooking.Booking.CreatedAt != "" {
			var err error
			bookingDate, err = time.Parse(time.RFC3339, userBooking.Booking.CreatedAt)
			if err != nil {
				// Try alternative format
				bookingDate, err = time.Parse("2006-01-02T15:04:05Z07:00", userBooking.Booking.CreatedAt)
				if err != nil {
					// Try simple date format
					bookingDate, err = time.Parse("2006-01-02", userBooking.Booking.CreatedAt)
					if err != nil {
						// Fallback to UserBooking's CreatedAt (time.Time)
						bookingDate = userBooking.CreatedAt
					}
				}
			}
		} else {
			// Fallback to UserBooking's CreatedAt if Booking.CreatedAt is empty
			bookingDate = userBooking.CreatedAt
		}

		// Check if booking date is within range
		if (bookingDate.After(*fromDate) || bookingDate.Equal(*fromDate)) &&
			(bookingDate.Before(*toDate) || bookingDate.Equal(*toDate)) {
			bookings = append(bookings, userBooking.Booking)
		}
	}

	return bookings, nil
}
