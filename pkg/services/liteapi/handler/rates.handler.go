package handler

import (
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strings"
	"time"

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
		//
		r.Post("/rates/stream", helpers.Make(h.GetHotelsRatesStream))
	})

	r.Route("/liteapi/rates", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Post("/prebook", helpers.Make(h.Prebook))
		r.With(middleware.Auth("authenticate")).Post("/book", helpers.Make(h.Book))

	})

	r.Route("/liteapi/bookings", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Put("/{bookingId}", helpers.Make(h.CancelBooking))
	})

	r.Route("/bookings", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/prebook/me", helpers.Make(h.MyPrebooks))
		r.With(middleware.Auth("authenticate")).Get("/book/me", helpers.Make(h.MyBookings))
		// User endpoint - gets bookings by email from token
		r.With(middleware.Auth("authenticate")).Get("/by-email", helpers.Make(h.GetBookingsByEmail))
		// Admin endpoint - requires admin capability
		r.With(middleware.Auth("authenticate"), middleware.CapabilityCheck("admin")).Get("/admin", helpers.Make(h.GetBookingsAdmin))
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

func (h *RatesHandler) MyPrebooks(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.ratessvcs.MyPrebooks(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) MyBookings(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.ratessvcs.MyBookings(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) CancelBooking(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	bookingId := chi.URLParam(r, "bookingId")

	result, err := h.ratessvcs.CancelBooking(ctx, bookingId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// stream
func (h *RatesHandler) GetHotelsRatesStream(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data map[string]any
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.ratessvcs.GetFullRatesStream(ctx, data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *RatesHandler) GetBookingsByEmail(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	var fromDate, toDate *time.Time

	// Parse fromDate if provided
	if fromDateStr := params.Get("fromDate"); fromDateStr != "" {
		// Trim quotes and whitespace
		fromDateStr = strings.Trim(fromDateStr, `"' `)
		parsed, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return fmt.Errorf("invalid fromDate format, use YYYY-MM-DD: %w", err)
		}
		fromDate = &parsed
	}

	// Parse toDate if provided
	if toDateStr := params.Get("toDate"); toDateStr != "" {
		// Trim quotes and whitespace
		toDateStr = strings.Trim(toDateStr, `"' `)
		parsed, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return fmt.Errorf("invalid toDate format, use YYYY-MM-DD: %w", err)
		}
		toDate = &parsed
	}

	result, err := h.ratessvcs.GetBookingsByEmail(ctx, fromDate, toDate)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

func (h *RatesHandler) GetBookingsAdmin(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	var email *string
	var fromDate, toDate *time.Time

	// Parse email if provided
	if emailStr := params.Get("email"); emailStr != "" {
		// Trim quotes and whitespace
		emailStr = strings.Trim(emailStr, `"' `)
		email = &emailStr
	}

	// Parse fromDate if provided
	if fromDateStr := params.Get("fromDate"); fromDateStr != "" {
		// Trim quotes and whitespace
		fromDateStr = strings.Trim(fromDateStr, `"' `)
		parsed, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return fmt.Errorf("invalid fromDate format, use YYYY-MM-DD: %w", err)
		}
		fromDate = &parsed
	}

	// Parse toDate if provided
	if toDateStr := params.Get("toDate"); toDateStr != "" {
		// Trim quotes and whitespace
		toDateStr = strings.Trim(toDateStr, `"' `)
		parsed, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return fmt.Errorf("invalid toDate format, use YYYY-MM-DD: %w", err)
		}
		toDate = &parsed
	}

	result, err := h.ratessvcs.GetBookingsAdmin(ctx, email, fromDate, toDate)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]interface{}{
		"data": result,
	})
}
