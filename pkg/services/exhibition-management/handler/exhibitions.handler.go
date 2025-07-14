package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	exhibition_management "larsa-tourism-microservices/pkg/services/exhibition-management"
	"larsa-tourism-microservices/pkg/services/exhibition-management/models"
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	marketingModels "larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"time"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExhibitionHandler struct {
	exhibitionSvcs       exhibition_management.ExhibitionSvcs
	validationInstance   *validator.Validate
	visitorSvcs          marketing.VisitorSvcs
	exhibitorProfileSvcs marketing.ExhibitorProfileSvcs
	inquirySvcs          marketing.InquirySvcs
}

func NewExhibitionHandler(i *do.Injector, r *chi.Mux) {
	h := &ExhibitionHandler{
		exhibitionSvcs:       do.MustInvoke[exhibition_management.ExhibitionSvcs](i),
		validationInstance:   do.MustInvoke[*validator.Validate](i),
		visitorSvcs:          do.MustInvoke[marketing.VisitorSvcs](i),
		exhibitorProfileSvcs: do.MustInvoke[marketing.ExhibitorProfileSvcs](i),
		inquirySvcs:          do.MustInvoke[marketing.InquirySvcs](i),
	}

	r.Route("/exhibitions", func(r chi.Router) {
		// Public routes
		r.Get("/", helpers.Make(h.Get))
		r.Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.OptionalAuth()).Get("/{id}", helpers.Make(h.GetById))
		r.Get("/{id}/related", helpers.Make(h.GetRelatedExhibitions))

		// Authenticated routes
		r.With(middleware.Auth("authenticate")).Get("/auth", helpers.Make(h.GetAuth))
		r.With(middleware.Auth("authenticate")).Get("/all/auth", helpers.Make(h.GetAllAuth))

		// Protected routes
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Save))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
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
		r.Route("/{exhibitionId}/visitors", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetExhibitionVisitors))
			r.Get("/stats", helpers.Make(h.GetExhibitionVisitorStats))
			r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddExhibitionVisitor))
		})

		r.Route("/{exhibitionId}/inquiries", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetExhibitionInquiries))
			r.Get("/stats", helpers.Make(h.GetExhibitionInquiryStats))
			r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddExhibitionInquiry))
		})

		r.Route("/{exhibitionId}/exhibitor-profiles", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetExhibitionProfiles))
			r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddExhibitionProfile))
		})

		// Overall marketing stats for exhibition
		r.Get("/{exhibitionId}/marketing/stats", helpers.Make(h.GetExhibitionMarketingStats))

		// Debug routes - remove in production
		r.Get("/debug/count", helpers.Make(h.DebugCount))
		r.Get("/debug/raw", helpers.Make(h.DebugRaw))
	})

	// v2 routes with filter support
	r.Route("/exhibitions/v2", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.GetV2))
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

