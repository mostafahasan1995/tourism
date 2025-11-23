package liteApiSdk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	Status string          `json:"status"` // success or failed
	Code   int             `json:"code"`   // http status code
	Data   json.RawMessage `json:"data"`   //response data on success
	Err    map[string]any  `json:"error"`  // response error on failure

}

// NewLiteApiSdk creates a new instance of the LiteApi SDK
func NewLiteApiSdk(apiKey string) *LiteApiSdk {
	return &LiteApiSdk{
		ApiKey:         apiKey,
		ServiceURL:     "https://api.liteapi.travel/v3.0",
		BookServiceURL: "https://book.liteapi.travel/v3.0",
		DashboardURL:   "https://da.liteapi.travel",
		Client: &http.Client{
			Timeout: 5 * time.Minute, // Increased from 30s to 5min for large data operations
		},
	}
}

func (sdk *LiteApiSdk) getBaseUrl(url string) string {
	urlMap := map[string]string{
		"data":                 sdk.ServiceURL,
		"hotels":               sdk.ServiceURL,
		"rates":                sdk.BookServiceURL,
		"bookings":             sdk.BookServiceURL,
		"guests":               sdk.ServiceURL,
		"loyalties":            sdk.ServiceURL,
		"vouchers":             sdk.DashboardURL,
		"analytics":            sdk.DashboardURL,
		"supply-customization": sdk.ServiceURL,
	}

	for key, val := range urlMap {
		if strings.HasPrefix(url, key) {
			return val
		}
	}

	return ""

}

// makeRequest handles HTTP requests with common headers and error handling
// Returns either error (for non-response errors) or APIResponse (for actual API responses)
func (sdk *LiteApiSdk) makeRequest(method, url string, body any) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", sdk.ApiKey)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// We have a valid response (either success or API error)
	if resp.StatusCode == http.StatusOK {
		return &APIResponse{
			Status: "success",
			Code:   resp.StatusCode,
			Data:   responseBody, // Success data goes into Data field as json.RawMessage
		}, nil
	} else {
		var result map[string]any
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %v", err)
		}
		return &APIResponse{
			Status: "failed",
			Code:   resp.StatusCode,
			Err:    result, // API error response goes into Err field
		}, nil
	}
}

// StreamEvent represents a single event from the streaming API
type StreamEvent struct {
	Data  json.RawMessage `json:"data"`
	Error error           `json:"error,omitempty"`
	Done  bool            `json:"done"`
}

// GetFullRates searches and returns all available rooms along with rates and cancellation policies
// The Full Rates API is to search and return all available rooms along with its rates, cancellation policies for a list of hotel ID's based on the search dates.
// For each hotel ID, all available room information is returned.
// The API also has a built in loyalty rewards system. The system rewards return users who have made prior bookings.
// If the search is coming from a known guest ID, the guest level is also returned along with the pricing that's appropriate for the guest level.
// If it is a new user, the guest ID will be generated at the time of the first confirmed booking.
func (sdk *LiteApiSdk) GetFullRates(data interface{}) (*APIResponse, error) {
	url := sdk.ServiceURL + "/hotels/rates"
	return sdk.makeRequest("POST", url, data)
}

// GetFullRatesStream searches and returns available rooms as a stream when stream=true is set in the request body
// This method returns a channel that yields streaming events. The channel will be closed when the stream ends or an error occurs.
// Usage:
//
//	eventChan, err := sdk.GetFullRatesStream(data)
//	if err != nil { handle error }
//	for event := range eventChan {
//	  if event.Error != nil { handle error }
//	  if event.Done { break }
//	  // process event.Data
//	}
func (sdk *LiteApiSdk) GetFullRatesStream(data interface{}) (<-chan StreamEvent, error) {
	url := sdk.ServiceURL + "/hotels/rates"
	return sdk.makeStreamRequest("POST", url, data)
}

