package liteapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"net/http"
	"strconv"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchV2Svcs interface {
	Search(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error)
}

type searchV2Svcs struct {
	lockRepo        repo.LockRepo
	dataFetchedRepo repo.DataFetchRepo
	dataSvcs        DataSvcs
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
	withtxn         *db.WithTxn
	workerPool      chan struct{}
}

func NewSearchV2Svcs(i *do.Injector) (SearchV2Svcs, error) {
	return &searchV2Svcs{
		lockRepo:        do.MustInvoke[repo.LockRepo](i),
		dataFetchedRepo: do.MustInvoke[repo.DataFetchRepo](i),
		dataSvcs:        do.MustInvoke[DataSvcs](i),
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		withtxn:         do.MustInvoke[*db.WithTxn](i),
		workerPool:      make(chan struct{}, 50),
	}, nil
}

func (s *searchV2Svcs) fetchHotelsInBackground(ctx context.Context, placeId, language string) error {
	if err := s.dataFetchedRepo.EnsureIndexes(ctx); err != nil {
		return errors.New("error create index")
	}

	dataFetched, err := s.getDataFetchedInfo(ctx, placeId, language)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("error getting data fetched info")
	}

	if dataFetched != nil && dataFetched.FetchedCount >= dataFetched.TotalCount {
		//we already have the full dataset
		return nil
	}

	lock, err := s.acquireLock(ctx, placeId, language)
	if err != nil {
		return err
	}
	if lock != nil {

		//we aquire the lock successfully - fetch the data from liteapi
		defer s.releaseLock(ctx, lock.Id)

		limit := 5000
		var numOfHotelFetched int
		if dataFetched != nil {
			numOfHotelFetched = dataFetched.FetchedCount
		}

		hotelsQuery := map[string]string{
			"placeId":  placeId,
			"language": language,
		}

		defer func() {
			if dataFetched != nil {
				filter := bson.M{
					"_id": dataFetched.Id,
				}
				update := bson.M{
					"$set": bson.M{
						"fetchedCount": numOfHotelFetched,
					},
				}
				_, err := s.dataFetchedRepo.Patch(ctx, filter, update)
				if err != nil {
					fmt.Println("error update fetch data ")
				}
			}

		}()

		for {
			offset := numOfHotelFetched

			hotelsQuery["offset"] = strconv.Itoa(offset)
			hotelsQuery["limit"] = strconv.Itoa(limit)

			result, err := s.dataSvcs.GetHotelsByPlaceId(ctx, hotelsQuery)
			if err != nil {
				fmt.Println("error fetching hotels from liteapi: %w", err)
				return fmt.Errorf("error fetching hotels from liteapi: %w", err)
			}

			numOfHotelFetched += len(result.Data)

			if dataFetched == nil {
				dataFetched = &models.DataFetch{
					Id:         primitive.NewObjectID(),
					PlaceId:    placeId,
					Language:   language,
					TotalCount: result.Total,
				}

				if err := s.dataFetchedRepo.Add(ctx, dataFetched); err != nil {
					return errors.New("error adding data fetched info")
				}
			}

			if len(result.Data) == 0 {
				//end of loop - no more data to fetch
				break
			}

		}

	}

	return nil

}

func (s *searchV2Svcs) getDataFetchedInfo(ctx context.Context, placeId, language string) (*models.DataFetch, error) {
	filter := bson.M{
		"placeId":  placeId,
		"language": language,
	}
	return s.dataFetchedRepo.GetByFilter(ctx, filter)
}

func (s *searchV2Svcs) pollFromDb(ctx context.Context, w http.ResponseWriter, query map[string]any, placeId, language string) error {
	numOfHotelFetched := 0
	limit := 5000
	var dataFetched *models.DataFetch
	var err error

	ctx2, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx2.Done():
			return errors.New("context timeout")
		default:
		}

		dataFetched, err = s.getDataFetchedInfo(ctx2, placeId, language)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("error getting total count")
		}
		if dataFetched == nil {
			continue
		}

		skip := numOfHotelFetched
		result, err := s.dataSvcs.GetHotelsFromDB(ctx2, placeId, language, skip, limit)
		if err != nil {
			return err
		}

		numOfHotelFetched += len(result.Data)

		if err := s.streamData(ctx2, w, query, result.Data); err != nil {
			return fmt.Errorf("error streaming data [fetch from db]: %w", err)
		}

		if numOfHotelFetched >= dataFetched.TotalCount {
			fmt.Println("fully fetched the hotels form db")
			break
		}

	}

	return nil
}

