package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsHandler struct {
	reviewsSvcs        interactions.ReviewsSvcs
	validationInstance *validator.Validate
}

func NewReviewsHandler(i *do.Injector, r *chi.Mux) {
	h := &ReviewsHandler{
		reviewsSvcs:        do.MustInvoke[interactions.ReviewsSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/reviews", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/stats", helpers.Make(h.GetStats))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.Add))
		r.Get("/status/{status}", helpers.Make(h.GetByStatus))
		r.Get("/user/{userId}", helpers.Make(h.GetByUserId))
		r.Route("/{entityType}/{refId}", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetEntityReviews))
			//wesite review
			r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.AddEntityReview))
		})
		//dashboard
		r.With(middleware.Auth("authenticate")).Post("/dashboard", helpers.Make(h.AddDashboardReview))
		r.Get("/all", helpers.Make(h.GetAllApproved))
		r.With(middleware.Auth("authenticate")).Get("/admin", helpers.Make(h.GetAllWithPagination))
		r.With(middleware.Auth("authenticate")).Get("/admin/all", helpers.Make(h.GetAllWithoutPagination))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
		r.With(middleware.Auth("authenticate")).Post("/{id}/replies", helpers.Make(h.AddReply))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/approve", helpers.Make(h.ApproveReview))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/reject", helpers.Make(h.RejectReview))
	})
}
func (h *ReviewsHandler) GetAllApproved(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.reviewsSvcs.GetAllApproved(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetAllWithoutPagination(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.reviewsSvcs.GetAllWithoutPagination(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetAllWithPagination(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.reviewsSvcs.GetAllWithPagination(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.reviewsSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.reviewsSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetEntityReviews(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	entityType := chi.URLParam(r, "entityType")
	refId := chi.URLParam(r, "refId")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := `{"type":"` + entityType + `","ref":"` + refId + `"}`

	if existingQuery := r.URL.Query().Get("query"); existingQuery != "" {
		query = `{"type":"` + entityType + `","ref":"` + refId + `"}`
	}

	result, err := h.reviewsSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.ReviewDto

	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	user, _ := r.Context().Value(util.ReqUser).(*types.User)
	if user == nil {
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
	}

	result, err := h.reviewsSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ReviewsHandler) AddEntityReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	entityType := chi.URLParam(r, "entityType")
	refId := chi.URLParam(r, "refId")

	var data models.ReviewDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	data.Type = entityType

	// Handle general type with no refId
	if entityType == "general" {
		data.Ref = primitive.NilObjectID
	} else {
		refObjectId, err := primitive.ObjectIDFromHex(refId)
		if err != nil {
			return helpers.BadRequest("Invalid reference ID")
		}
		data.Ref = refObjectId
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	user, _ := r.Context().Value(util.ReqUser).(*types.User)
	if user == nil {
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
	}

	result, err := h.reviewsSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *ReviewsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.ReviewDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.reviewsSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.reviewsSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.reviewsSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review deleted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *ReviewsHandler) AddReply(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var reply models.ReviewReply
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &reply); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if reply.Text == "" {
		return helpers.BadRequest("Reply text is required")
	}

	err := h.reviewsSvcs.AddReply(ctx, id, reply)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Reply added successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *ReviewsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var statusUpdate models.ReviewStatusDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &statusUpdate); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.reviewsSvcs.UpdateReviewStatus(ctx, id, statusUpdate.Status)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) ApproveReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.reviewsSvcs.ApproveReview(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) RejectReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.reviewsSvcs.RejectReview(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.reviewsSvcs.GetStats(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetByStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	status := chi.URLParam(r, "status")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := `{"status":"` + status + `"}`

	result, err := h.reviewsSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetByUserId(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	userId := chi.URLParam(r, "userId")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := `{"userId":"` + userId + `"}`

	result, err := h.reviewsSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *ReviewsHandler) AddDashboardReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.ReviewDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate review type
	// if data.Type != "agent" && data.Type != "hotel" && data.Type != "destination" && data.Type != "general" {
	// 	return helpers.BadRequest("Invalid review type. Must be one of: agent, hotel, destination, general")
	// }

	// // Handle general type with no refId
	if data.Type == "general" {
		data.Ref = primitive.NilObjectID
	} else if data.Ref.IsZero() {
		// Validate refId for non-general types
		return helpers.BadRequest("refId is required for non-general review types")
	}

	// Validate countries for destination type
	if data.Type == "destination" && (len(data.Countries) == 0) {
		return helpers.BadRequest("countries are required for destination reviews")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.reviewsSvcs.AddFromDashboard(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}
