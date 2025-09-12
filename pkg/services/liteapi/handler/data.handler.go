package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type DataHandler struct {
	datasvcs liteapi.DataSvcs
}

func NewDataHandler(i *do.Injector, r *chi.Mux) {
	h := &DataHandler{
		datasvcs: do.MustInvoke[liteapi.DataSvcs](i),
	}

	r.Route("/liteapi/data", func(r chi.Router) {
		r.Get("/hotels", helpers.Make(h.GetHotels))
		r.Get("/hotel", helpers.Make(h.GetHotelDetails))
		//
		r.Get("/cities", helpers.Make(h.GetCities))
		r.Get("/countries", helpers.Make(h.GetCountries))
		r.Get("/currencies", helpers.Make(h.GetCurrencies))
		r.Get("/iataCodes", helpers.Make(h.GetIatas))
		r.Get("/chains", helpers.Make(h.GetHotelChains))
		r.Get("/hotelTypes", helpers.Make(h.GetHotelTypes))
		r.Get("/facilities", helpers.Make(h.GetHotelFacilities))

	})

}

func (h *DataHandler) GetHotels(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	query := make(map[string]string)
	for key, value := range params {
		query[key] = value[0]
	}

	result, err := h.datasvcs.GetHotels(ctx, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelDetails(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	id := params.Get("hotelId")
	language := params.Get("language")
	advancedAccessibilityOnly := params.Get("advancedAccessibilityOnly")

	result, err := h.datasvcs.GetHotelDetails(ctx, id, language, advancedAccessibilityOnly)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetCities(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()
	countryCode := params.Get("countryCode")

	result, err := h.datasvcs.GetCitiesByCountryCode(ctx, countryCode)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetCountries(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetCountries(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetCurrencies(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetCurrencies(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetIatas(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetIatas(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelChains(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetHotelChains(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelTypes(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetHotelTypes(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelFacilities(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.datasvcs.GetHotelFacilities(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