func (s *searchV2Svcs) Search(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error) {
	language, ok := query["language"].(string)
	if !ok {
		language = "en"
	}
	// country, ok := query["countryCode"].(string)
	// if !ok {
	// 	return nil, errors.New("countryCode is required")
	// }

	placeId, ok := query["placeId"].(string)
	if !ok {
		return nil, errors.New("placeId is required")
	}

	go s.fetchHotelsInBackground(context.WithoutCancel(ctx), placeId, language)

	if err := s.pollFromDb(ctx, w, query, placeId, language); err != nil {
		return nil, err
	}

	return nil, nil

}

func (s *searchV2Svcs) acquireLock(ctx context.Context, placeId, language string) (*models.Lock, error) {
	if err := s.lockRepo.EnsureIndexes(ctx); err != nil {
		fmt.Println("error create index")
		return nil, err
	}

	filter := bson.M{
		"placeId":  placeId,
		"language": language,
		"isLocked": false,
	}

	update := bson.M{
		"$set": bson.M{"isLocked": true},
	}

	upsert := true
	after := options.After
	opts := []*options.FindOneAndUpdateOptions{
		{
			Upsert:         &upsert,
			ReturnDocument: &after,
		},
	}

	doc, err := s.lockRepo.Patch(ctx, filter, update, opts...)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, nil
		}
		fmt.Println("error patch the lock", err)
		return nil, err
	}

	fmt.Println("lock acquired successfully for", placeId, language)
	return doc, nil
}

func (s *searchV2Svcs) releaseLock(ctx context.Context, lockId primitive.ObjectID) error {
	filter := bson.M{
		"_id": lockId,
	}
	update := bson.M{
		"$set": bson.M{"isLocked": false},
	}
	_, err := s.lockRepo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *searchV2Svcs) streamData(ctx context.Context, w http.ResponseWriter, data map[string]any, hotels []models.Hotel) error {
	// Create a map for quick hotel lookup by ID
	hotelMap := make(map[string]models.Hotel, len(hotels))
	ids := make([]string, 0, len(hotels))
	for _, hotel := range hotels {
		ids = append(ids, hotel.Id)
		hotelMap[hotel.Id] = hotel
	}

	ratesData := map[string]any{
		"hotelIds":         ids,
		"occupancies":      data["occupancies"],
		"currency":         "USD",
		"guestNationality": "US",
		"checkin":          data["checkin"],
		"checkout":         data["checkout"],
		//"countryCode":      "USD",
		"stream": true,
	}

	// ratesData := map[string]any{
	// 	"hotelIds":         []string{"lp6556ccb8"},
	// 	"occupancies":      []any{map[string]any{"adults": 1, "children": []any{}}},
	// 	"currency":         "USD",
	// 	"guestNationality": "US",
	// 	"checkin":          "2025-11-23",
	// 	"checkout":         "2025-11-24",
	// 	"countryCode":      "USD",
	// 	"stream":           true,
	// }

	liteApiSdk, err := s.liteApiInitFunc(ctx)
	if err != nil {
		return err
	}

	ch, err := liteApiSdk.GetFullRatesStream(ratesData)
	if err != nil {
		return err
	}

	// Set headers for Server-Sent Events (SSE)
	// These headers tell the browser to keep the connection open and expect streaming data
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the flusher to flush data immediately
	// The flusher is a type assertion that checks if 'w' supports immediate flushing
	// Most HTTP servers (like Go's net/http) support this
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming not supported")
	}

	for event := range ch {
		if event.Error != nil {
			// Send error and terminate stream
			errorData, _ := json.Marshal(map[string]any{
				"error": event.Error.Error(),
			})
			fmt.Fprintf(w, "data: %s\n\n", errorData)
			flusher.Flush()
			return event.Error
		}

		if event.Done {
			// Send [DONE] signal to indicate completion
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			break
		}

		var result map[string]any
		if err := json.Unmarshal(event.Data, &result); err != nil {
			// Send parse error and terminate
			errorData, _ := json.Marshal(map[string]any{
				"error": "failed to parse rate data",
			})
			fmt.Fprintf(w, "data: %s\n\n", errorData)
			flusher.Flush()
			return err
		}

		// Add related hotel information to the rate
		if hotelId, ok := result["hotelId"].(string); ok {
			if hotel, exists := hotelMap[hotelId]; exists {
				result["hotel"] = hotel
			}
		}

		// Stream each rate to the frontend
		rateJSON, err := json.Marshal(result)
		if err != nil {
			return err
		}

		// Write SSE data event to the HTTP response buffer
		// Format: "data: {json}\n\n" - compatible with JS client
		// The client expects simple "data: " prefix without event types
		fmt.Fprintf(w, "data: %s\n\n", rateJSON)

		// Flush immediately to send data in real-time
		// This ensures the client receives each rate as it becomes available
		flusher.Flush()
	}

	return nil
}
