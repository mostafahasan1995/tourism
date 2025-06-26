package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/statistics"
	"larsa-tourism-microservices/pkg/services/statistics/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type StatisticsHandler struct {
	statisticsSvcs statistics.StatisticsSvcs
}

func NewStatisticsHandler(i *do.Injector, r *chi.Mux) {
	h := &StatisticsHandler{
		statisticsSvcs: do.MustInvoke[statistics.StatisticsSvcs](i),
	}

	r.Route("/statistics", func(r chi.Router) {

		r.Get("/", helpers.Make(h.GetStatistics))

		r.Get("/programs", helpers.Make(h.GetProgramsCount))
		r.Get("/trips", helpers.Make(h.GetTripsCount))
		r.Get("/countries", helpers.Make(h.GetCountriesCount))
		r.Get("/hotels", helpers.Make(h.GetHotelsCount))
		r.Get("/happy-travelers", helpers.Make(h.GetHappyTravelersCount))

		// Manual statistics endpoints
		r.Get("/manual", helpers.Make(h.GetManualStatistics))
		r.Put("/manual", helpers.Make(h.UpdateManualStatistics))
		r.Post("/manual/auto-calculate", helpers.Make(h.ToggleAutoCalculate))
	})
}

func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.statisticsSvcs.GetStatistics(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *StatisticsHandler) GetProgramsCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	count, err := h.statisticsSvcs.GetProgramsCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int64{"count": count})
}

func (h *StatisticsHandler) GetTripsCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	count, err := h.statisticsSvcs.GetTripsCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int64{"count": count})
}

func (h *StatisticsHandler) GetCountriesCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	count, err := h.statisticsSvcs.GetCountriesCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int64{"count": count})
}

func (h *StatisticsHandler) GetHotelsCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	count, err := h.statisticsSvcs.GetHotelsCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int64{"count": count})
}

func (h *StatisticsHandler) GetHappyTravelersCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	count, err := h.statisticsSvcs.GetHappyTravelersCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]int64{"count": count})
}

// GetManualStatistics retrieves the current manual statistics configuration
func (h *StatisticsHandler) GetManualStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.statisticsSvcs.GetManualStatistics(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

// UpdateManualStatistics updates manual statistics values
func (h *StatisticsHandler) UpdateManualStatistics(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var req models.ManualStatisticsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	result, err := h.statisticsSvcs.UpdateManualStatistics(ctx, &req)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

// ToggleAutoCalculate toggles the auto-calculate mode
func (h *StatisticsHandler) ToggleAutoCalculate(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var req models.AutoCalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return helpers.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	result, err := h.statisticsSvcs.ToggleAutoCalculate(ctx, req.AutoCalculate)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
