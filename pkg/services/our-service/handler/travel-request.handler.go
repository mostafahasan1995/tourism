package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type TravelRequestHandler struct {
	travelreqsvcs      ourservice.TravelRequestSvcs
	validationInstance *validator.Validate
}

func NewTravelRequestHandler(i *do.Injector, r *chi.Mux) {
	h := &TravelRequestHandler{
		travelreqsvcs:      do.MustInvoke[ourservice.TravelRequestSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/travel-requests", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/", helpers.Make(h.Get))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Get("/my-requests/{status}", helpers.Make(h.MyRequests))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
	})
}

func (h *TravelRequestHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := h.travelreqsvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TravelRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	var data models.TravelRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelRequestHandler) MyRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	status := chi.URLParam(r, "status")

	result, err := h.travelreqsvcs.MyRequests(ctx, status)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelRequestHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // travel request id

	var data models.ChangeStatusDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.UpdateStatus(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
