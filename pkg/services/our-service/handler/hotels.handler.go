package handler

import (
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourService "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type HotelsHandler struct {
	hotelssvcs ourService.HotelsSvcs
}

func NewHotelsHandler(i *do.Injector, r *chi.Mux) {
	h := &HotelsHandler{
		hotelssvcs: do.MustInvoke[ourService.HotelsSvcs](i),
	}

	r.Route("/hotels", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/", helpers.Make(h.GetAll))
		r.Get("/all", helpers.Make(h.GetAllHotels))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (l *HotelsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.hotelssvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) GetAllHotels(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Create a new filter
	hotelFilter := filter.HotelsFilter{}

	// Parse JSON filter parameter if present
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		if err := json.Unmarshal([]byte(filterParam), &hotelFilter); err != nil {

			return fmt.Errorf("invalid filter format: %v", err)
		}
	}

	// Parse direct query parameters into filter
	if searchWord := r.URL.Query().Get("searchWord"); searchWord != "" {
		hotelFilter.SearchWord = searchWord
	}

	if location := r.URL.Query().Get("location"); location != "" {
		hotelFilter.Locations = []string{location}
	}

	// Parse hotel type (both single and array)
	if hotelType := r.URL.Query().Get("hotelType"); hotelType != "" {
		hotelFilter.HotelType = hotelType
	}

	// Parse hotel types as array
	if hotelTypes := r.URL.Query().Get("hotelTypes"); hotelTypes != "" {
		hotelFilter.HotelTypes = strings.Split(hotelTypes, ",")
	}

	// Parse room amenities
	if roomAmenities := r.URL.Query().Get("roomAmenities"); roomAmenities != "" {
		hotelFilter.RoomAmenities = strings.Split(roomAmenities, ",")
	}

	// Parse nearby attractions
	if nearbyAttractions := r.URL.Query().Get("nearbyAttractions"); nearbyAttractions != "" {
		hotelFilter.NearbyAttractions = strings.Split(nearbyAttractions, ",")
	}

	// Parse ratings
	if ratings := r.URL.Query().Get("ratings"); ratings != "" {
		if rating, err := strconv.ParseFloat(ratings, 64); err == nil {
			hotelFilter.Ratings = rating
		}
	}

	// Parse price range parameters
	if priceFrom := r.URL.Query().Get("priceFrom"); priceFrom != "" {
		if from, err := strconv.Atoi(priceFrom); err == nil {
			hotelFilter.PriceRange.From = from
		}
	}

	if priceTo := r.URL.Query().Get("priceTo"); priceTo != "" {
		if to, err := strconv.Atoi(priceTo); err == nil {
			hotelFilter.PriceRange.To = to
		}
	}

	// Parse distance from city center
	if greaterThan := r.URL.Query().Get("distanceGreaterThan"); greaterThan != "" {
		if gt, err := strconv.Atoi(greaterThan); err == nil {
			hotelFilter.DistanceFromCityCenter.GreaterThan = gt
		}
	}

	if lessThan := r.URL.Query().Get("distanceLessThan"); lessThan != "" {
		if lt, err := strconv.Atoi(lessThan); err == nil {
			hotelFilter.DistanceFromCityCenter.LessThan = lt
		}
	}

	// Parse isDisplayInPerfectStay if present
	if perfectStay := r.URL.Query().Get("isDisplayInPerfectStay"); perfectStay != "" {
		if val, err := strconv.ParseBool(perfectStay); err == nil {
			hotelFilter.IsDisplayInPerfectStay = &val
		}
	}

	// Validate price range
	if err := hotelFilter.PriceRange.Validate(); err != nil {

		return fmt.Errorf("invalid price range: %v", err)
	}

	result, err := l.hotelssvcs.GetAllHotels(ctx, hotelFilter)
	if err != nil {

		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Create a new filter
	hotelFilter := filter.HotelsFilter{}

	// Parse JSON filter parameter if present
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		if err := json.Unmarshal([]byte(filterParam), &hotelFilter); err != nil {

			return fmt.Errorf("invalid filter format: %v", err)
		}
	}

	// Parse direct query parameters into filter
	if searchWord := r.URL.Query().Get("searchWord"); searchWord != "" {
		hotelFilter.SearchWord = searchWord
	}

	if location := r.URL.Query().Get("location"); location != "" {
		hotelFilter.Locations = []string{location}
	}

	// Parse hotel type (both single and array)
	if hotelType := r.URL.Query().Get("hotelType"); hotelType != "" {
		hotelFilter.HotelType = hotelType
	}

	// Parse hotel types as array
	if hotelTypes := r.URL.Query().Get("hotelTypes"); hotelTypes != "" {
		hotelFilter.HotelTypes = strings.Split(hotelTypes, ",")
	}

	// Parse room amenities
	if roomAmenities := r.URL.Query().Get("roomAmenities"); roomAmenities != "" {
		hotelFilter.RoomAmenities = strings.Split(roomAmenities, ",")
	}

	// Parse nearby attractions
	if nearbyAttractions := r.URL.Query().Get("nearbyAttractions"); nearbyAttractions != "" {
		hotelFilter.NearbyAttractions = strings.Split(nearbyAttractions, ",")
	}

	// Parse ratings
	if ratings := r.URL.Query().Get("ratings"); ratings != "" {
		if rating, err := strconv.ParseFloat(ratings, 64); err == nil {
			hotelFilter.Ratings = rating
		}
	}

	// Parse price range parameters
	if priceFrom := r.URL.Query().Get("priceFrom"); priceFrom != "" {
		if from, err := strconv.Atoi(priceFrom); err == nil {
			hotelFilter.PriceRange.From = from
		}
	}

	if priceTo := r.URL.Query().Get("priceTo"); priceTo != "" {
		if to, err := strconv.Atoi(priceTo); err == nil {
			hotelFilter.PriceRange.To = to
		}
	}

	// Parse distance from city center
	if greaterThan := r.URL.Query().Get("distanceGreaterThan"); greaterThan != "" {
		if gt, err := strconv.Atoi(greaterThan); err == nil {
			hotelFilter.DistanceFromCityCenter.GreaterThan = gt
		}
	}

	if lessThan := r.URL.Query().Get("distanceLessThan"); lessThan != "" {
		if lt, err := strconv.Atoi(lessThan); err == nil {
			hotelFilter.DistanceFromCityCenter.LessThan = lt
		}
	}

	// Parse isDisplayInPerfectStay if present
	if perfectStay := r.URL.Query().Get("isDisplayInPerfectStay"); perfectStay != "" {
		if val, err := strconv.ParseBool(perfectStay); err == nil {
			hotelFilter.IsDisplayInPerfectStay = &val
		}
	}

	// Validate price range
	if err := hotelFilter.PriceRange.Validate(); err != nil {

		return fmt.Errorf("invalid price range: %v", err)
	}

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int64((skip / limit) + 1)
	perPage := int64(limit)

	result, err := l.hotelssvcs.GetAll(ctx, hotelFilter, page, perPage)
	if err != nil {

		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := l.hotelssvcs.Add(ctx, &data)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.hotelssvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	id := chi.URLParam(r, "id")
	result, err := l.hotelssvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}