// makeStreamRequest handles HTTP streaming requests
func (sdk *LiteApiSdk) makeStreamRequest(method, url string, body any) (<-chan StreamEvent, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", sdk.ApiKey)

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Check if response is not OK
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response: %v", err)
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Create channel for streaming events
	eventChan := make(chan StreamEvent, 10)

	// Start goroutine to read stream
	go func() {
		defer close(eventChan)
		defer resp.Body.Close()

		// Check if response is SSE format by looking at Content-Type
		contentType := resp.Header.Get("Content-Type")
		isSSE := strings.Contains(contentType, "text/event-stream")

		// Also check for streaming parameter in the request to detect SSE
		// LiteAPI might return SSE when stream=true but not set proper Content-Type
		if !isSSE {
			// Read a small sample to detect format
			peek := make([]byte, 10)
			n, _ := resp.Body.Read(peek)
			if n > 0 {
				// Check if response starts with "data:" which indicates SSE
				if strings.HasPrefix(string(peek[:n]), "data:") {
					isSSE = true
				}
			}

			// Create a new reader with the peeked data prepended
			resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(peek[:n]), resp.Body))
		}

		if isSSE {
			// Handle SSE (Server-Sent Events) format
			// Using bufio.Reader instead of Scanner to handle UNLIMITED line sizes
			// Scanner has buffer limits, Reader.ReadString does not
			reader := bufio.NewReader(resp.Body)

			for {
				// Read line by line with NO SIZE LIMIT
				line, err := reader.ReadString('\n')

				// Handle EOF
				if err == io.EOF {
					// If we got data before EOF, process it
					if line != "" {
						line = strings.TrimSpace(line)
						if strings.HasPrefix(line, "data:") {
							jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
							if jsonData != "" {
								eventChan <- StreamEvent{Data: json.RawMessage(jsonData), Done: false}
							}
						}
					}
					// Stream ended normally
					eventChan <- StreamEvent{Done: true}
					return
				}

				// Handle other errors
				if err != nil {
					eventChan <- StreamEvent{Error: fmt.Errorf("error reading SSE stream: %v", err), Done: true}
					return
				}

				// Trim whitespace and newlines
				line = strings.TrimSpace(line)

				// Skip empty lines (SSE uses empty lines as separators)
				if line == "" {
					continue
				}

				// SSE format: "data: {json}"
				if strings.HasPrefix(line, "data:") {
					// Extract JSON after "data: "
					jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

					// Skip if data is empty
					if jsonData == "" {
						continue
					}

					// Send the JSON data - NO SIZE LIMIT!
					eventChan <- StreamEvent{Data: json.RawMessage(jsonData), Done: false}
				}
				// Skip event type lines and other SSE metadata
			}

		} else {
			// Handle newline-delimited JSON format
			decoder := json.NewDecoder(resp.Body)

			for {
				var chunk json.RawMessage
				if err := decoder.Decode(&chunk); err != nil {
					if err == io.EOF {
						// Stream ended normally
						eventChan <- StreamEvent{Done: true}
						return
					}
					// Error occurred
					eventChan <- StreamEvent{Error: fmt.Errorf("failed to decode stream: %v", err), Done: true}
					return
				}

				// Send the chunk
				eventChan <- StreamEvent{Data: chunk, Done: false}
			}
		}
	}()

	return eventChan, nil
}

// GetMinRates gets minimum rates for hotels
func (sdk *LiteApiSdk) GetMinRates(data interface{}) (*APIResponse, error) {
	url := sdk.ServiceURL + "/hotels/min-rates"
	return sdk.makeRequest("POST", url, data)
}

