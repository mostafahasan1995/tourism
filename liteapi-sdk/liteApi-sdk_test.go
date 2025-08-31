package liteApiSdk

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// Test configuration
const testAPIKey = "sand_c0155ab8-c683-4f26-8f94-b5e92c5797b9"

var (
	sdk             *LiteApiSdk
	latestVoucherId string
	bookingId       string
	offer           string
)

func init() {
	sdk = NewLiteApiSdk(testAPIKey)
	rand.Seed(time.Now().UnixNano())
}

// Helper function to generate voucher codes
func getVoucherCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 9)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return "VOUCHER_" + string(code)
}

func TestGetFullRates(t *testing.T) {
	data := map[string]interface{}{
		"hotelIds":         []string{"lp19d4c"},
		"occupancies":      []map[string]interface{}{{"adults": 1, "children": []int{}}},
		"currency":         "USD",
		"guestNationality": "US",
		"checkin":          "2025-11-19",
		"checkout":         "2025-11-20",
		"countryCode":      "USD",
	}

	result := sdk.GetFullRates(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Extract offer if available
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if dataArray, ok := dataMap["data"].([]interface{}); ok && len(dataArray) > 0 {
			if hotel, ok := dataArray[0].(map[string]interface{}); ok {
				if roomTypes, ok := hotel["roomTypes"].([]interface{}); ok && len(roomTypes) > 0 {
					if roomType, ok := roomTypes[0].(map[string]interface{}); ok {
						if offerId, ok := roomType["offerId"].(string); ok {
							offer = offerId
						}
					}
				}
			}
		}
	}
}

