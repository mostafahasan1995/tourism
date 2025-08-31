package liteApiSdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// LiteApiSdk represents the main SDK client
type LiteApiSdk struct {
	ApiKey         string
	ServiceURL     string
	BookServiceURL string
	DashboardURL   string
	Client         *http.Client
}

// APIResponse represents a standard API response
type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

// NewLiteApiSdk creates a new instance of the LiteApi SDK
func NewLiteApiSdk(apiKey string) *LiteApiSdk {
	return &LiteApiSdk{
		ApiKey:         apiKey,
		ServiceURL:     "https://api.liteapi.travel/v3.0",
		BookServiceURL: "https://book.liteapi.travel/v3.0",
		DashboardURL:   "https://da.liteapi.travel",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest handles HTTP requests with common headers and error handling
func (sdk *LiteApiSdk) makeRequest(method, url string, body interface{}) *APIResponse {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("Failed to marshal request body: %v", err),
			}
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", sdk.ApiKey)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check status code first, before reading/parsing response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// For error responses, try to read the body for error details
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("HTTP %d: Failed to read error response", resp.StatusCode),
			}
		}

		var errorResult map[string]interface{}
		if err := json.Unmarshal(responseBody, &errorResult); err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody)),
			}
		}

		return &APIResponse{
			Status: "failed",
			Error:  errorResult["error"],
		}
	}

	// Only parse response body for successful requests
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	return &APIResponse{
		Status: "success",
		Data:   result,
	}
}

// makeRequestWithRetry handles requests with retry logic for rate limiting
func (sdk *LiteApiSdk) makeRequestWithRetry(method, url string, body interface{}, retries int, delay time.Duration) *APIResponse {
	resp := sdk.makeRequest(method, url, body)

	// Check for rate limit errors
	if resp.Status == "failed" && retries > 0 {
		if errorMap, ok := resp.Error.(map[string]interface{}); ok {
			if code, exists := errorMap["code"]; exists {
				if codeFloat, ok := code.(float64); ok && codeFloat == 4290 {
					time.Sleep(delay)
					return sdk.makeRequestWithRetry(method, url, body, retries-1, delay*2)
				}
			}
		}
	}

	return resp
}

// GetFullRates searches and returns all available rooms along with rates and cancellation policies
// The Full Rates API is to search and return all available rooms along with its rates, cancellation policies for a list of hotel ID's based on the search dates.
// For each hotel ID, all available room information is returned.
// The API also has a built in loyalty rewards system. The system rewards return users who have made prior bookings.
// If the search is coming from a known guest ID, the guest level is also returned along with the pricing that's appropriate for the guest level.
// If it is a new user, the guest ID will be generated at the time of the first confirmed booking.
func (sdk *LiteApiSdk) GetFullRates(data interface{}) *APIResponse {
	url := sdk.ServiceURL + "/hotels/rates"
	return sdk.makeRequest("POST", url, data)
}

// GetMinRates gets minimum rates for hotels
func (sdk *LiteApiSdk) GetMinRates(data interface{}) *APIResponse {
	url := sdk.ServiceURL + "/hotels/min-rates"
	return sdk.makeRequest("POST", url, data)
}

