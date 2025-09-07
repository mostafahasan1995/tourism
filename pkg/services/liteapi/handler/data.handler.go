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
	hotelsvcs         liteapi.HotelSvcs
	referencedatasvcs liteapi.ReferenceDataSvcs
}

func NewDataHandler(i *do.Injector, r *chi.Mux) {
	h := &DataHandler{
		hotelsvcs:         do.MustInvoke[liteapi.HotelSvcs](i),
		referencedatasvcs: do.MustInvoke[liteapi.ReferenceDataSvcs](i),
	}

	r.Route("/liteapi/data", func(r chi.Router) {
		r.Get("/hotels", helpers.Make(h.GetHotels))
		r.Get("/hotels/streaming", helpers.Make(h.TestStreaming))
		r.Get("/hotels/streaming2", helpers.Make(h.TestStreaming2))
		//
		r.Get("/cities", helpers.Make(h.GetCities))
		r.Get("/countries", helpers.Make(h.GetCountries))
		r.Get("/currencies", helpers.Make(h.GetCurrencies))
		r.Get("/iatas", helpers.Make(h.GetIatas))
		r.Get("/hotel-chains", helpers.Make(h.GetHotelChains))
		r.Get("/hotel-types", helpers.Make(h.GetHotelTypes))
	})

}

func (h *DataHandler) GetHotels(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	query := make(map[string]string)
	for key, value := range params {
		query[key] = value[0]
	}

	result, err := h.hotelsvcs.SearchHotels(ctx, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) TestStreaming(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()
	query := make(map[string]string)
	for key, value := range params {
		query[key] = value[0]
	}

	err := h.hotelsvcs.TestStreaming(ctx, w, query)
	if err != nil {
		return err
	}
	return nil
}

func (h *DataHandler) TestStreaming2(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()
	query := make(map[string]string)
	for key, value := range params {
		query[key] = value[0]
	}

	err := h.hotelsvcs.TestStreaming2(ctx, w, query)
	if err != nil {
		return err
	}
	return nil
}

func (h *DataHandler) GetCities(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()
	countryCode := params.Get("countryCode")

	result, err := h.referencedatasvcs.GetCitiesByCountryCode(ctx, countryCode)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetCountries(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.referencedatasvcs.GetCountries(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetCurrencies(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.referencedatasvcs.GetCurrencies(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetIatas(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.referencedatasvcs.GetIatas(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelChains(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.referencedatasvcs.GetHotelChains(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *DataHandler) GetHotelTypes(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.referencedatasvcs.GetHotelTypes(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
