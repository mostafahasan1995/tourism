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

	if dataFetched != nil && dataFetched.IsFullyFetched() {
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

func (s *searchV2Svcs) pollFromDbV2(ctx context.Context, placeId, language string, callback func(hotels []models.Hotel) error) error {
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

		callback(result.Data)

		numOfHotelFetched += len(result.Data)

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

	placeId, ok := query["placeId"].(string)
	if !ok {
		return nil, errors.New("placeId is required")
	}

	go s.fetchHotelsInBackground(context.WithoutCancel(ctx), placeId, language)

	s.streamDataV2(context.WithoutCancel(ctx), w, placeId, language, query)

	// if err := s.pollFromDb(ctx, w, query, placeId, language); err != nil {
	// 	return nil, err
	// }

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

func (s *searchV2Svcs) streamDataV2(ctx context.Context, w http.ResponseWriter, placeId, language string, data map[string]any) error {
	dataFetched, err := s.getDataFetchedInfo(ctx, placeId, language)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("error getting data fetched info: %w -%s-%s", err, placeId, language)
	}

	ratesData := map[string]any{
		"placeId":          data["placeId"],
		"occupancies":      data["occupancies"],
		"currency":         "USD",
		"guestNationality": data["guestNationality"],
		"checkin":          data["checkin"],
		"checkout":         data["checkout"],
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

	// Set headers for Server-Sent Events (SSE) before writing any data
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, x-client, x-access-token, x-service, x-expire, x-service-token, x-user-id")

	// Get the flusher to flush data immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming not supported")
	}

	if dataFetched != nil && dataFetched.IsFullyFetched() {
		hotelIds := make([]string, 0, len(result.Data))
		for _, rate := range result.Data {
			if hotelId, ok := rate["hotelId"].(string); ok {
				hotelIds = append(hotelIds, hotelId)
			}
		}
		hotels, err := s.dataSvcs.GetHotelsByIds(ctx, hotelIds)
		if err != nil {
			return fmt.Errorf("error getting hotels by ids")
		}
		m := make(map[string]models.Hotel, len(hotels))
		for _, hotel := range hotels {
			m[hotel.Id] = hotel
		}

		for _, rate := range result.Data {
			if hotelId, ok := rate["hotelId"].(string); ok {
				if hotel, exists := m[hotelId]; exists {
					rate["hotel"] = hotel
					rateJSON, err := json.Marshal(rate)
					if err != nil {
						return err
					}
					fmt.Fprintf(w, "data: %s\n\n", rateJSON)
					flusher.Flush()
				}
			}
		}

	} else {

		s.pollFromDbV2(ctx, placeId, language, func(hotels []models.Hotel) error {

			m := make(map[string]models.Hotel, len(hotels))
			for _, hotel := range hotels {
				m[hotel.Id] = hotel
			}

			for _, rate := range result.Data {
				if hotelId, ok := rate["hotelId"].(string); ok {
					if hotel, exists := m[hotelId]; exists {
						rate["hotel"] = hotel
						rateJSON, err := json.Marshal(rate)
						if err != nil {
							return err
						}
						fmt.Fprintf(w, "data: %s\n\n", rateJSON)
						flusher.Flush()
					}
				}

			}

			return nil
		})

	}

	// Send [DONE] signal to indicate completion
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()

	return nil
}
