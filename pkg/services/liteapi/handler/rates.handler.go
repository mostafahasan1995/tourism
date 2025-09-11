package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type RatesHandler struct {
	ratessvcs liteapi.RatesSvcs
}

func NewRatesHandler(i *do.Injector, r *chi.Mux) {
	h := &RatesHandler{
		ratessvcs: do.MustInvoke[liteapi.RatesSvcs](i),
	}

	r.Route("/liteapi/hotels", func(r chi.Router) {
		r.Post("/rates", helpers.Make(h.GetHotelsRates))
		r.Post("/min-rates", helpers.Make(h.GetMinRates))
	})

	r.Route("/liteapi/rates", func(r chi.Router) {
		r.Post("/prebook", helpers.Make(h.Prebook))
		r.Post("/book", helpers.Make(h.Book))
	})

}

func (h *RatesHandler) GetHotelsRates(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data map[string]any
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.ratessvcs.GetHotelsRates(ctx, data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) GetMinRates(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data map[string]any
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.ratessvcs.GetMinRates(ctx, data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) Prebook(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data map[string]any
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.ratessvcs.PreBook(ctx, data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) Book(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data map[string]any
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.ratessvcs.Book(ctx, data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
