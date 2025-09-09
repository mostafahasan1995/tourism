package liteapi

import (
	"context"
	"encoding/json"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi/models"
	"larsa-tourism-microservices/pkg/services/liteapi/repo"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReferenceDataSvcs interface {
	GetCitiesByCountryCode(ctx context.Context, countryCode string) (*models.CityList, error)
	GetCountries(ctx context.Context) (*models.CountryList, error)
	GetCurrencies(ctx context.Context) (*models.CurrencyList, error)
	GetIatas(ctx context.Context) (*models.IataList, error)
	GetHotelChains(ctx context.Context) (*models.HotelChainList, error)
	GetHotelTypes(ctx context.Context) (*models.HotelTypeList, error)
}

type referencedatasvcs struct {
	cityrepo        repo.CityRepo
	countryrepo     repo.CountryRepo
	currencyrepo    repo.CurrencyRepo
	iatarepo        repo.IataRepo
	hotelchainrepo  repo.HotelChainRepo
	hoteltyperepo   repo.HotelTypeRepo
	liteApiInitFunc liteApiSdk.LiteApiInitFunc
}

func NewReferenceDataSvcs(i *do.Injector) (ReferenceDataSvcs, error) {
	return &referencedatasvcs{
		cityrepo:        do.MustInvoke[repo.CityRepo](i),
		countryrepo:     do.MustInvoke[repo.CountryRepo](i),
		currencyrepo:    do.MustInvoke[repo.CurrencyRepo](i),
		iatarepo:        do.MustInvoke[repo.IataRepo](i),
		hotelchainrepo:  do.MustInvoke[repo.HotelChainRepo](i),
		hoteltyperepo:   do.MustInvoke[repo.HotelTypeRepo](i),
		liteApiInitFunc: do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
	}, nil
}

const ExpireTime = 20 * time.Second

func (r *referencedatasvcs) GetCitiesByCountryCode(ctx context.Context, countryCode string) (*models.CityList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"country": countryCode, "expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}
	var result []models.City
	err := r.cityrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		fmt.Printf("fetch form liteapi\n")
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetCitiesByCountryCode(countryCode)
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.CityList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var cities []models.City
		for _, city := range result.Data {
			cities = append(cities, models.City{
				City:      city.City,
				Country:   countryCode,
				ExpiresAt: time.Now().Add(ExpireTime),
			})
		}

		go func(ctx context.Context, cities []models.City) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{"country": countryCode})
			writeOps = append(writeOps, deleteOp)
			for _, city := range cities {
				insertOp := mongo.NewInsertOneModel().SetDocument(city)
				writeOps = append(writeOps, insertOp)
			}
			r.cityrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), cities)

		return &result, nil

	}

	fmt.Printf("fetch from db\n")
	return &models.CityList{
		Data: result,
	}, nil
}

func (r *referencedatasvcs) GetCountries(ctx context.Context) (*models.CountryList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Country
	err := r.countryrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		fmt.Printf("fetch form liteapi\n")
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetCountries()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.CountryList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var countries []models.Country
		for _, country := range result.Data {
			country.ExpiresAt = time.Now().Add(ExpireTime)
			countries = append(countries, country)
		}

		go func(ctx context.Context, countries []models.Country) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, country := range countries {
				insertOp := mongo.NewInsertOneModel().SetDocument(country)
				writeOps = append(writeOps, insertOp)
			}
			r.countryrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), countries)

		return &result, nil
	}
	return &models.CountryList{
		Data: result,
	}, nil
}

func (r *referencedatasvcs) GetCurrencies(ctx context.Context) (*models.CurrencyList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Currency
	err := r.currencyrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetCurrencies()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.CurrencyList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var currencies []models.Currency
		for _, currency := range result.Data {
			currency.ExpiresAt = time.Now().Add(ExpireTime)
			currencies = append(currencies, currency)
		}

		go func(ctx context.Context, currencies []models.Currency) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, currency := range currencies {
				insertOp := mongo.NewInsertOneModel().SetDocument(currency)
				writeOps = append(writeOps, insertOp)
			}
			r.currencyrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), currencies)

		return &result, nil
	}

	return &models.CurrencyList{
		Data: result,
	}, nil
}

func (r *referencedatasvcs) GetIatas(ctx context.Context) (*models.IataList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Iata
	err := r.iatarepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetIataCodes()
		if err != nil {
			return nil, err
		}

		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.IataList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var iatas []models.Iata
		for _, iata := range result.Data {
			iata.ExpiresAt = time.Now().Add(ExpireTime)
			iatas = append(iatas, iata)
		}

		go func(ctx context.Context, iatas []models.Iata) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, iata := range iatas {
				insertOp := mongo.NewInsertOneModel().SetDocument(iata)
				writeOps = append(writeOps, insertOp)
			}
			r.iatarepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), iatas)

		return &result, nil
	}

	return &models.IataList{
		Data: result,
	}, nil
}

func (r *referencedatasvcs) GetHotelChains(ctx context.Context) (*models.HotelChainList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelChain
	err := r.hotelchainrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetHotelChains()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.HotelChainList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var hotelChains []models.HotelChain
		for _, hotelChain := range result.Data {
			hotelChain.ExpiresAt = time.Now().Add(ExpireTime)
			hotelChains = append(hotelChains, hotelChain)
		}

		go func(ctx context.Context, hotelChains []models.HotelChain) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, hotelChain := range hotelChains {
				insertOp := mongo.NewInsertOneModel().SetDocument(hotelChain)
				writeOps = append(writeOps, insertOp)
			}
			r.hotelchainrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), hotelChains)

		return &result, nil
	}

	return &models.HotelChainList{
		Data: result,
	}, nil
}

func (r *referencedatasvcs) GetHotelTypes(ctx context.Context) (*models.HotelTypeList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelType
	err := r.hoteltyperepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := r.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetHotelTypes()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.HotelTypeList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}
		var hotelTypes []models.HotelType
		for _, hotelType := range result.Data {
			hotelType.ExpiresAt = time.Now().Add(ExpireTime)
			hotelTypes = append(hotelTypes, hotelType)
		}

		go func(ctx context.Context, hotelTypes []models.HotelType) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, hotelType := range hotelTypes {
				insertOp := mongo.NewInsertOneModel().SetDocument(hotelType)
				writeOps = append(writeOps, insertOp)
			}
			r.hoteltyperepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), hotelTypes)

		return &result, nil
	}

	return &models.HotelTypeList{
		Data: result,
	}, nil
}