// PreBook confirms if the room and rates for the search criterion
// This API is used to confirm if the room and rates for the search criterion. The input to the endpoint is an array of rate Ids coming from the GET hotel full rates availability API.
// In response, the API generates a prebook Id, a new rate Id and contains information if price, cancellation policy or boarding information has changed.
func (sdk *LiteApiSdk) PreBook(data map[string]any) (*APIResponse, error) {
	var errors []string

	// Validate offerId
	if offerId, exists := data["offerId"]; !exists || offerId == "" {
		errors = append(errors, "The offerId is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := sdk.BookServiceURL + "/rates/prebook"
	return sdk.makeRequest("POST", url, data)
}

// Book confirms a booking when the prebook Id and the rate Id from the pre book stage along with the guest and payment information are passed
// This API confirms a booking when the prebook Id and the rate Id from the pre book stage along with the guest and payment information are passed.
// The guest information is an object that should include the guest first name, last name and email.
// The payment information is an object that should include the name, credit card number, expiry and CVC number.
// The response will confirm the booking along with a booking Id and a hotel confirmation code. It will also include the booking details including the dates, price and the cancellation policies.
func (sdk *LiteApiSdk) Book(data map[string]any) (*APIResponse, error) {
	var errors []string

	// Validate prebookId
	if prebookId, exists := data["prebookId"]; !exists || prebookId == "" {
		errors = append(errors, "The prebookId is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := sdk.BookServiceURL + "/rates/book"
	return sdk.makeRequest("POST", url, data)
}

// GetBookingsList returns the list of booking Id's for a given guest Id
// The API returns the list of booking Id's for a given guest Id.
func (sdk *LiteApiSdk) GetBookingsList(clientReference string) (*APIResponse, error) {
	url := fmt.Sprintf("%s/bookings?clientReference=%s", sdk.BookServiceURL, clientReference)
	return sdk.makeRequest("GET", url, nil)
}

// GetBookingsByEmail returns bookings for a given email (clientReference)
// This method uses email as clientReference to fetch bookings
func (sdk *LiteApiSdk) GetBookingsByEmail(email string) (*APIResponse, error) {
	if email == "" {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": []string{"The email is required"},
			},
		}, nil
	}
	return sdk.GetBookingsList(email)
}

// GetBookingsByDateRange returns all bookings for a given date range
// This method calls the LiteAPI bookings endpoint with startDate and endDate parameters
func (sdk *LiteApiSdk) GetBookingsByDateRange(startDate, endDate string) (*APIResponse, error) {
	params := url.Values{}
	if startDate != "" {
		params.Add("startDate", startDate)
	}
	if endDate != "" {
		params.Add("endDate", endDate)
	}

	url := fmt.Sprintf("%s/bookings/?%s", sdk.ServiceURL, params.Encode())
	return sdk.makeRequest("GET", url, nil)
}

// RetrieveBooking returns the status and the details for a specific booking Id
// The API returns the status and the details for the a specific booking Id.
func (sdk *LiteApiSdk) RetrieveBooking(bookingId string) (*APIResponse, error) {
	var errors []string
	if bookingId == "" {
		errors = append(errors, "The booking ID is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := fmt.Sprintf("%s/bookings/%s", sdk.BookServiceURL, bookingId)
	return sdk.makeRequest("GET", url, nil)
}

// CancelBooking requests a cancellation of an existing confirmed booking
// This API is used to request a cancellation of an existing confirmed booking. Cancellation policies and conditions will be used to determine the success of the cancellation. For example a booking with non-refundable (NRFN) tag or a booking with a cancellation policy that was requested past the cancellation date will not be able to cancel the confirmed booking.
func (sdk *LiteApiSdk) CancelBooking(bookingId string) (*APIResponse, error) {
	var errors []string
	if bookingId == "" {
		errors = append(errors, "The booking ID is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := fmt.Sprintf("%s/bookings/%s", sdk.BookServiceURL, bookingId)
	return sdk.makeRequest("PUT", url, nil)
}

// GetCitiesByCountryCode returns a list of city names from a specific country
// The API returns a list of city names from a specific country. The country codes needs be is in ISO-2 format. To get the country codes in ISO-2 for all countries please use the GET Country list endpoint
func (sdk *LiteApiSdk) GetCitiesByCountryCode(countryCode string) (*APIResponse, error) {
	var errors []string
	if countryCode == "" {
		errors = append(errors, "The country code is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := fmt.Sprintf("%s/data/cities?countryCode=%s", sdk.ServiceURL, countryCode)
	return sdk.makeRequest("GET", url, nil)
}

// GetPlaces looks up for a list of places and areas, given a search query
// Look up for a list of places and areas, given a search query. Places can be used to search for hotels within a location and restrict the list to results within the boundaries of a selected place.
func (sdk *LiteApiSdk) GetPlaces(textQuery, placeType, language string) (*APIResponse, error) {
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

func (sdk *LiteApiSdk) GetPlace(placeId string) (*APIResponse, error) {
	url := fmt.Sprintf("%s/data/places/%s", sdk.ServiceURL, placeId)
	return sdk.makeRequest("GET", url, nil)
}

// GetCurrencies returns all available currency codes along with its name and supported countries
// The API returns all available currency codes along with its name and the list of supported countries that the currency applies to.
func (sdk *LiteApiSdk) GetCurrencies() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/currencies"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelFacilities returns the list of hotel facilities available in the system
// The API returns the list of hotel facilities available in the system.
func (sdk *LiteApiSdk) GetHotelFacilities() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/facilities"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelTypes returns a list of available hotel types
// The API returns a list of available hotel types.
func (sdk *LiteApiSdk) GetHotelTypes() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/hotelTypes"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelChains returns a list of available hotel chains
// The API returns a list of available hotel chains.
func (sdk *LiteApiSdk) GetHotelChains() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/chains"
	return sdk.makeRequest("GET", url, nil)
}

// GetCountries returns the list of countries available along with its ISO-2 code
// The API returns the list of countries available along with its ISO-2 code.
func (sdk *LiteApiSdk) GetCountries() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/countries"
	return sdk.makeRequest("GET", url, nil)
}

// GetIataCodes returns the IATA codes for all available airports
// The API returns the IATA (International Air Transport Association) codes for all available airports along with the name of the airport, geographical coordinates and country code in ISO-2 format.
func (sdk *LiteApiSdk) GetIataCodes() (*APIResponse, error) {
	url := sdk.ServiceURL + "/data/iataCodes"
	return sdk.makeRequest("GET", url, nil)
}

// GetHotels returns a list of hotels available based on different search criteria
// This API endpoint returns a list of hotels available based on different search criterion.
// The minimum required information is the country code in ISO-2 format. The API supports additional search criteria such as city name, geo coordinates, and radius.
// This endpoint provides detailed hotel metadata, including names, addresses, ratings, amenities, and images, facilitating robust hotel search and display features within applications.
func (sdk *LiteApiSdk) GetHotels(parameters map[string]string, retries int, delay time.Duration) (*APIResponse, error) {
	params := url.Values{}
	for key, value := range parameters {
		params.Add(key, value)
	}

	url := fmt.Sprintf("%s/data/hotels?%s", sdk.ServiceURL, params.Encode())
	resp, err := sdk.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		// Check for rate limiting
		if resp.Code == 429 || (resp.Err != nil && resp.Err["code"] == 4290) {
			if retries > 0 {
				time.Sleep(delay)
				return sdk.GetHotels(parameters, retries-1, delay*2)
			} else {
				return &APIResponse{
					Status: "failed",
					Code:   resp.Code,
					Err: map[string]any{
						"error": "Rate limit exceeded",
					},
				}, nil
			}
		}
		return resp, nil
	}
	return resp, nil

}

// GetHotelDetails returns all the static content details of a hotel or property
// The hotel details API returns all the static contents details of a hotel or property if the hotel ID is provided. The static content include name, description, address, amenities, cancellation policies, images and more.
func (sdk *LiteApiSdk) GetHotelDetails(hotelId, language, advancedAccessibilityOnly string) (*APIResponse, error) {
	var errors []string
	if hotelId == "" {
		errors = append(errors, "The Hotel id is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	params := url.Values{}
	params.Add("hotelId", hotelId)
	if language != "" {
		params.Add("language", language)
	}
	if advancedAccessibilityOnly != "" {
		params.Add("advancedAccessibilityOnly", advancedAccessibilityOnly)
	}

	url := fmt.Sprintf("%s/data/hotel?%s", sdk.ServiceURL, params.Encode())
	return sdk.makeRequest("GET", url, nil)
}

// GetHotelReviews retrieves a list of reviews for a specific hotel identified by hotelId
func (sdk *LiteApiSdk) GetHotelReviews(parameters map[string]string) (*APIResponse, error) {
	var errors []string
	if parameters["hotelId"] == "" {
		errors = append(errors, "The Hotel id is required")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	params := url.Values{}
	for key, value := range parameters {
		params.Add(key, value)
	}

	url := fmt.Sprintf("%s/data/reviews?%s", sdk.ServiceURL, params.Encode())

	return sdk.makeRequest("GET", url, nil)
}

// GetDataReviews retrieves a list of reviews for a specific hotel identified by hotelId
// func (sdk *LiteApiSdk) GetDataReviews(hotelId string, limit int, getSentiment bool) (*APIResponse, error) {
// 	var errors []string
// 	if hotelId == "" {
// 		errors = append(errors, "The Hotel id is required")
// 	}

// 	if len(errors) > 0 {
// 		return &APIResponse{
// 			Status: "failed",
// 			Code:   http.StatusBadRequest,
// 			Err: map[string]any{
// 				"errors": errors,
// 			},
// 		}, nil
// 	}

// 	params := url.Values{}
// 	params.Add("hotelId", hotelId)
// 	params.Add("limit", strconv.Itoa(limit))
// 	params.Add("getSentiment", strconv.FormatBool(getSentiment))

// 	url := fmt.Sprintf("%s/data/reviews?%s", sdk.ServiceURL, params.Encode())

// 	resp, err := sdk.makeRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Handle the special response structure for reviews
// 	if resp.Status == "success" && resp.Data != nil {
// 		// Parse the raw data
// 		var reviewData map[string]any
// 		if err := json.Unmarshal(resp.Data, &reviewData); err == nil {
// 			resultData := make(map[string]any)

// 			// Add sentiment analysis if present
// 			if sentiment, exists := reviewData["sentimentAnalysis"]; exists {
// 				resultData["data"] = reviewData["data"]
// 				resultData["sentimentAnalysis"] = sentiment
// 			} else {
// 				resultData = reviewData
// 			}

// 			resultBytes, _ := json.Marshal(resultData)
// 			return &APIResponse{
// 				Status: "success",
// 				Code:   resp.Code,
// 				Data:   resultBytes,
// 			}, nil
// 		}
// 	}

// 	return resp, nil
// }

// GetGuestsIds returns the unique guest ID of a user based on the users email ID
// The guests API returns the unique guest ID of a user based on the users email ID.
func (sdk *LiteApiSdk) GetGuestsIds(guestId string) (*APIResponse, error) {
	var errors []string
	if guestId == "" {
		errors = append(errors, "The guestId is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := fmt.Sprintf("%s/guests/%s", sdk.ServiceURL, guestId)
	return sdk.makeRequest("GET", url, nil)
}

// GetGuestsBookings retrieves a list of all bookings associated with a specific guest
// Retrieves a list of all bookings associated with a specific guest, including details about the points earned and cashback applied for each booking.
func (sdk *LiteApiSdk) GetGuestsBookings(guestId string) (*APIResponse, error) {
	var errors []string
	if guestId == "" {
		errors = append(errors, "The guestId is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	url := fmt.Sprintf("%s/guests/%s/bookings", sdk.ServiceURL, guestId)
	return sdk.makeRequest("GET", url, nil)
}

// GetVoucherById retrieves a voucher by its ID
func (sdk *LiteApiSdk) GetVoucherById(voucherID string) (*APIResponse, error) {
	var errors []string
	if voucherID == "" {
		errors = append(errors, "The voucherID is required.")
	}

	if len(errors) > 0 {
		return &APIResponse{
			Status: "failed",
			Code:   http.StatusBadRequest,
			Err: map[string]any{
				"errors": errors,
			},
		}, nil
	}

	return sdk.makeDashboardRequest("GET", fmt.Sprintf("/vouchers/%s", voucherID), nil)
}

// GetVouchers retrieves all available vouchers
func (sdk *LiteApiSdk) GetVouchers() (*APIResponse, error) {
	return sdk.makeDashboardRequest("GET", "/vouchers", nil)
}

// CreateVoucher creates a new voucher with the specified details
// Create a new voucher with the specified details, including the voucher code, discount type, value, and validity period. This voucher can then be used by customers.
func (sdk *LiteApiSdk) CreateVoucher(data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("POST", "/vouchers", data)
}

// UpdateVoucher updates the details of an existing voucher
// Update the details of an existing voucher, including the voucher code, discount value, validity period, and more.
func (sdk *LiteApiSdk) UpdateVoucher(id string, data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("PUT", fmt.Sprintf("/vouchers/%s", id), data)
}

// UpdateVoucherStatus updates the status of a voucher
// Update the status of a voucher, typically to activate or deactivate it.
func (sdk *LiteApiSdk) UpdateVoucherStatus(id string, data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("PUT", fmt.Sprintf("/vouchers/%s/status", id), data)
}

// makeDashboardRequest handles requests to the dashboard API with special headers
func (sdk *LiteApiSdk) makeDashboardRequest(method, path string, body interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := sdk.DashboardURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", sdk.ApiKey) // Note: Different header case for dashboard

	resp, err := sdk.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// We have a valid response (either success or API error)
	if resp.StatusCode == http.StatusOK {
		return &APIResponse{
			Status: "success",
			Code:   resp.StatusCode,
			Data:   responseBody, // Success data goes into Data field as json.RawMessage
		}, nil
	} else {
		return &APIResponse{
			Status: "failed",
			Code:   resp.StatusCode,
			Err:    result, // API error response goes into Err field
		}, nil
	}
}

// GetLoyalty fetches the current loyalty program information
func (sdk *LiteApiSdk) GetLoyalty() (*APIResponse, error) {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("GET", url, nil)
}

// EnableLoyalty enables the loyalty program with specified status and cashback rate
// Once enable the loyalty program with specified status enabled/disabled and cashback rate (e.g. 0.03 = 3% cashback).
func (sdk *LiteApiSdk) EnableLoyalty(data interface{}) (*APIResponse, error) {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("POST", url, data)
}

// UpdateLoyalty updates the loyalty program settings
// Updates the loyalty program settings, including status and cashback rates.
func (sdk *LiteApiSdk) UpdateLoyalty(data interface{}) (*APIResponse, error) {
	url := sdk.ServiceURL + "/loyalties/"
	return sdk.makeRequest("PUT", url, data)
}

// RetrieveWeeklyAnalytics fetches weekly analytics data for the specified date range
func (sdk *LiteApiSdk) RetrieveWeeklyAnalytics(data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("POST", "/analytics/weekly", data)
}

// RetrieveAnalyticsReport fetches a detailed analytics report for the specified date range
func (sdk *LiteApiSdk) RetrieveAnalyticsReport(data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("POST", "/analytics/report", data)
}

// RetrieveMarketAnalytics fetches market analytics data for the specified date range
func (sdk *LiteApiSdk) RetrieveMarketAnalytics(data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("POST", "/analytics/markets", data)
}

// RetrieveMostBookedHotels fetches hotel analytics data for most booked hotels in the specified date range
func (sdk *LiteApiSdk) RetrieveMostBookedHotels(data interface{}) (*APIResponse, error) {
	return sdk.makeDashboardRequest("POST", "/analytics/hotels", data)
}

func (sdk *LiteApiSdk) Request(method, url string, body any) (*APIResponse, error) {

	url = strings.TrimPrefix(url, "/liteapi/")

	baseUrl := sdk.getBaseUrl(url)

	completeUrl := fmt.Sprintf("%s/%s", baseUrl, url)

	return sdk.makeRequest(method, completeUrl, body)
}
