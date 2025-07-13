package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type TravelExperHandler struct {
	travelexpersvcs interactions.TravelExperSvcs
}

func NewTravelExperHandler(i *do.Injector, r *chi.Mux) {
	h := &TravelExperHandler{
		travelexpersvcs: do.MustInvoke[interactions.TravelExperSvcs](i),
	}

	r.Route("/travel-exper", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Get("/traveler-stories/{id}", helpers.Make(h.GetTravelerStory))
		r.With(middleware.OptionalAuth()).Get("/traveler-stories/", helpers.Make(h.GetTravelerStories))
		r.With(middleware.Auth("authenticate")).Post("/traveler-stories/", helpers.Make(h.AddTravelerStory))
		r.With(middleware.Auth("authenticate")).Patch("/traveler-stories/{id}/status", helpers.Make(h.SetTravelerStoryStatus))
		r.With(middleware.Auth("authenticate")).Patch("/traveler-stories/{id}/feedback", helpers.Make(h.SendFeedback))
		r.With(middleware.Auth("authenticate")).Put("/traveler-stories/{id}", helpers.Make(h.UpdateTravelerStory))
		r.With(middleware.Auth("authenticate")).Delete("/traveler-stories/{id}", helpers.Make(h.DeleteTravelerStory))
		r.With(middleware.Auth("authenticate")).Patch("/traveler-stories/{id}/restore", helpers.Make(h.RestoreTravelerStory))
		//client
		r.Get("/client-stories/{id}", helpers.Make(h.GetClientStory))
		r.With(middleware.OptionalAuth()).Get("/client-stories/", helpers.Make(h.GetClientStories))
		r.With(middleware.Auth("authenticate")).Post("/client-stories/", helpers.Make(h.AddClientStory))
		r.With(middleware.Auth("authenticate")).Put("/client-stories/{id}", helpers.Make(h.UpdateClientStory))
		r.With(middleware.Auth("authenticate")).Delete("/client-stories/{id}", helpers.Make(h.DeleteClientStory))
	})

	r.Route("/travel-exper/v2", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Post("/traveler-stories/", helpers.Make(h.GetTravelerStoriesV2))
		r.With(middleware.OptionalAuth()).Post("/client-stories/", helpers.Make(h.GetClientStoriesV2))
	})
}

func (h *TravelExperHandler) GetTravelerStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // story id

	result, err := h.travelexpersvcs.GetTravelerStory(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) GetTravelerStories(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	query := r.URL.Query().Get("query")

	result, err := h.travelexpersvcs.GetTravelerStories(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) AddTravelerStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TravelerStoryDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.AddTravelerStory(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) SetTravelerStoryStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	var data models.TravelerStoryStatusDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.SetTravelerStoryStatus(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) SendFeedback(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	var data models.TravelerStoryFeedback
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.SendFeedback(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) UpdateTravelerStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	var data models.TravelerStoryDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.UpdateTravelerStory(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) DeleteTravelerStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	err := h.travelexpersvcs.DeleteTravlerStory(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

func (h *TravelExperHandler) RestoreTravelerStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	err := h.travelexpersvcs.RestoreTravlerStory(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

//client

func (h *TravelExperHandler) GetClientStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // story id

	result, err := h.travelexpersvcs.GetClientStory(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) GetClientStories(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	query := r.URL.Query().Get("query")

	result, err := h.travelexpersvcs.GetClientStories(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) AddClientStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.ClientStoryDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.AddClientStory(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) UpdateClientStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	var data models.ClientStoryDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.UpdateClientStory(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) DeleteClientStory(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id") // story id

	err := h.travelexpersvcs.DeleteClientStory(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

// v2
func (h *TravelExperHandler) GetTravelerStoriesV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.GetTravelerStoriesV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelExperHandler) GetClientStoriesV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.travelexpersvcs.GetClientStoriesV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
