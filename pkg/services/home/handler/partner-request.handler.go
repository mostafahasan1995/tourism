package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type PartnerRequestHandler struct {
	partnerRequestSvcs home.PartnerRequestSvcs
	validationInstance *validator.Validate
}

func NewPartnerRequestHandler(i *do.Injector, r *chi.Mux) {
	h := &PartnerRequestHandler{
		partnerRequestSvcs: do.MustInvoke[home.PartnerRequestSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/partner-requests", func(r chi.Router) {
		// Public endpoints
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))

		// Protected endpoints
		r.With(middleware.Auth("authenticate")).Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.Get))
		r.With(middleware.Auth("authenticate")).Get("/{id}", helpers.Make(h.GetOne))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (h *PartnerRequestHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := h.partnerRequestSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *PartnerRequestHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	cursor, count, err := h.partnerRequestSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var partnerRequests []models.PartnerRequest
	if err := cursor.All(ctx, &partnerRequests); err != nil {
		return err
	}

	// Create response with pagination info
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (count + limit - 1) / limit
	}

	response := struct {
		Data       []models.PartnerRequest `json:"data"`
		Pagination struct {
			TotalCount int64 `json:"totalCount"`
			Page       int64 `json:"page"`
			PerPage    int64 `json:"perPage"`
			TotalPages int64 `json:"totalPages"`
		} `json:"pagination"`
	}{
		Data: partnerRequests,
	}
	response.Pagination.TotalCount = count
	response.Pagination.Page = (skip / limit) + 1
	response.Pagination.PerPage = limit
	response.Pagination.TotalPages = totalPages

	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *PartnerRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.partnerRequestSvcs.GetAll(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *PartnerRequestHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.PartnerRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := helpers.GenericValidation(h.validationInstance, data); err != nil {
		return err
	}

	result, err := h.partnerRequestSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *PartnerRequestHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.PartnerRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := helpers.GenericValidation(h.validationInstance, data); err != nil {
		return err
	}

	id := chi.URLParam(r, "id")
	result, err := h.partnerRequestSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *PartnerRequestHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.partnerRequestSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *PartnerRequestHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		return helpers.InvalidJSON()
	}

	if statusUpdate.Status == "" {
		return helpers.BadRequest("Status is required")
	}

	result, err := h.partnerRequestSvcs.UpdateStatus(ctx, id, statusUpdate.Status)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *PartnerRequestHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.partnerRequestSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Partner request deleted successfully",
	})
}
