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

type DataSvcs interface {
	GetHotels(ctx context.Context, query map[string]string) (*models.HotelListRes, error)
	GetHotelDetails(ctx context.Context, id, language, advancedAccessibilityOnly string) (*models.HotelDetailsData, error)
	GetCitiesByCountryCode(ctx context.Context, countryCode string) (*models.CityList, error)
	GetCountries(ctx context.Context) (*models.CountryList, error)
	GetCurrencies(ctx context.Context) (*models.CurrencyList, error)
	GetIatas(ctx context.Context) (*models.IataList, error)
	GetHotelChains(ctx context.Context) (*models.HotelChainList, error)
	GetHotelTypes(ctx context.Context) (*models.HotelTypeList, error)
	GetHotelFacilities(ctx context.Context) (*models.FacilityList, error)
}

type datasvcs struct {
	hotelrepo        repo.HotelRepo
	hoteldetailsrepo repo.HotelDetailsRepo
	cityrepo         repo.CityRepo
	countryrepo      repo.CountryRepo
	currencyrepo     repo.CurrencyRepo
	iatarepo         repo.IataRepo
	hotelchainrepo   repo.HotelChainRepo
	hoteltyperepo    repo.HotelTypeRepo
	facilityrepo     repo.FacilityRepo
	liteApiInitFunc  liteApiSdk.LiteApiInitFunc
}

func NewDataSvcs(i *do.Injector) (DataSvcs, error) {
	return &datasvcs{
		hotelrepo:        do.MustInvoke[repo.HotelRepo](i),
		hoteldetailsrepo: do.MustInvoke[repo.HotelDetailsRepo](i),
		cityrepo:         do.MustInvoke[repo.CityRepo](i),
		countryrepo:      do.MustInvoke[repo.CountryRepo](i),
		currencyrepo:     do.MustInvoke[repo.CurrencyRepo](i),
		iatarepo:         do.MustInvoke[repo.IataRepo](i),
		hotelchainrepo:   do.MustInvoke[repo.HotelChainRepo](i),
		hoteltyperepo:    do.MustInvoke[repo.HotelTypeRepo](i),
		facilityrepo:     do.MustInvoke[repo.FacilityRepo](i),
		liteApiInitFunc:  do.MustInvoke[liteApiSdk.LiteApiInitFunc](i),
	}, nil
}

const ExpireTime = 20 * time.Second

