package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	exhibition_management "larsa-tourism-microservices/pkg/services/exhibition-management"
	"larsa-tourism-microservices/pkg/services/exhibition-management/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type ExhibitionHandler struct {
	exhibitionSvcs     exhibition_management.ExhibitionSvcs
	validationInstance *validator.Validate
}

func NewExhibitionHandler(i *do.Injector, r *chi.Mux) {
	h := &ExhibitionHandler{
		exhibitionSvcs:     do.MustInvoke[exhibition_management.ExhibitionSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/exhibitions", func(r chi.Router) {
		// Public routes
		r.Get("/", helpers.Make(h.Get))
		r.Get("/{id}", helpers.Make(h.GetById))
		r.Get("/{id}/related", helpers.Make(h.GetRelatedExhibitions))

		// Protected routes
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Save))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/toggle", helpers.Make(h.Toggle))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		// Ad management routes
		r.Route("/{exhibitionId}/ads", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetExhibitionAds))
			r.With(middleware.Auth("authenticate")).Put("/{adId}", helpers.Make(h.UpdateAd))
			r.With(middleware.Auth("authenticate")).Delete("/{adId}", helpers.Make(h.DeleteAd))
			r.With(middleware.Auth("authenticate")).Patch("/{adId}/toggle", helpers.Make(h.ToggleAd))
		})

		// Marketing-related routes
		r.Route("/{exhibitionId}/marketing", func(r chi.Router) {
			r.Get("/visitors", helpers.Make(h.GetExhibitionVisitors))
			r.Get("/requests", helpers.Make(h.GetExhibitionRequests))
			r.Get("/inquiries", helpers.Make(h.GetExhibitionInquiries))
			r.Get("/stats", helpers.Make(h.GetExhibitionStats))
		})
	})

	// v2 routes with filter support
	r.Route("/exhibitions/v2", func(r chi.Router) {
		r.Post("/", helpers.Make(h.GetV2))
	})

	// Debug routes (remove in production)
	r.Route("/exhibitions/debug", func(r chi.Router) {
		r.Get("/count", helpers.Make(h.DebugCount))
		r.Get("/raw", helpers.Make(h.DebugRaw))
	})
}

func (h *ExhibitionHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.Get(ctx, skip, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitionSvcs.GetById(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetRelatedExhibitions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.exhibitionSvcs.GetRelatedExhibitions(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) Save(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.ExhibitionDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate the data
	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.Save(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ExhibitionHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.ExhibitionDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate the data
	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) Toggle(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data struct {
		IsActive bool `json:"isActive"`
	}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	result, err := h.exhibitionSvcs.Toggle(ctx, id, data.IsActive)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	if err := h.exhibitionSvcs.Delete(ctx, id); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{
		"message": "Exhibition deleted successfully",
	})
}

// Ad Management Handlers

func (h *ExhibitionHandler) GetExhibitionAds(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	result, err := h.exhibitionSvcs.GetExhibitionAds(ctx, exhibitionId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) UpdateAd(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")
	adId := chi.URLParam(r, "adId")

	var data models.AdDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.UpdateAd(ctx, exhibitionId, adId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) DeleteAd(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")
	adId := chi.URLParam(r, "adId")

	if err := h.exhibitionSvcs.DeleteAd(ctx, exhibitionId, adId); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{
		"message": "Ad deleted successfully",
	})
}

func (h *ExhibitionHandler) ToggleAd(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")
	adId := chi.URLParam(r, "adId")

	var data struct {
		IsActive bool `json:"isActive"`
	}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	result, err := h.exhibitionSvcs.ToggleAd(ctx, exhibitionId, adId, data.IsActive)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// v2 - Get exhibitions with filters in request body
func (h *ExhibitionHandler) GetV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// Debug methods - remove in production
func (h *ExhibitionHandler) DebugCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.exhibitionSvcs.DebugCount(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) DebugRaw(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.exhibitionSvcs.DebugRaw(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// Marketing-related methods (These should be moved to marketing module or inject marketing services)
func (h *ExhibitionHandler) GetExhibitionVisitors(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// TODO: Query marketing.visitors collection where exhibitionId = exhibitionId
	// For now, return placeholder response
	response := map[string]interface{}{
		"message":      "This endpoint should query visitors collection by exhibitionId: " + exhibitionId,
		"exhibitionId": exhibitionId,
		"note":         "Implement by querying marketing.visitors where exhibitionId = " + exhibitionId,
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExhibitionHandler) GetExhibitionRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// TODO: Query marketing.exhibitor_requests collection where exhibitionId = exhibitionId
	response := map[string]interface{}{
		"message":      "This endpoint should query exhibitor requests collection by exhibitionId: " + exhibitionId,
		"exhibitionId": exhibitionId,
		"note":         "Implement by querying marketing.exhibitor_requests where exhibitionId = " + exhibitionId,
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExhibitionHandler) GetExhibitionInquiries(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// TODO: Query marketing.inquiries collection where exhibitionId = exhibitionId
	response := map[string]interface{}{
		"message":      "This endpoint should query inquiries collection by exhibitionId: " + exhibitionId,
		"exhibitionId": exhibitionId,
		"note":         "Implement by querying marketing.inquiries where exhibitionId = " + exhibitionId,
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExhibitionHandler) GetExhibitionStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// TODO: Aggregate stats from marketing collections
	response := map[string]interface{}{
		"message":      "This endpoint should return aggregated stats for exhibitionId: " + exhibitionId,
		"exhibitionId": exhibitionId,
		"stats": map[string]interface{}{
			"totalVisitors":  0, // Count from visitors collection
			"totalRequests":  0, // Count from exhibitor_requests collection
			"totalInquiries": 0, // Count from inquiries collection
			"note":           "Implement by aggregating counts from marketing collections",
		},
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}
