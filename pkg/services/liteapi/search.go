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

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchSvcs interface {
	SearchHotels(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error)
}

type searchSvcs struct {
	lockRepo        repo.LockRepo
	dataFetchedRepo repo.DataFetchRepo
	dataSvcs        DataSvcs
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
	withtxn         *db.WithTxn
}

func NewSearchSvcs(i *do.Injector) (SearchSvcs, error) {
	return &searchSvcs{
		lockRepo:        do.MustInvoke[repo.LockRepo](i),
		dataFetchedRepo: do.MustInvoke[repo.DataFetchRepo](i),
		dataSvcs:        do.MustInvoke[DataSvcs](i),
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
		withtxn:         do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (s *searchSvcs) CheckIfDataExist(ctx context.Context, country, language string) (*models.DataFetch, error) {
	filter := bson.M{
		"country":      country,
		"language":     language,
		"fullyFetched": true,
	}
	return s.dataFetchedRepo.GetByFilter(ctx, filter)
}

func (s *searchSvcs) serveHotelsFromDb(ctx context.Context, w http.ResponseWriter, query map[string]any, datafetched *models.DataFetch, country, language string) error {

	totalCount := datafetched.Count
	limit := 5000
	numOfHotelFetched := 0
	for {
		skip := numOfHotelFetched
		result, err := s.dataSvcs.GetHotelsFromDB(ctx, country, language, skip, limit)
		if err != nil {
			return err
		}
		numOfHotelFetched += len(result.Data)

		if err := s.streamData(ctx, w, query, result.Data); err != nil {
			return errors.New("error streaming data [fetch from db]")
		}

		if numOfHotelFetched >= totalCount {
			break
		}
	}

	return nil

}

func (s *searchSvcs) SearchHotels(ctx context.Context, w http.ResponseWriter, query map[string]any) (any, error) {
	//query e.g {"country":"US", "city":"New York", "checkIn":"2025-01-01", "checkOut":"2025-01-02"}

	lang, ok := query["language"].(string)
	if !ok {
		lang = "en"
	}
	country, ok := query["countryCode"].(string)
	if !ok {
		return nil, errors.New("countryCode is required")
	}

	dataFetched, err := s.CheckIfDataExist(ctx, country, lang)

	if err == nil {
		//the data is already fetched and ready - get form db

		if err := s.serveHotelsFromDb(ctx, w, query, dataFetched, country, lang); err != nil {
			return nil, errors.New("error serving hotels from db")
		}

	} else {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//the data is either not fetched at all - or is not fully fetched yet
			// try aquire the lock

			lock, err := s.acquireLock(ctx, country, lang)
			if err != nil {
				return nil, err
			}
			if lock != nil {

				defer func() {
					s.releaseLock(ctx, lock.Id)
				}()
				//we aquire the lock successfully - fetch the data from liteapi

				limit := 5000
				numOfHotelFetched := 0

				var dataFetched *models.DataFetch

				hotelsQuery := map[string]string{
					"countryCode": country,
					"language":    lang,
				}

				for {
					offset := numOfHotelFetched

					hotelsQuery["offset"] = strconv.Itoa(offset)
					hotelsQuery["limit"] = strconv.Itoa(limit)

					result, err := s.dataSvcs.GetHotels(ctx, hotelsQuery)
					if err != nil {
						return nil, fmt.Errorf("error fetching hotels from liteapi: %w", err)
					}

					numOfHotelFetched += len(result.Data)

					if dataFetched == nil {
						dataFetched = &models.DataFetch{
							Id:           primitive.NewObjectID(),
							Country:      country,
							Language:     lang,
							FullyFetched: false,
							Count:        result.Total,
						}
						if err := s.dataFetchedRepo.Add(ctx, dataFetched); err != nil {
							return nil, err
						}

					}

					if err := s.streamData(ctx, w, query, result.Data); err != nil {
						return nil, fmt.Errorf("error streaming data [fetch from liteapi]: %w", err)
					}

					if len(result.Data) == 0 {
						//end of loop - no more data to fetch
						break
					}

				}

				filter := bson.M{
					"_id": dataFetched.Id,
				}
				update := bson.M{
					"$set": bson.M{
						"fullyFetched": true,
						"count":        numOfHotelFetched,
					},
				}
				_, err := s.dataFetchedRepo.Patch(ctx, filter, update)
				if err != nil {
					return nil, err
				}

			} else {
				//fetch form db
				numOfHotelFetched := 0
				limit := 5000
				var dataFetched *models.DataFetch

				for {

					if dataFetched == nil || !dataFetched.FullyFetched {
						f, err := s.getDataFetchedInfo(ctx, country, lang)
						if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
							return nil, errors.New("error getting total count")
						}
						if f == nil {
							continue
						}

						dataFetched = f
					}
					if dataFetched != nil {
						skip := numOfHotelFetched
						result, err := s.dataSvcs.GetHotelsFromDB(ctx, country, lang, skip, limit)
						if err != nil {
							return nil, err
						}

						if err := s.streamData(ctx, w, query, result.Data); err != nil {
							return nil, errors.New("error streaming data [fetch from db]")
						}
						numOfHotelFetched += len(result.Data)

						if dataFetched.FullyFetched {
							if numOfHotelFetched >= dataFetched.Count {
								break
							}
						}

					}

				}
			}

		} else {
			return nil, err
		}

	}

	return nil, nil
}