// PreBook confirms if the room and rates for the search criterion
// This API is used to confirm if the room and rates for the search criterion. The input to the endpoint is an array of rate Ids coming from the GET hotel full rates availability API.
// In response, the API generates a prebook Id, a new rate Id and contains information if price, cancellation policy or boarding information has changed.
func (sdk *LiteApiSdk) PreBook(data map[string]interface{}) *APIResponse {
	var errors []string

	// Validate offerId
	if offerId, exists := data["offerId"]; !exists || offerId == "" {
		errors = append(errors, "The offerId is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := sdk.BookServiceURL + "/rates/prebook"
	return sdk.makeRequest("POST", url, data)
}

// Book confirms a booking when the prebook Id and the rate Id from the pre book stage along with the guest and payment information are passed
// This API confirms a booking when the prebook Id and the rate Id from the pre book stage along with the guest and payment information are passed.
// The guest information is an object that should include the guest first name, last name and email.
// The payment information is an object that should include the name, credit card number, expiry and CVC number.
// The response will confirm the booking along with a booking Id and a hotel confirmation code. It will also include the booking details including the dates, price and the cancellation policies.
func (sdk *LiteApiSdk) Book(data map[string]interface{}) *APIResponse {
	var errors []string

	// Validate prebookId
	if prebookId, exists := data["prebookId"]; !exists || prebookId == "" {
		errors = append(errors, "The prebookId is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := sdk.BookServiceURL + "/rates/book"
	return sdk.makeRequest("POST", url, data)
}

// GetBookingsList returns the list of booking Id's for a given guest Id
// The API returns the list of booking Id's for a given guest Id.
func (sdk *LiteApiSdk) GetBookingsList(clientReference string) *APIResponse {
	url := fmt.Sprintf("%s/bookings?clientReference=%s", sdk.BookServiceURL, clientReference)
	return sdk.makeRequest("GET", url, nil)
}

// RetrieveBooking returns the status and the details for a specific booking Id
// The API returns the status and the details for the a specific booking Id.
func (sdk *LiteApiSdk) RetrieveBooking(bookingId string) *APIResponse {
	var errors []string
	if bookingId == "" {
		errors = append(errors, "The booking ID is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := fmt.Sprintf("%s/bookings/%s", sdk.BookServiceURL, bookingId)
	return sdk.makeRequest("GET", url, nil)
}

// CancelBooking requests a cancellation of an existing confirmed booking
// This API is used to request a cancellation of an existing confirmed booking. Cancellation policies and conditions will be used to determine the success of the cancellation. For example a booking with non-refundable (NRFN) tag or a booking with a cancellation policy that was requested past the cancellation date will not be able to cancel the confirmed booking.
func (sdk *LiteApiSdk) CancelBooking(bookingId string) *APIResponse {
	var errors []string
	if bookingId == "" {
		errors = append(errors, "The booking ID is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := fmt.Sprintf("%s/bookings/%s", sdk.BookServiceURL, bookingId)
	return sdk.makeRequest("PUT", url, nil)
}

// GetCitiesByCountryCode returns a list of city names from a specific country
// The API returns a list of city names from a specific country. The country codes needs be is in ISO-2 format. To get the country codes in ISO-2 for all countries please use the GET Country list endpoint
func (sdk *LiteApiSdk) GetCitiesByCountryCode(countryCode string) *APIResponse {
	var errors []string
	if countryCode == "" {
		errors = append(errors, "The country code is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := fmt.Sprintf("%s/data/cities?countryCode=%s", sdk.ServiceURL, countryCode)
	return sdk.makeRequest("GET", url, nil)
}

// GetPlaces looks up for a list of places and areas, given a search query
// Look up for a list of places and areas, given a search query. Places can be used to search for hotels within a location and restrict the list to results within the boundaries of a selected place.
func (sdk *LiteApiSdk) GetPlaces(textQuery, placeType, language string) *APIResponse {
	if language == "" {
		language = "en"
	}

	params := url.Values{}
	params.Add("textQuery", textQuery)
	if placeType != "" {
		params.Add("type", placeType)
	}
	params.Add("language", language)

	url := fmt.Sprintf("%s/data/places?%s", sdk.ServiceURL, params.Encode())
	return sdk.makeRequest("GET", url, nil)
}

// GetCurrencies returns all available currency codes along with its name and supported countries
// The API returns all available currency codes along with its name and the list of supported countries that the currency applies to.
func (sdk *LiteApiSdk) GetCurrencies() *APIResponse {
	url := sdk.ServiceURL + "/data/currencies"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelFacilities returns the list of hotel facilities available in the system
// The API returns the list of hotel facilities available in the system.
func (sdk *LiteApiSdk) GetHotelFacilities() *APIResponse {
	url := sdk.ServiceURL + "/data/facilities"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelTypes returns a list of available hotel types
// The API returns a list of available hotel types.
func (sdk *LiteApiSdk) GetHotelTypes() *APIResponse {
	url := sdk.ServiceURL + "/data/hotelTypes"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelChains returns a list of available hotel chains
// The API returns a list of available hotel chains.
func (sdk *LiteApiSdk) GetHotelChains() *APIResponse {
	url := sdk.ServiceURL + "/data/chains"
	return sdk.makeRequest("GET", url, nil)
}

// GetCountries returns the list of countries available along with its ISO-2 code
// The API returns the list of countries available along with its ISO-2 code.
func (sdk *LiteApiSdk) GetCountries() *APIResponse {
	url := sdk.ServiceURL + "/data/countries"
	return sdk.makeRequest("GET", url, nil)
}

// GetIataCodes returns the IATA codes for all available airports
// The API returns the IATA (International Air Transport Association) codes for all available airports along with the name of the airport, geographical coordinates and country code in ISO-2 format.
func (sdk *LiteApiSdk) GetIataCodes() *APIResponse {
	url := sdk.ServiceURL + "/data/iataCodes"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotels returns a list of hotels available based on different search criteria
// This API endpoint returns a list of hotels available based on different search criterion.
// The minimum required information is the country code in ISO-2 format. The API supports additional search criteria such as city name, geo coordinates, and radius.
// This endpoint provides detailed hotel metadata, including names, addresses, ratings, amenities, and images, facilitating robust hotel search and display features within applications.
func (sdk *LiteApiSdk) GetHotels(parameters map[string]string, language string) *APIResponse {
	return sdk.GetHotelsWithRetry(parameters, language, 3, time.Second)
}

// GetHotelsWithRetry returns a list of hotels with retry logic for rate limiting
func (sdk *LiteApiSdk) GetHotelsWithRetry(parameters map[string]string, language string, retries int, delay time.Duration) *APIResponse {
	params := url.Values{}
	for key, value := range parameters {
		params.Add(key, value)
	}

	if language != "" {
		params.Add("language", language)
	}

	url := fmt.Sprintf("%s/data/hotels?%s", sdk.ServiceURL, params.Encode())
	return sdk.makeRequestWithRetry("GET", url, nil, retries, delay)
}

// GetHotelDetails returns all the static content details of a hotel or property
// The hotel details API returns all the static contents details of a hotel or property if the hotel ID is provided. The static content include name, description, address, amenities, cancellation policies, images and more.
func (sdk *LiteApiSdk) GetHotelDetails(hotelId, language string) *APIResponse {
	var errors []string
	if hotelId == "" {
		errors = append(errors, "The Hotel code is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	params := url.Values{}
	params.Add("hotelId", hotelId)
	if language != "" {
		params.Add("language", language)
	}

	url := fmt.Sprintf("%s/data/hotel?%s", sdk.ServiceURL, params.Encode())
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelReviews retrieves a list of reviews for a specific hotel identified by hotelId
// Deprecated: This method is deprecated and will be removed in future versions. Use GetDataReviews instead.
func (sdk *LiteApiSdk) GetHotelReviews(hotelId string, limit int, getSentiment bool) *APIResponse {
	return sdk.GetDataReviews(hotelId, limit, getSentiment)
}

// GetDataReviews retrieves a list of reviews for a specific hotel identified by hotelId
func (sdk *LiteApiSdk) GetDataReviews(hotelId string, limit int, getSentiment bool) *APIResponse {
	var errors []string
	if hotelId == "" {
		errors = append(errors, "The Hotel code is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	params := url.Values{}
	params.Add("hotelId", hotelId)
	params.Add("limit", strconv.Itoa(limit))
	params.Add("getSentiment", strconv.FormatBool(getSentiment))

	url := fmt.Sprintf("%s/data/reviews?%s", sdk.ServiceURL, params.Encode())

	resp := sdk.makeRequest("GET", url, nil)

	// Handle the special response structure for reviews
	if resp.Status == "success" {
		if dataMap, ok := resp.Data.(map[string]interface{}); ok {
			result := &APIResponse{
				Status: "success",
				Data:   dataMap["data"],
			}
			// Add sentiment analysis if present
			if sentiment, exists := dataMap["sentimentAnalysis"]; exists {
				resultMap := map[string]interface{}{
					"data":              dataMap["data"],
					"sentimentAnalysis": sentiment,
				}
				result.Data = resultMap
			}
			return result
		}
	}

	return resp
}

// GetGuestsIds returns the unique guest ID of a user based on the users email ID
// The guests API returns the unique guest ID of a user based on the users email ID.
func (sdk *LiteApiSdk) GetGuestsIds(guestId string) *APIResponse {
	var errors []string
	if guestId == "" {
		errors = append(errors, "The guestId is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := fmt.Sprintf("%s/guests/%s", sdk.ServiceURL, guestId)
	return sdk.makeRequest("GET", url, nil)
}

// GetGuestsBookings retrieves a list of all bookings associated with a specific guest
// Retrieves a list of all bookings associated with a specific guest, including details about the points earned and cashback applied for each booking.
func (sdk *LiteApiSdk) GetGuestsBookings(guestId string) *APIResponse {
	var errors []string
	if guestId == "" {
		errors = append(errors, "The guestId is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	url := fmt.Sprintf("%s/guests/%s/bookings", sdk.ServiceURL, guestId)
	return sdk.makeRequest("GET", url, nil)
}

// GetVoucherById retrieves a voucher by its ID
func (sdk *LiteApiSdk) GetVoucherById(voucherID string) *APIResponse {
	var errors []string
	if voucherID == "" {
		errors = append(errors, "The voucherID is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Errors: errors,
		}
	}

	return sdk.makeDashboardRequest("GET", fmt.Sprintf("/vouchers/%s", voucherID), nil)
}

// GetVouchers retrieves all available vouchers
func (sdk *LiteApiSdk) GetVouchers() *APIResponse {
	return sdk.makeDashboardRequest("GET", "/vouchers", nil)
}

// CreateVoucher creates a new voucher with the specified details
// Create a new voucher with the specified details, including the voucher code, discount type, value, and validity period. This voucher can then be used by customers.
func (sdk *LiteApiSdk) CreateVoucher(data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("POST", "/vouchers", data)
}

// UpdateVoucher updates the details of an existing voucher
// Update the details of an existing voucher, including the voucher code, discount value, validity period, and more.
func (sdk *LiteApiSdk) UpdateVoucher(id string, data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("PUT", fmt.Sprintf("/vouchers/%s", id), data)
}

// UpdateVoucherStatus updates the status of a voucher
// Update the status of a voucher, typically to activate or deactivate it.
func (sdk *LiteApiSdk) UpdateVoucherStatus(id string, data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("PUT", fmt.Sprintf("/vouchers/%s/status", id), data)
}

// makeDashboardRequest handles requests to the dashboard API with special headers
func (sdk *LiteApiSdk) makeDashboardRequest(method, path string, body interface{}) *APIResponse {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("Failed to marshal request body: %v", err),
			}
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := sdk.DashboardURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", sdk.ApiKey) // Note: Different header case for dashboard

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check status code first, before reading/parsing response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// For error responses, try to read the body for error details
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("HTTP %d: Failed to read error response", resp.StatusCode),
			}
		}

		var errorResult map[string]interface{}
		if err := json.Unmarshal(responseBody, &errorResult); err != nil {
			return &APIResponse{
				Status: "failed",
				Error:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody)),
			}
		}

		if errorResult["error"] != nil {
			return &APIResponse{
				Status: "failed",
				Error:  errorResult["error"],
			}
		}

		return &APIResponse{
			Status: "failed",
			Error:  "Request failed",
		}
	}

	// Only parse response body for successful requests
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	var result interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return &APIResponse{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	return &APIResponse{
		Status: "success",
		Data:   result,
	}
}

// GetLoyalty fetches the current loyalty program information
func (sdk *LiteApiSdk) GetLoyalty() *APIResponse {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("GET", url, nil)
}

// EnableLoyalty enables the loyalty program with specified status and cashback rate
// Once enable the loyalty program with specified status enabled/disabled and cashback rate (e.g. 0.03 = 3% cashback).
func (sdk *LiteApiSdk) EnableLoyalty(data interface{}) *APIResponse {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("POST", url, data)
}

// UpdateLoyalty updates the loyalty program settings
// Updates the loyalty program settings, including status and cashback rates.
func (sdk *LiteApiSdk) UpdateLoyalty(data interface{}) *APIResponse {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("PUT", url, data)
}

// RetrieveWeeklyAnalytics fetches weekly analytics data for the specified date range
func (sdk *LiteApiSdk) RetrieveWeeklyAnalytics(data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("POST", "/analytics/weekly", data)
}

// RetrieveAnalyticsReport fetches a detailed analytics report for the specified date range
func (sdk *LiteApiSdk) RetrieveAnalyticsReport(data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("POST", "/analytics/report", data)
}

// RetrieveMarketAnalytics fetches market analytics data for the specified date range
func (sdk *LiteApiSdk) RetrieveMarketAnalytics(data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("POST", "/analytics/markets", data)
}

// RetrieveMostBookedHotels fetches hotel analytics data for most booked hotels in the specified date range
func (sdk *LiteApiSdk) RetrieveMostBookedHotels(data interface{}) *APIResponse {
	return sdk.makeDashboardRequest("POST", "/analytics/hotels", data)
}
