package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type ExhibitorProfileHandler struct {
	exhibitorProfileSvcs marketing.ExhibitorProfileSvcs
	validationInstance   *validator.Validate
}

func NewExhibitorProfileHandler(i *do.Injector, r *chi.Mux) {
	h := &ExhibitorProfileHandler{
		exhibitorProfileSvcs: do.MustInvoke[marketing.ExhibitorProfileSvcs](i),
		validationInstance:   do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/marketing/exhibitor-profiles", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/hotel/{hotelId}", helpers.Make(h.GetByHotelId))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Post("/request", helpers.Make(h.AddExhibitorRequest))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		r.With(middleware.Auth("authenticate")).Get("/requests", helpers.Make(h.GetExhibitorRequests))
		r.With(middleware.Auth("authenticate")).Patch("/requests/{id}/status", helpers.Make(h.UpdateExhibitorRequestStatus))

		r.With(middleware.Auth("authenticate")).Post("/{id}/dynamic-sections", helpers.Make(h.AddDynamicSection))
		r.With(middleware.Auth("authenticate")).Put("/{id}/dynamic-sections/{sectionId}", helpers.Make(h.UpdateDynamicSection))
		r.With(middleware.Auth("authenticate")).Delete("/{id}/dynamic-sections/{sectionId}", helpers.Make(h.DeleteDynamicSection))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/dynamic-sections/reorder", helpers.Make(h.ReorderDynamicSections))

		r.With(middleware.Auth("authenticate")).Patch("/{id}/facility", helpers.Make(h.UpdateFacility))
		r.With(middleware.Auth("authenticate")).Put("/{id}/hero-section", helpers.Make(h.UpdateHeroSection))

		r.With(middleware.Auth("authenticate")).Patch("/{id}/publish", helpers.Make(h.Publish))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/unpublish", helpers.Make(h.Unpublish))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/toggle-active", helpers.Make(h.ToggleActive))
	})
}

func (h *ExhibitorProfileHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Parse filter parameters
	var filterQuery filter.ExhibitorProfileFilter
	if err := json.NewDecoder(r.Body).Decode(&filterQuery); err != nil {
		// If no body or invalid JSON, use query parameters
		filterQuery = filter.ExhibitorProfileFilter{
			Page: 1,
			Size: int(limit),
		}

	}

	result, err := h.exhibitorProfileSvcs.Get(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitorProfileSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) GetByHotelId(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	hotelId := chi.URLParam(r, "hotelId")

	result, err := h.exhibitorProfileSvcs.GetByHotelId(ctx, hotelId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.ExhibitorProfileDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *ExhibitorProfileHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.ExhibitorProfileDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.exhibitorProfileSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.exhibitorProfileSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Exhibitor profile deleted successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ExhibitorProfileHandler) AddDynamicSection(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.AddDynamicSectionDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.AddDynamicSection(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) UpdateDynamicSection(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	sectionId := chi.URLParam(r, "sectionId")

	var data models.UpdateDynamicSectionDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.exhibitorProfileSvcs.UpdateDynamicSection(ctx, id, sectionId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) DeleteDynamicSection(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	sectionId := chi.URLParam(r, "sectionId")

	result, err := h.exhibitorProfileSvcs.DeleteDynamicSection(ctx, id, sectionId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) ReorderDynamicSections(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var sectionOrders map[string]int
	if err := json.NewDecoder(r.Body).Decode(&sectionOrders); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.exhibitorProfileSvcs.ReorderDynamicSections(ctx, id, sectionOrders)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) UpdateFacility(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.FacilityUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.UpdateFacility(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) UpdateHeroSection(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.HeroSection
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.exhibitorProfileSvcs.UpdateHeroSection(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) Publish(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitorProfileSvcs.Publish(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) Unpublish(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitorProfileSvcs.Unpublish(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) ToggleActive(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitorProfileSvcs.ToggleActive(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) AddExhibitorRequest(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.ExhibitorRequestDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.AddExhibitorRequest(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *ExhibitorProfileHandler) GetExhibitorRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Parse filter parameters
	var filterQuery filter.ExhibitorRequestFilter

	// Use the filter's ParseQueryParams method
	filterQuery.ParseQueryParams(r.URL.Query())

	result, err := h.exhibitorProfileSvcs.GetExhibitorRequests(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ExhibitorProfileHandler) UpdateExhibitorRequestStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data struct {
		Status string `json:"status" validate:"required,oneof=Pending Replied Closed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := helpers.GenericValidation(h.validationInstance, &data); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.UpdateExhibitorRequestStatus(ctx, id, data.Status)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
