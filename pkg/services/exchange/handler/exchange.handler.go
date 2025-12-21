package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	exchangeService "larsa-tourism-microservices/pkg/services/exchange"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type ExchangeHandler struct {
	exchangeSvc exchangeService.ExchangeService
}

func NewExchangeHandler(i *do.Injector, r *chi.Mux) {
	h := &ExchangeHandler{
		exchangeSvc: do.MustInvoke[exchangeService.ExchangeService](i),
	}

	r.Route("/exchange", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/rates", helpers.Make(h.GetRates))
		r.With(middleware.Auth("authenticate")).Post("/convert", helpers.Make(h.Convert))
	})
}

// ConvertRequest represents the request body for currency conversion
type ConvertRequest struct {
	Amount       float64 `json:"amount"`
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
}

// ConvertResponse represents the response for currency conversion
type ConvertResponse struct {
	Amount          float64 `json:"amount"`
	FromCurrency    string  `json:"fromCurrency"`
	ToCurrency      string  `json:"toCurrency"`
	ConvertedAmount float64 `json:"convertedAmount"`
}

// RatesResponse represents the response for exchange rates
type RatesResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
}

func (h *ExchangeHandler) GetRates(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	rates, err := h.exchangeSvc.GetRates()
	if err != nil {
		return err
	}

	response := RatesResponse{
		Rates: rates,
		Base:  "USD",
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExchangeHandler) Convert(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); err != nil {
		return err
	}

	convertedAmount, err := h.exchangeSvc.Convert(req.Amount, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return err
	}

	response := ConvertResponse{
		Amount:          req.Amount,
		FromCurrency:    req.FromCurrency,
		ToCurrency:      req.ToCurrency,
		ConvertedAmount: convertedAmount,
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}
