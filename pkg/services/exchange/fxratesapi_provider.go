package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"larsa-tourism-microservices/pkg/caching"
)

const (
	cacheTTL = 4 * time.Hour
)

// FxRatesAPIProvider implements RateProvider using fxratesapi.com
type FxRatesAPIProvider struct {
	httpClient *http.Client
	baseURL    string
}

// NewFxRatesAPIProvider creates a new FxRatesAPIProvider
func NewFxRatesAPIProvider() *FxRatesAPIProvider {
	return &FxRatesAPIProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.fxratesapi.com/latest?base=USD&resolution=1m&amount=1&places=6&format=json",
	}
}

// Name returns the provider name
func (p *FxRatesAPIProvider) Name() string {
	return "fxratesapi.com"
}

// FxRatesAPIResponse represents the response from fxratesapi.com
type FxRatesAPIResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

// buildCacheKey builds a cache key following the pattern: academy:exchange:rates:c:{clientID}:fxratesapi
func buildCacheKey() string {
	return fmt.Sprintf("tourism:exchange:rates:fxratesapi")
}

// getCachedRates attempts to retrieve rates from cache
func (p *FxRatesAPIProvider) getCachedRates(cacheKey string) (map[string]float64, bool) {
	cached, err := caching.Rdb.GetByKey(cacheKey)
	if err != nil {
		slog.Debug("Cache miss or error", "key", cacheKey, "error", err)
		return nil, false
	}

	// Check if cached value is non-empty
	if strings.TrimSpace(cached) == "" {
		slog.Debug("Cache returned empty value", "key", cacheKey)
		return nil, false
	}

	var rates map[string]float64
	if err := json.Unmarshal([]byte(cached), &rates); err != nil {
		slog.Warn("Failed to unmarshal cached rates", "key", cacheKey, "error", err)
		return nil, false
	}

	// Validate rates map is not empty
	if len(rates) == 0 {
		slog.Debug("Cached rates map is empty", "key", cacheKey)
		return nil, false
	}

	slog.Info("Cache hit", "key", cacheKey, "currencies", len(rates))
	return rates, true
}

// GetRates fetches exchange rates from fxratesapi.com with 4-hour caching
func (p *FxRatesAPIProvider) GetRates() (map[string]float64, error) {
	// Get clientID from context for cache key

	cacheKey := buildCacheKey()

	// Try to get from cache first
	if rates, found := p.getCachedRates(cacheKey); found {
		return rates, nil
	}

	slog.Info("Cache miss, fetching from API", "key", cacheKey)

	// Cache miss, fetch from API
	req, err := http.NewRequest("GET", p.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		// On network error, try to return stale cache if available
		slog.Warn("Network error fetching rates, checking for stale cache", "error", err)
		if staleRates, found := p.getCachedRates(cacheKey); found {
			slog.Info("Returning stale cache due to network error")
			return staleRates, nil
		}
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limit errors (429) by returning stale cache if available
	if resp.StatusCode == http.StatusTooManyRequests {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Warn("Rate limit hit (429), checking for stale cache", "key", cacheKey)

		// Try to get stale cache before failing
		if staleRates, found := p.getCachedRates(cacheKey); found {
			slog.Info("Returning stale cache due to rate limit")
			return staleRates, nil
		}

		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		// For other errors, also try stale cache as fallback
		slog.Warn("API error, checking for stale cache", "status", resp.StatusCode, "key", cacheKey)
		if staleRates, found := p.getCachedRates(cacheKey); found {
			slog.Info("Returning stale cache due to API error")
			return staleRates, nil
		}
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp FxRatesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		// On decode error, try stale cache
		slog.Warn("Failed to decode API response, checking for stale cache", "error", err)
		if staleRates, found := p.getCachedRates(cacheKey); found {
			slog.Info("Returning stale cache due to decode error")
			return staleRates, nil
		}
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		// On API failure, try stale cache
		slog.Warn("API returned success=false, checking for stale cache")
		if staleRates, found := p.getCachedRates(cacheKey); found {
			slog.Info("Returning stale cache due to API failure")
			return staleRates, nil
		}
		return nil, fmt.Errorf("API returned success=false")
	}

	if apiResp.Base != "USD" {
		return nil, fmt.Errorf("expected base currency USD, got %s", apiResp.Base)
	}

	if apiResp.Rates == nil {
		return nil, fmt.Errorf("rates map is nil")
	}

	// Ensure USD is in the rates map
	apiResp.Rates["USD"] = 1.0

	// Cache the result for 4 hours
	caching.Rdb.CacheByKey(cacheKey, apiResp.Rates, cacheTTL)
	slog.Info("Successfully cached rates", "key", cacheKey, "currencies", len(apiResp.Rates))

	return apiResp.Rates, nil
}
