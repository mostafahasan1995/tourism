package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/filter"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InquiryHandler struct {
	inquirySvcs        marketing.InquirySvcs
	validationInstance *validator.Validate
}

func NewInquiryHandler(i *do.Injector, r *chi.Mux) {
	h := &InquiryHandler{
		inquirySvcs:        do.MustInvoke[marketing.InquirySvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/marketing/inquiries", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/stats", helpers.Make(h.GetStats))
		r.Get("/unread-count", helpers.Make(h.GetUnreadCount))
		r.Get("/{id}", helpers.Make(h.GetOne))

		r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.Add))

		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		r.With(middleware.Auth("authenticate")).Post("/{id}/replies", helpers.Make(h.AddReply))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/read", helpers.Make(h.MarkAsRead))

		// r.With(middleware.Auth("authenticate")).Patch("/bulk/status", helpers.Make(h.BulkUpdateStatus))
		// r.With(middleware.Auth("authenticate")).Patch("/bulk/read", helpers.Make(h.BulkMarkAsRead))

		r.Get("/hotel/{hotelId}", helpers.Make(h.GetByHotelId))
		r.Get("/hotel/{hotelId}/stats", helpers.Make(h.GetHotelStats))
	})
}

func (h *InquiryHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Parse filter parameters
	var filterQuery filter.InquiryFilter
	if err := json.NewDecoder(r.Body).Decode(&filterQuery); err != nil {
		// If no body or invalid JSON, use query parameters
		filterQuery = filter.InquiryFilter{
			Page: 1,
			Size: int(limit),
		}

		if status := r.URL.Query().Get("status"); status != "" {
			inquiryStatus := models.InquiryStatus(status)
			filterQuery.Status = &inquiryStatus
		}
		if priority := r.URL.Query().Get("priority"); priority != "" {
			filterQuery.Priority = priority
		}
		if source := r.URL.Query().Get("source"); source != "" {
			filterQuery.Source = source
		}
		if visitorName := r.URL.Query().Get("visitorName"); visitorName != "" {
			filterQuery.VisitorName = visitorName
		}
		if searchText := r.URL.Query().Get("search"); searchText != "" {
			filterQuery.SearchText = searchText
		}
	}

	result, err := h.inquirySvcs.Get(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.inquirySvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) GetByHotelId(w http.ResponseWriter, r *http.Request) error {
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

	filterQuery := filter.InquiryFilter{
		Page:    1,
		Size:    int(limit),
		HotelId: hotelObjectId,
	}

	result, err := h.inquirySvcs.Get(ctx, skip, limit, filterQuery)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.InquiryDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.inquirySvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *InquiryHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.InquiryDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.inquirySvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.inquirySvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.inquirySvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Inquiry deleted successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *InquiryHandler) AddReply(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.AddInquiryReplyDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.inquirySvcs.AddReply(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.UpdateInquiryStatusDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.inquirySvcs.UpdateStatus(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.MarkAsReadDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.inquirySvcs.MarkAsRead(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.inquirySvcs.GetStats(ctx, nil)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) GetHotelStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	hotelId := chi.URLParam(r, "hotelId")

	hotelObjectId, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return helpers.BadRequest("Invalid hotel ID")
	}

	result, err := h.inquirySvcs.GetStats(ctx, &hotelObjectId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InquiryHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var hotelObjectId *primitive.ObjectID
	if hotelId := r.URL.Query().Get("hotelId"); hotelId != "" {
		objId, err := primitive.ObjectIDFromHex(hotelId)
		if err != nil {
			return helpers.BadRequest("Invalid hotel ID")
		}
		hotelObjectId = &objId
	}

	count, err := h.inquirySvcs.GetUnreadCount(ctx, hotelObjectId)
	if err != nil {
		return err
	}

	response := map[string]int64{
		"unreadCount": count,
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

// func (h *InquiryHandler) BulkUpdateStatus(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data struct {
// 		InquiryIds []string             `json:"inquiryIds"`
// 		Status     models.InquiryStatus `json:"status"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	if len(data.InquiryIds) == 0 {
// 		return helpers.BadRequest("No inquiry IDs provided")
// 	}

// 	err := h.inquirySvcs.BulkUpdateStatus(ctx, data.InquiryIds, data.Status)
// 	if err != nil {
// 		return err
// 	}

// 	response := map[string]string{
// 		"message": "Status updated successfully",
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, response)
// }

// func (h *InquiryHandler) BulkMarkAsRead(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data struct {
// 		InquiryIds []string `json:"inquiryIds"`
// 		IsRead     bool     `json:"isRead"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	if len(data.InquiryIds) == 0 {
// 		return helpers.BadRequest("No inquiry IDs provided")
// 	}

// 	err := h.inquirySvcs.BulkMarkAsRead(ctx, data.InquiryIds, data.IsRead)
// 	if err != nil {
// 		return err
// 	}

// 	response := map[string]string{
// 		"message": "Read status updated successfully",
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, response)
// }
