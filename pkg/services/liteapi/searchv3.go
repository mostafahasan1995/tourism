package liteapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/samber/do"
)

type Searchv3Svcs interface {
	Search(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error)
}

type searchv3Svcs struct {
	lockRepo        repo.LockRepo
	dataFetchedRepo repo.DataFetchRepo
	dataSvcs        DataSvcs
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
	withtxn         *db.WithTxn
	workerPool      chan struct{}
	testCh          chan int
}

func NewSearchv3Svcs(i *do.Injector) (Searchv3Svcs, error) {
	return &searchv3Svcs{
		lockRepo:        do.MustInvoke[repo.LockRepo](i),
		dataFetchedRepo: do.MustInvoke[repo.DataFetchRepo](i),
		dataSvcs:        do.MustInvoke[DataSvcs](i),
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		withtxn:         do.MustInvoke[*db.WithTxn](i),
		workerPool:      make(chan struct{}, 2),
		testCh:          make(chan int, 100),
	}, nil
}

func (s *searchv3Svcs) streamRates(ctx context.Context, w http.ResponseWriter, placeId, language string, query map[string]any) error {
	numOfFetchedHotels := 0
	limit := 5000

	for {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		offset := numOfFetchedHotels
		hotels, err := s.dataSvcs.GetHotelsByPlaceIdFromDB(ctx, placeId, language, offset, limit)
		if err != nil {
			return err
		}
		numOfFetchedHotels += len(hotels.Data)
		if len(hotels.Data) == 0 {
			break
		}

		if err := s.stream(ctx, w, hotels.Data, placeId, language, query); err != nil {
			return err
		}

	}

	if err := s.serveFromLiteApiv2(ctx, w, numOfFetchedHotels, placeId, language, query); err != nil {
		return err
	}

	return nil
}

func (s *searchv3Svcs) serveFromLiteApiv2(ctx context.Context, w http.ResponseWriter, start int, placeId, language string, query map[string]any) error {

	numOfFetchedHotels := start

	limit := 5000

	hotelsQuery := map[string]string{
		"placeId":  placeId,
		"language": language,
	}

	for {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		offset := numOfFetchedHotels

		hotelsQuery["offset"] = strconv.Itoa(offset)
		hotelsQuery["limit"] = strconv.Itoa(limit)

		result, err := s.dataSvcs.GetHotelsByPlaceId(ctx, hotelsQuery)
		if err != nil {
			return err
		}
		numOfFetchedHotels += len(result.Data)
		if len(result.Data) == 0 {
			break
		}

		if err := s.stream(ctx, w, result.Data, placeId, language, query); err != nil {
			return err
		}

	}
	return nil
}

func (s *searchv3Svcs) Search(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error) {
	language, ok := query["language"].(string)
	if !ok {
		language = "en"
	}

	placeId, ok := query["placeId"].(string)
	if !ok {
		return nil, errors.New("placeId is required")
	}

	// Set SSE headers once at the start
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, x-client, x-access-token, x-service, x-expire, x-service-token, x-user-id")

	// Get the flusher for sending completion signal
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
	}

	if err := s.streamRates(ctx, w, placeId, language, query); err != nil {
		return nil, err
	}

	// Send [DONE] signal to indicate completion of entire stream
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()

	return nil, nil

}

func (s *searchv3Svcs) stream(ctx context.Context, w http.ResponseWriter, hotels []models.Hotel, placeId, language string, data map[string]any) error {
	// Check if context is canceled before processing
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Skip if no hotels to process
	if len(hotels) == 0 {
		return nil
	}

	m := make(map[string]models.Hotel, len(hotels))

	hotelIds := make([]string, 0, len(hotels))

	for _, hotel := range hotels {
		m[hotel.Id] = hotel
		hotelIds = append(hotelIds, hotel.Id)
	}

	ratesData := map[string]any{
		"hotelIds":         hotelIds,
		"occupancies":      data["occupancies"],
		"currency":         "USD",
		"guestNationality": data["guestNationality"],
		"checkin":          data["checkin"],
		"checkout":         data["checkout"],
		"maxRatesPerHotel": 1,
	}

	liteApiSdk, err := s.liteApiInitFunc(ctx)
	if err != nil {
		return err
	}

	resp, err := liteApiSdk.GetFullRates(ratesData)
	if err != nil {
		return err
	}

	if resp.Status == "failed" {
		return helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.RatesList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return err
	}

	// Get the flusher to flush data immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming not supported")
	}

	streamData := [][]byte{}

	for _, rate := range result.Data {
		if hotelId, ok := rate["hotelId"].(string); ok {
			if hotel, exists := m[hotelId]; exists {
				rate["hotel"] = hotel

				rateJSON, err := json.Marshal(rate)
				if err != nil {
					return err
				}

				streamData = append(streamData, rateJSON)
			}
		}
	}

	// Only stream if there's data to send
	if len(streamData) > 0 {
		chunks := util.SliceChunk(streamData, 50)
		for _, chunk := range chunks {
			// Check if context is canceled before streaming each chunk
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Manually construct JSON array from already-encoded JSON objects
			fmt.Fprint(w, "data: [")
			for i, item := range chunk {
				if i > 0 {
					fmt.Fprint(w, ",")
				}
				fmt.Fprintf(w, "%s", item)
			}
			fmt.Fprint(w, "]\n\n")
			flusher.Flush()
		}
	}

	return nil
}
