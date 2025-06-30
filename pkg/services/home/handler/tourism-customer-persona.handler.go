package handler

import (
	"github.com/goccy/go-json"

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

type CustomerPersonaHandler struct {
	personaSvcs        home.CustomerPersonaSvcs
	validationInstance *validator.Validate
}

func NewCustomerPersonaHandler(i *do.Injector, r *chi.Mux) {
	h := &CustomerPersonaHandler{
		personaSvcs:        do.MustInvoke[home.CustomerPersonaSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/customer-personas", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/all", helpers.Make(h.GetAll))
		r.Get("/", helpers.Make(h.Get))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (h *CustomerPersonaHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := h.personaSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *CustomerPersonaHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	cursor, count, err := h.personaSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var personas []models.CustomerPersona
	if err := cursor.All(ctx, &personas); err != nil {
		return err
	}

	// Create response with pagination info
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (count + limit - 1) / limit
	}

	response := struct {
		Data       []models.CustomerPersona `json:"data"`
		Pagination struct {
			TotalCount int64 `json:"totalCount"`
			Page       int64 `json:"page"`
			PerPage    int64 `json:"perPage"`
			TotalPages int64 `json:"totalPages"`
		} `json:"pagination"`
	}{
		Data: personas,
	}
	response.Pagination.TotalCount = count
	response.Pagination.Page = (skip / limit) + 1
	response.Pagination.PerPage = limit
	response.Pagination.TotalPages = totalPages

	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *CustomerPersonaHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.personaSvcs.GetAll(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *CustomerPersonaHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CustomerPersonaDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := helpers.GenericValidation(h.validationInstance, data); err != nil {
		return err
	}

	result, err := h.personaSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *CustomerPersonaHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CustomerPersonaDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := helpers.GenericValidation(h.validationInstance, data); err != nil {
		return err
	}

	id := chi.URLParam(r, "id")
	result, err := h.personaSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *CustomerPersonaHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.personaSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *CustomerPersonaHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.personaSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Customer persona deleted successfully",
	})
}