func (h *ExhibitionHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.exhibitionSvcs.GetAll(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetAuth(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := h.exhibitionSvcs.GetAuth(ctx, skip, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetAllAuth(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.exhibitionSvcs.GetAllAuth(ctx)
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

func (h *ExhibitionHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	result, err := h.exhibitionSvcs.Patch(ctx, id, updates)
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
func (h *ExhibitionHandler) GetExhibitionMarketing(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	response := map[string]interface{}{
		"exhibitionId": exhibitionId,
		"availableEndpoints": map[string]string{
			"visitors":  "/exhibitions/" + exhibitionId + "/marketing/visitors",
			"requests":  "/exhibitions/" + exhibitionId + "/marketing/requests",
			"inquiries": "/exhibitions/" + exhibitionId + "/marketing/inquiries",
			"stats":     "/exhibitions/" + exhibitionId + "/marketing/stats",
		},
		"description": "Marketing data and analytics for exhibition",
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExhibitionHandler) GetExhibitionVisitors(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := h.visitorSvcs.GetV2(ctx, skip, limit, &query.Conditions{
		Columns: []query.Column{
			{Name: "exhibitionId", Value: exhibitionId, Exp: "="},
		},
	})
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetExhibitionVisitorStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	// Convert exhibitionId to ObjectID
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	// Get stats for this specific exhibition
	result, err := h.visitorSvcs.GetStats(ctx, &exhibitionObjectId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) AddExhibitionVisitor(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	var data marketingModels.VisitorDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Set exhibitionId from URL parameter
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}
	data.ExhibitionId = exhibitionObjectId

	// If no hotelId provided, use zero ObjectID (exhibition visitor without hotel)
	if data.HotelId.IsZero() {
		data.HotelId = primitive.NilObjectID
	}

	// Validate the data (hotelId is still required in struct but we allow NilObjectID)
	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ExhibitionHandler) GetExhibitionRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Create empty filter - note: current filter doesn't support exhibition filtering
	// This would need to be implemented by adding ExhibitionId to ExhibitorRequestFilter
	filterQuery := filter.ExhibitorRequestFilter{}

	result, err := h.exhibitorProfileSvcs.GetExhibitorRequests(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	// Add note about exhibition filtering
	response := map[string]interface{}{
		"note":         "Currently showing all exhibitor requests - exhibition-specific filtering needs to be implemented",
		"exhibitionId": exhibitionId,
		"data":         result,
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ExhibitionHandler) GetExhibitionInquiries(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Get inquiries for this specific exhibition
	result, err := h.inquirySvcs.GetV2(ctx, skip, limit, &query.Conditions{
		Columns: []query.Column{
			{Name: "exhibitionId", Value: exhibitionId, Exp: "="},
		},
	})
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) GetExhibitionInquiryStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	// Convert exhibitionId to ObjectID
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	// Get inquiry stats for this specific exhibition
	result, err := h.inquirySvcs.GetStats(ctx, &exhibitionObjectId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) AddExhibitionInquiry(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	var data marketingModels.InquiryDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Set exhibitionId from URL parameter
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}
	data.ExhibitionId = exhibitionObjectId

	// If no hotelId provided, use zero ObjectID (exhibition inquiry without hotel)
	if data.HotelId.IsZero() {
		data.HotelId = primitive.NilObjectID
	}

	// Validate the data (hotelId is still required in struct but we allow NilObjectID)
	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.inquirySvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ExhibitionHandler) GetExhibitionProfiles(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Get exhibitor profiles for this specific exhibition
	result, err := h.exhibitorProfileSvcs.GetV2(ctx, skip, limit, &query.Conditions{
		Columns: []query.Column{
			{Name: "exhibitionId", Value: exhibitionId, Exp: "="},
		},
	})
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ExhibitionHandler) AddExhibitionProfile(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	var data marketingModels.ExhibitorProfileDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Set exhibitionId from URL parameter
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}
	data.ExhibitionId = exhibitionObjectId

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.exhibitorProfileSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ExhibitionHandler) GetExhibitionMarketingStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	exhibitionId := chi.URLParam(r, "exhibitionId")

	// Validate exhibition exists
	if err := h.exhibitionSvcs.ValidateExhibitionExists(ctx, exhibitionId); err != nil {
		return err
	}

	// Convert exhibitionId to ObjectID
	exhibitionObjectId, err := primitive.ObjectIDFromHex(exhibitionId)
	if err != nil {
		return helpers.BadRequest("Invalid exhibition ID")
	}

	// Get comprehensive marketing stats for this exhibition
	visitorStats, err := h.visitorSvcs.GetStats(ctx, &exhibitionObjectId)
	if err != nil {
		visitorStats = nil
	}

	inquiryStats, err := h.inquirySvcs.GetStats(ctx, &exhibitionObjectId)
	if err != nil {
		inquiryStats = nil
	}

	response := map[string]interface{}{
		"exhibitionId": exhibitionId,
		"visitorStats": visitorStats,
		"inquiryStats": inquiryStats,
		"lastUpdated":  time.Now(),
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}