func TestPreBook(t *testing.T) {
	if offer == "" {
		t.Skip("Skipping prebook test - no offer available from getFullRates")
	}

	data := map[string]interface{}{
		"offerId": offer,
	}

	result := sdk.PreBook(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestBook(t *testing.T) {
	data := map[string]interface{}{
		"holder": map[string]interface{}{
			"firstName": "Steve",
			"lastName":  "Doe",
			"email":     "s.doe@liteapi.travel",
		},
		"payment": map[string]interface{}{
			"method": "ACC_CREDIT_CARD",
		},
		"prebookId": "6-xUGK8_C",
		"guests": []map[string]interface{}{
			{
				"occupancyNumber": 1,
				"remarks":         "quiet room please",
				"firstName":       "Sunny",
				"lastName":        "Mars",
				"email":           "s.mars@liteapi.travel",
			},
		},
	}

	result := sdk.Book(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Extract booking ID if available
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if bId, ok := dataMap["bookingId"].(string); ok {
			bookingId = bId
		}
	}
}

func TestGetBookingsList(t *testing.T) {
	result := sdk.GetBookingsList("testref")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestRetrieveBooking(t *testing.T) {
	result := sdk.RetrieveBooking("XE1Bxh1bS")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestCancelBooking(t *testing.T) {
	if bookingId == "" {
		t.Skip("Skipping cancel booking test - no booking ID available")
	}

	result := sdk.CancelBooking(bookingId)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetCitiesByCountryCode(t *testing.T) {
	result := sdk.GetCitiesByCountryCode("SG")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetCurrencies(t *testing.T) {
	result := sdk.GetCurrencies()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetHotelFacilities(t *testing.T) {
	result := sdk.GetHotelFacilities()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetHotelTypes(t *testing.T) {
	result := sdk.GetHotelTypes()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetHotelChains(t *testing.T) {
	result := sdk.GetHotelChains()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetPlaces(t *testing.T) {
	result := sdk.GetPlaces("Rome", "", "en")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetHotels(t *testing.T) {
	parameters := map[string]string{
		"countryCode": "IT",
		"cityName":    "Rome",
	}

	result := sdk.GetHotels(parameters, "en")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetHotelDetails(t *testing.T) {
	result := sdk.GetHotelDetails("lp1897", "fr")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetDataReviews(t *testing.T) {
	result := sdk.GetDataReviews("lp1897", 5, true)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Check for sentiment analysis in the response
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if _, hasSentiment := dataMap["sentimentAnalysis"]; !hasSentiment {
			t.Error("Expected sentimentAnalysis to be present in response")
		}
	}
}

func TestGetCountries(t *testing.T) {
	result := sdk.GetCountries()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetIataCodes(t *testing.T) {
	result := sdk.GetIataCodes()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetGuestsIds(t *testing.T) {
	result := sdk.GetGuestsIds("10")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetGuestsBookings(t *testing.T) {
	result := sdk.GetGuestsBookings("10")

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetVouchers(t *testing.T) {
	result := sdk.GetVouchers()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Extract latest voucher ID
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if vouchers, ok := dataMap["vouchers"].([]interface{}); ok && len(vouchers) > 0 {
			var latestVoucher map[string]interface{}
			var latestTime time.Time

			for _, v := range vouchers {
				if voucher, ok := v.(map[string]interface{}); ok {
					if createdAt, ok := voucher["created_at"].(string); ok {
						if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
							if t.After(latestTime) {
								latestTime = t
								latestVoucher = voucher
							}
						}
					}
				}
			}

			if latestVoucher != nil {
				if id, ok := latestVoucher["id"].(string); ok {
					latestVoucherId = id
				}
			}
		}
	}
}

func TestGetVoucherById(t *testing.T) {
	if latestVoucherId == "" {
		t.Skip("Skipping voucher by ID test - no voucher ID available")
	}

	result := sdk.GetVoucherById(latestVoucherId)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestCreateVoucher(t *testing.T) {
	data := map[string]interface{}{
		"voucher_code":            getVoucherCode(),
		"discount_type":           "percentage",
		"discount_value":          12,
		"minimum_spend":           60,
		"maximum_discount_amount": 20,
		"currency":                "USD",
		"validity_start":          "2024-06-03",
		"validity_end":            "2024-07-30",
		"usages_limit":            10,
		"status":                  "active",
	}

	result := sdk.CreateVoucher(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestUpdateVoucher(t *testing.T) {
	if latestVoucherId == "" {
		t.Skip("Skipping update voucher test - no voucher ID available")
	}

	data := map[string]interface{}{
		"voucher_code":            getVoucherCode(),
		"discount_type":           "percentage",
		"discount_value":          12,
		"minimum_spend":           60,
		"maximum_discount_amount": 20,
		"currency":                "USD",
		"validity_start":          "2024-06-03",
		"validity_end":            "2024-07-30",
		"usages_limit":            10,
		"status":                  "active",
	}

	result := sdk.UpdateVoucher(latestVoucherId, data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestUpdateVoucherStatus(t *testing.T) {
	if latestVoucherId == "" {
		t.Skip("Skipping update voucher status test - no voucher ID available")
	}

	data := map[string]interface{}{
		"status": "inactive",
	}

	result := sdk.UpdateVoucherStatus(latestVoucherId, data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetLoyalty(t *testing.T) {
	result := sdk.GetLoyalty()

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestEnableLoyalty(t *testing.T) {
	data := map[string]interface{}{
		"status":       "enabled",
		"cashbackRate": 0.03,
	}

	result := sdk.EnableLoyalty(data)

	// Handle case where loyalty is already created
	if result.Status == "failed" {
		if errorMap, ok := result.Error.(map[string]interface{}); ok {
			if message, ok := errorMap["message"].(string); ok {
				if strings.Contains(message, "loyalty already created") {
					t.Logf("Loyalty already created - expected behavior")
					return
				}
			}
		}
		t.Errorf("Expected success or 'loyalty already created' error, got: %v", result.Error)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestUpdateLoyalty(t *testing.T) {
	data := map[string]interface{}{
		"status":       "enable",
		"cashbackRate": 0.03,
	}

	result := sdk.UpdateLoyalty(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestRetrieveWeeklyAnalytics(t *testing.T) {
	data := map[string]interface{}{
		"from": "2024-01-01",
		"to":   "2024-01-07",
	}

	result := sdk.RetrieveWeeklyAnalytics(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Check for arr field
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if _, hasArr := dataMap["arr"]; !hasArr {
			t.Error("Expected 'arr' field to be present in response")
		}
	}
}

func TestRetrieveAnalyticsReport(t *testing.T) {
	data := map[string]interface{}{
		"from": "2024-01-01",
		"to":   "2024-01-07",
	}

	result := sdk.RetrieveAnalyticsReport(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}

	// Check for required fields
	if dataMap, ok := result.Data.(map[string]interface{}); ok {
		if _, hasTotalRevenue := dataMap["totalRevenue"]; !hasTotalRevenue {
			t.Error("Expected 'totalRevenue' field to be present in response")
		}
		if _, hasSalesRevenue := dataMap["salesRevenue"]; !hasSalesRevenue {
			t.Error("Expected 'salesRevenue' field to be present in response")
		}
	}
}

func TestRetrieveMarketAnalytics(t *testing.T) {
	data := map[string]interface{}{
		"from": "2024-01-01",
		"to":   "2024-01-07",
	}

	result := sdk.RetrieveMarketAnalytics(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestRetrieveMostBookedHotels(t *testing.T) {
	data := map[string]interface{}{
		"from": "2024-01-01",
		"to":   "2024-01-07",
	}

	result := sdk.RetrieveMostBookedHotels(data)

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestGetMinRates(t *testing.T) {
	data := map[string]interface{}{
		"hotelIds":         []string{"lp1897"},
		"occupancies":      []map[string]interface{}{{"adults": 2, "children": []int{5}}},
		"currency":         "USD",
		"guestNationality": "US",
		"checkin":          "2025-12-30",
		"checkout":         "2025-12-31",
		"countryCode":      "USD",
	}

	result := sdk.GetMinRates(data)

	fmt.Printf("getMinRates endpoint response: %+v\n", result)
	if result.Data != nil {
		fmt.Printf("getMinRates endpoint response data: %+v\n", result.Data)
	}

	if result.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", result.Status)
	}

	if result.Data == nil {
		t.Error("Expected data to be present")
	}
}

// Validation tests
func TestPreBookValidation(t *testing.T) {
	// Test with missing offerId
	data := map[string]interface{}{}

	result := sdk.PreBook(data)

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestBookValidation(t *testing.T) {
	// Test with missing prebookId
	data := map[string]interface{}{}

	result := sdk.Book(data)

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestRetrieveBookingValidation(t *testing.T) {
	// Test with empty booking ID
	result := sdk.RetrieveBooking("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestCancelBookingValidation(t *testing.T) {
	// Test with empty booking ID
	result := sdk.CancelBooking("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetCitiesByCountryCodeValidation(t *testing.T) {
	// Test with empty country code
	result := sdk.GetCitiesByCountryCode("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetHotelDetailsValidation(t *testing.T) {
	// Test with empty hotel ID
	result := sdk.GetHotelDetails("", "en")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetDataReviewsValidation(t *testing.T) {
	// Test with empty hotel ID
	result := sdk.GetDataReviews("", 5, true)

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetGuestsIdsValidation(t *testing.T) {
	// Test with empty guest ID
	result := sdk.GetGuestsIds("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetGuestsBookingsValidation(t *testing.T) {
	// Test with empty guest ID
	result := sdk.GetGuestsBookings("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestGetVoucherByIdValidation(t *testing.T) {
	// Test with empty voucher ID
	result := sdk.GetVoucherById("")

	if result.Status != "failed" {
		t.Errorf("Expected status 'failed', got '%s'", result.Status)
	}

	if result.Errors == nil || len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}
