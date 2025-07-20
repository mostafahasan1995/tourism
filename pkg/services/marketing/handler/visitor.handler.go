package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"

	"larsa-tourism-microservices/pkg/query"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VisitorHandler struct {
	visitorSvcs        marketing.VisitorSvcs
	validationInstance *validator.Validate
}

func NewVisitorHandler(i *do.Injector, r *chi.Mux) {
	h := &VisitorHandler{
		visitorSvcs:        do.MustInvoke[marketing.VisitorSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/marketing/visitors", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/stats", helpers.Make(h.GetStats))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/{id}/activities", helpers.Make(h.GetVisitorActivities))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		r.With(middleware.Auth("authenticate")).Post("/{id}/tags", helpers.Make(h.AddTag))
		r.With(middleware.Auth("authenticate")).Delete("/{id}/tags/{tag}", helpers.Make(h.RemoveTag))
		r.With(middleware.Auth("authenticate")).Put("/{id}/tags", helpers.Make(h.UpdateTags))

		r.With(middleware.Auth("authenticate")).Post("/{id}/interests", helpers.Make(h.AddInterest))
		r.With(middleware.Auth("authenticate")).Delete("/{id}/interests/{interest}", helpers.Make(h.RemoveInterest))
		r.With(middleware.Auth("authenticate")).Put("/{id}/interests", helpers.Make(h.UpdateInterests))

		r.With(middleware.Auth("authenticate")).Patch("/bulk", helpers.Make(h.BulkUpdate))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/vip", helpers.Make(h.MarkAsVIP))
		r.With(middleware.Auth("authenticate")).Delete("/{id}/vip", helpers.Make(h.RemoveVIP))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/toggle-active", helpers.Make(h.ToggleActive))

		r.With(middleware.Auth("authenticate")).Post("/{id}/activities", helpers.Make(h.AddActivity))

		r.Get("/hotel/{hotelId}", helpers.Make(h.GetByHotelId))
		r.Get("/hotel/{hotelId}/stats", helpers.Make(h.GetHotelStats))
	})

	// v2 routes with filter support
	r.Route("/marketing/visitors/v2", func(r chi.Router) {
		r.Post("/", helpers.Make(h.GetV2))
	})
}

func (h *VisitorHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Initialize filter with default values
	filterQuery := filter.VisitorFilter{
		Page: 1,
		Size: int(limit),
	}

	// Parse query parameters (GET requests typically use query params, not body)
	if fullName := r.URL.Query().Get("fullName"); fullName != "" {
		filterQuery.FullName = fullName
	}
	if nationality := r.URL.Query().Get("nationality"); nationality != "" {
		filterQuery.Nationality = nationality
	}
	if email := r.URL.Query().Get("email"); email != "" {
		filterQuery.Email = email
	}
	if phone := r.URL.Query().Get("phone"); phone != "" {
		filterQuery.Phone = phone
	}
	if searchText := r.URL.Query().Get("search"); searchText != "" {
		filterQuery.SearchText = searchText
	}
	if source := r.URL.Query().Get("source"); source != "" {
		filterQuery.Source = source
	}
	if hotelId := r.URL.Query().Get("hotelId"); hotelId != "" {
		if objId, err := primitive.ObjectIDFromHex(hotelId); err == nil {
			filterQuery.HotelId = objId
		}
	}

	// Handle boolean filters
	if isActive := r.URL.Query().Get("isActive"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filterQuery.IsActive = &active
		}
	}
	if isVIP := r.URL.Query().Get("isVIP"); isVIP != "" {
		if vip, err := strconv.ParseBool(isVIP); err == nil {
			filterQuery.IsVIP = &vip
		}
	}

	result, err := h.visitorSvcs.Get(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.visitorSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) GetByHotelId(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	hotelId := chi.URLParam(r, "hotelId")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	hotelObjectId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return helpers.BadRequest("Invalid hotel ID")
	}

	filterQuery := filter.VisitorFilter{
		Page:    1,
		Size:    int(limit),
		HotelId: hotelObjectId,
	}

	result, err := h.visitorSvcs.Get(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.VisitorDto

	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *VisitorHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.VisitorDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.visitorSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.visitorSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Visitor deleted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *VisitorHandler) AddTag(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.AddVisitorTagDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.AddTag(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) RemoveTag(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	tag := chi.URLParam(r, "tag")

	result, err := h.visitorSvcs.RemoveTag(ctx, id, tag)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) UpdateTags(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.UpdateVisitorTagsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.UpdateTags(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) AddInterest(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.AddVisitorInterestDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.AddInterest(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) RemoveInterest(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	interest := chi.URLParam(r, "interest")

	result, err := h.visitorSvcs.RemoveInterest(ctx, id, interest)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) UpdateInterests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.UpdateVisitorInterestsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.visitorSvcs.UpdateInterests(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.visitorSvcs.GetStats(ctx, nil)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) GetHotelStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	hotelId := chi.URLParam(r, "hotelId")

	hotelObjectId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return helpers.BadRequest("Invalid hotel ID")
	}

	result, err := h.visitorSvcs.GetStats(ctx, &hotelObjectId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) BulkUpdate(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.BulkUpdateVisitorsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	err := h.visitorSvcs.BulkUpdate(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Visitors updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *VisitorHandler) MarkAsVIP(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.visitorSvcs.MarkAsVIP(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) RemoveVIP(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.visitorSvcs.RemoveVIP(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) ToggleActive(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.visitorSvcs.ToggleActive(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *VisitorHandler) AddActivity(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.VisitorActivity
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	visitorId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.BadRequest("Invalid visitor ID")
	}

	data.VisitorId = visitorId

	err = h.visitorSvcs.AddActivity(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Activity added successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *VisitorHandler) GetVisitorActivities(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := h.visitorSvcs.GetVisitorActivities(ctx, id, skip, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// GetV2 - Get visitors with filters in request body
func (h *VisitorHandler) GetV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.visitorSvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