func (d *datasvcs) GetHotels(ctx context.Context, query map[string]string) (*models.HotelListRes, error) {
	lang := "en"
	if query["language"] != "" {
		lang = query["language"]
	}

	liteApiSdk, err := d.liteApiInitFunc(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := liteApiSdk.GetHotels(query, lang, 3, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if resp.Status == "failed" {
		return nil, helpers.LiteApiError(resp.Code, resp.Err)
	}

	var result models.HotelListRes
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	go func(ctx context.Context, hotels []models.Hotel) {
		writeOps := []mongo.WriteModel{}
		for _, hotel := range hotels {
			updateOp := mongo.NewUpdateOneModel().SetFilter(bson.M{"id": hotel.Id}).SetUpdate(bson.M{"$set": hotel}).SetUpsert(true)
			writeOps = append(writeOps, updateOp)
		}

		if _, err := d.hotelrepo.BulkWrite(ctx, writeOps); err != nil {
			fmt.Printf("error bulk writing hotels: %v", err)
		}

	}(context.WithoutCancel(ctx), result.Data)

	return &result, nil

}

func (d *datasvcs) GetHotelDetails(ctx context.Context, id, language, advancedAccessibilityOnly string) (*models.HotelDetailsData, error) {
	filter := bson.M{
		"id":        id,
		"expiresAt": bson.M{"$gt": time.Now()},
	}

	//todo : check for language code

	result, err := d.hoteldetailsrepo.GetByFilter(ctx, filter)
	if err != nil {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		lang := "en"
		if language != "" {
			lang = language
		}

		resp, err := liteApiSdk.GetHotelDetails(id, lang, advancedAccessibilityOnly)
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.HotelDetailsData
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		result.Data["expiresAt"] = time.Now().Add(ExpireTime)

		go func(ctx context.Context, data models.HotelDetails) {
			writeOps := []mongo.WriteModel{}
			updateOp := mongo.NewUpdateOneModel().SetFilter(bson.M{"id": data["id"]}).SetUpdate(bson.M{"$set": data}).SetUpsert(true)
			writeOps = append(writeOps, updateOp)

			if _, err := d.hoteldetailsrepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing hotel details: %v", err)
			}
		}(context.WithoutCancel(ctx), result.Data)

		return &result, nil
	}

	return &models.HotelDetailsData{
		Data: *result,
	}, nil

}

func (d *datasvcs) GetCitiesByCountryCode(ctx context.Context, countryCode string) (*models.CityList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"country": countryCode, "expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}
	var result []models.City
	err := d.cityrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		fmt.Printf("fetch form liteapi\n")
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			d.cityrepo.BulkWrite(ctx, writeOps)
		}(context.WithoutCancel(ctx), cities)

		return &result, nil

	}

	fmt.Printf("fetch from db\n")
	return &models.CityList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetCountries(ctx context.Context) (*models.CountryList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Country
	err := d.countryrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		fmt.Printf("fetch form liteapi\n")
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			if _, err := d.countryrepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing countries: %v", err)
			}
		}(context.WithoutCancel(ctx), countries)

		return &result, nil
	}
	return &models.CountryList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetCurrencies(ctx context.Context) (*models.CurrencyList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Currency
	err := d.currencyrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			if _, err := d.currencyrepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing currencies: %v", err)
			}
		}(context.WithoutCancel(ctx), currencies)

		return &result, nil
	}

	return &models.CurrencyList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetIatas(ctx context.Context) (*models.IataList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Iata
	err := d.iatarepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			if _, err := d.iatarepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing iatas: %v", err)
			}
		}(context.WithoutCancel(ctx), iatas)

		return &result, nil
	}

	return &models.IataList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetHotelChains(ctx context.Context) (*models.HotelChainList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelChain
	err := d.hotelchainrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			if _, err := d.hotelchainrepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing hotel chains: %v", err)
			}
		}(context.WithoutCancel(ctx), hotelChains)

		return &result, nil
	}

	return &models.HotelChainList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetHotelTypes(ctx context.Context) (*models.HotelTypeList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.HotelType
	err := d.hoteltyperepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
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
			if _, err := d.hoteltyperepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing hotel types: %v", err)
			}
		}(context.WithoutCancel(ctx), hotelTypes)

		return &result, nil
	}

	return &models.HotelTypeList{
		Data: result,
	}, nil
}

func (d *datasvcs) GetHotelFacilities(ctx context.Context) (*models.FacilityList, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"expiresAt": bson.M{"$gt": time.Now().UTC()}}},
	}

	var result []models.Facility
	err := d.facilityrepo.Aggregate(ctx, pipeline, func(cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil || len(result) == 0 {
		liteApiSdk, err := d.liteApiInitFunc(ctx)
		if err != nil {
			return nil, err
		}
		resp, err := liteApiSdk.GetHotelFacilities()
		if err != nil {
			return nil, err
		}
		if resp.Status == "failed" {
			return nil, helpers.LiteApiError(resp.Code, resp.Err)
		}

		var result models.FacilityList
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}
		var facilities []models.Facility
		for _, facility := range result.Data {
			facility.ExpiresAt = time.Now().Add(ExpireTime)
			facilities = append(facilities, facility)
		}

		go func(ctx context.Context, facilities []models.Facility) {
			writeOps := []mongo.WriteModel{}
			deleteOp := mongo.NewDeleteManyModel().SetFilter(bson.M{})
			writeOps = append(writeOps, deleteOp)
			for _, facility := range facilities {
				insertOp := mongo.NewInsertOneModel().SetDocument(facility)
				writeOps = append(writeOps, insertOp)
			}
			if _, err := d.hoteltyperepo.BulkWrite(ctx, writeOps); err != nil {
				fmt.Printf("error bulk writing facilities: %v", err)
			}
		}(context.WithoutCancel(ctx), facilities)

		return &result, nil
	}

	return &models.FacilityList{
		Data: result,
	}, nil
}
