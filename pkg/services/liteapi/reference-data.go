package liteapi

import (
	"context"
	"encoding/json"
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
	GetCitiesByCountryCode(ctx context.Context, countryCode string) ([]models.City, error)
	GetCountries(ctx context.Context) ([]models.Country, error)
	GetCurrencies(ctx context.Context) ([]models.Currency, error)
	GetIatas(ctx context.Context) ([]models.Iata, error)
	GetHotelChains(ctx context.Context) ([]models.HotelChain, error)
	GetHotelTypes(ctx context.Context) ([]models.HotelType, error)
}

type referencedatasvcs struct {
	cityrepo       repo.CityRepo
	countryrepo    repo.CountryRepo
	currencyrepo   repo.CurrencyRepo
	iatarepo       repo.IataRepo
	hotelchainrepo repo.HotelChainRepo
	hoteltyperepo  repo.HotelTypeRepo
	liteApiSdk     *liteApiSdk.LiteApiSdk
}

func NewReferenceDataSvcs(i *do.Injector) (ReferenceDataSvcs, error) {
	return &referencedatasvcs{
		cityrepo:       do.MustInvoke[repo.CityRepo](i),
		countryrepo:    do.MustInvoke[repo.CountryRepo](i),
		currencyrepo:   do.MustInvoke[repo.CurrencyRepo](i),
		iatarepo:       do.MustInvoke[repo.IataRepo](i),
		hotelchainrepo: do.MustInvoke[repo.HotelChainRepo](i),
		hoteltyperepo:  do.MustInvoke[repo.HotelTypeRepo](i),
		liteApiSdk:     do.MustInvoke[*liteApiSdk.LiteApiSdk](i),
	}, nil
}

func (r *referencedatasvcs) GetCitiesByCountryCode(ctx context.Context, countryCode string) ([]models.City, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"country": countryCode, "expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}
	var result []models.City
	err := r.cityrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetCitiesByCountryCode(countryCode)
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.City `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var cities []models.City
		for _, city := range result.Data {
			cities = append(cities, models.City{
				City:      city.City,
				Country:   countryCode,
				ExpiresAt: time.Now().Add(20 * time.Second),
			})
		}

		go func(ctx context.Context) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{"country": countryCode})
			writeOps = append(writeOps, deleteOp)
			for _, city := range cities {
				insertOp := mongo.NewInsertOneModel().SetDocument(city)
				writeOps = append(writeOps, insertOp)
			}
			r.cityrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx))

		return cities, nil

	}
	return result, nil
}

func (r *referencedatasvcs) GetCountries(ctx context.Context) ([]models.Country, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Country
	err := r.countryrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetCountries()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.Country `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var countries []models.Country
		for _, country := range result.Data {
			country.ExpiresAt = time.Now().Add(20 * time.Second)
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

		return result.Data, nil
	}
	return result, nil
}

func (r *referencedatasvcs) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Currency
	err := r.currencyrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetCurrencies()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.Currency `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var currencies []models.Currency
		for _, currency := range result.Data {
			currency.ExpiresAt = time.Now().Add(20 * time.Second)
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

		return result.Data, nil
	}

	return result, nil
}

func (r *referencedatasvcs) GetIatas(ctx context.Context) ([]models.Iata, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Iata
	err := r.iatarepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetIataCodes()
		if err != nil {
			return nil, err
		}

		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.Iata `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		var iatas []models.Iata
		for _, iata := range result.Data {
			iata.ExpiresAt = time.Now().Add(20 * time.Second)
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

		return iatas, nil
	}

	return result, nil
}

func (r *referencedatasvcs) GetHotelChains(ctx context.Context) ([]models.HotelChain, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelChain
	err := r.hotelchainrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetHotelChains()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.HotelChain `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}
		var hotelChains []models.HotelChain
		for _, hotelChain := range result.Data {
			hotelChain.ExpiresAt = time.Now().Add(20 * time.Second)
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

		return hotelChains, nil
	}

	return result, nil
}

func (r *referencedatasvcs) GetHotelTypes(ctx context.Context) ([]models.HotelType, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelType
	err := r.hoteltyperepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		resp, err := r.liteApiSdk.GetHotelTypes()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}
		type aux struct {
			Data []models.HotelType `json:"data"`
		}
		var result aux
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}
		var hotelTypes []models.HotelType
		for _, hotelType := range result.Data {
			hotelType.ExpiresAt = time.Now().Add(20 * time.Second)
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

		return hotelTypes, nil
	}

	return result, nil
}