func (s *searchSvcs) getDataFetchedInfo(ctx context.Context, country, language string) (*models.DataFetch, error) {
	filter := bson.M{
		"country":  country,
		"language": language,
	}
	return s.dataFetchedRepo.GetByFilter(ctx, filter)
}

func (s *searchSvcs) acquireLock(ctx context.Context, country, language string) (*models.Lock, error) {
	if err := s.lockRepo.EnsureIndexes(ctx); err != nil {
		fmt.Println("error create index")
		return nil, err
	}

	filter := bson.M{
		"country":  country,
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

	fmt.Println("lock acquired successfully for", country, language)
	return doc, nil
}

// StreamData streams hotel rates to the frontend in real-time using Server-Sent Events (SSE)
//
// HOW THE CONNECTION WORKS:
// 1. The 'w http.ResponseWriter' parameter IS the HTTP connection established when client made the request
// 2. This connection stays OPEN for the duration of this function
// 3. When you write to 'w', you're writing to a BUFFER in the HTTP connection
// 4. Normally, the buffer only sends when the handler completes
// 5. The Flusher bypasses buffering and sends data IMMEDIATELY over the open connection
func (s *searchSvcs) streamData(ctx context.Context, w http.ResponseWriter, data map[string]any, hotels []models.Hotel) error {
	ids := make([]string, 0, len(hotels))
	for _, hotel := range hotels {
		ids = append(ids, hotel.Id)
	}

	ratesData := map[string]any{
		"hotelIds":         ids,
		"occupancies":      data["occupancies"],
		"currency":         "USD",
		"guestNationality": "US",
		"checkin":          data["checkin"],
		"checkout":         data["checkout"],
		"countryCode":      "USD",
		"stream":           true,
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

	rateCount := 0

	for event := range ch {
		if event.Error != nil {
			// Send error event to client
			errorData, _ := json.Marshal(map[string]any{
				"error": event.Error.Error(),
			})
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", errorData)
			flusher.Flush()
			return event.Error
		}

		if event.Done {
			// Send completion event
			fmt.Fprintf(w, "event: done\ndata: {\"total\": %d}\n\n", rateCount)
			flusher.Flush()
			break
		}

		var result map[string]any
		if err := json.Unmarshal(event.Data, &result); err != nil {
			// Send parse error
			errorData, _ := json.Marshal(map[string]any{
				"error": "failed to parse rate data",
			})
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", errorData)
			flusher.Flush()
			return err
		}

		// Stream each rate to the frontend
		rateJSON, err := json.Marshal(result)
		if err != nil {
			return err
		}

		// STEP 1: Write the SSE event to the HTTP response buffer
		// fmt.Fprintf writes to 'w', which is the HTTP connection's buffer
		// Format: "event: rate\ndata: {json}\n\n" is the SSE protocol format
		fmt.Fprintf(w, "event: rate\ndata: %s\n\n", rateJSON)

		// STEP 2: Flush immediately sends the buffered data over the TCP connection
		// Without Flush(), the data would sit in the buffer until the function returns
		// With Flush(), the client receives the data RIGHT NOW, in real-time
		flusher.Flush()

		// This creates a streaming effect:
		// Client ← [TCP Connection stays open] ← Server
		//   ↑                                       ↑
		//   Receives rate 1                     Flush sends rate 1
		//   Receives rate 2                     Flush sends rate 2
		//   Receives rate 3                     Flush sends rate 3
		//   ... and so on in real-time

		rateCount++
	}

	return nil
}

func (s *searchSvcs) releaseLock(ctx context.Context, lockId primitive.ObjectID) error {
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
