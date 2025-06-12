package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/statistics"
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
