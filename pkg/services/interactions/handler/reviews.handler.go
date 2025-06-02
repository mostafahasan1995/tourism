package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/filter"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsHandler struct {
	reviewsSvcs interactions.ReviewsSvcs
}

func NewReviewsHandler(i *do.Injector, r *chi.Mux) {
	h := &ReviewsHandler{
		reviewsSvcs: do.MustInvoke[interactions.ReviewsSvcs](i),
	}

	r.Route("/reviews", func(r chi.Router) {
		// Public endpoints
		r.Get("/", helpers.Make(h.GetAll))
		r.Get("/stats", helpers.Make(h.GetStats))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.Add)) // Public endpoint with optional auth

		// Filter endpoints
		r.Get("/status/{status}", helpers.Make(h.GetByStatus))
		r.Get("/user/{userId}", helpers.Make(h.GetByUserId))

		// Entity-specific endpoints (e.g., /reviews/hotels/{hotelId})
		r.Route("/{entityType}/{refId}", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetEntityReviews))
			r.Get("/stats", helpers.Make(h.GetEntityStats))
			r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.AddEntityReview))
		})

		// Protected endpoints (require authentication)

		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		// Reply management endpoints
		r.With(middleware.Auth("authenticate")).Post("/{id}/replies", helpers.Make(h.AddReply))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/approve", helpers.Make(h.ApproveReview))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/reject", helpers.Make(h.RejectReview))
	})
}

// Helper function to parse filter parameters
func (h *ReviewsHandler) parseReviewsFilter(r *http.Request) filter.ReviewsFilter {
	var reviewsFilter filter.ReviewsFilter

	// Handle JSON query parameter
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		json.Unmarshal([]byte(filterParam), &reviewsFilter)
	}

	// Handle direct parameters (removed page and size)
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		reviewsFilter.Status = statusParam
	}

	// Entity filters
	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		reviewsFilter.Type = typeParam
	}
	if refParam := r.URL.Query().Get("ref"); refParam != "" {
		if refId, err := primitive.ObjectIDFromHex(refParam); err == nil {
			reviewsFilter.Ref = refId
		}
	}

	// Customer filters
	if customerParam := r.URL.Query().Get("customer"); customerParam != "" {
		reviewsFilter.Customer = customerParam
	}
	if usernameParam := r.URL.Query().Get("username"); usernameParam != "" {
		reviewsFilter.Username = usernameParam
	}
	if userIdParam := r.URL.Query().Get("userId"); userIdParam != "" {
		reviewsFilter.UserId = userIdParam
	}

	// Program filters (legacy support)
	if programIdParam := r.URL.Query().Get("programId"); programIdParam != "" {
		if programId, err := primitive.ObjectIDFromHex(programIdParam); err == nil {
			reviewsFilter.ProgramId = programId
		}
	}

	// Independent destination and countries
	if destinationParam := r.URL.Query().Get("destination"); destinationParam != "" {
		reviewsFilter.Destination = destinationParam
	}
	if countryParam := r.URL.Query().Get("country"); countryParam != "" {
		reviewsFilter.Country = countryParam
	}

	// Rating filters
	if minRatingParam := r.URL.Query().Get("minRating"); minRatingParam != "" {
		if minRating, err := strconv.ParseFloat(minRatingParam, 64); err == nil {
			reviewsFilter.MinRating = minRating
		}
	}
	if maxRatingParam := r.URL.Query().Get("maxRating"); maxRatingParam != "" {
		if maxRating, err := strconv.ParseFloat(maxRatingParam, 64); err == nil {
			reviewsFilter.MaxRating = maxRating
		}
	}

	// Content filters
	if searchParam := r.URL.Query().Get("search"); searchParam != "" {
		reviewsFilter.Search = searchParam
	}
	if hasImagesParam := r.URL.Query().Get("hasImages"); hasImagesParam != "" {
		if hasImages, err := strconv.ParseBool(hasImagesParam); err == nil {
			reviewsFilter.HasImages = &hasImages
		}
	}
	if hasAdviceParam := r.URL.Query().Get("hasAdvice"); hasAdviceParam != "" {
		if hasAdvice, err := strconv.ParseBool(hasAdviceParam); err == nil {
			reviewsFilter.HasAdvice = &hasAdvice
		}
	}

	return reviewsFilter
}

func (h *ReviewsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	reviewsFilter := h.parseReviewsFilter(r)

	// Default to approved reviews if no status filter is provided
	if reviewsFilter.Status == "" {
		reviewsFilter.Status = "approved"
	}

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := h.reviewsSvcs.GetAll(ctx, reviewsFilter, page, perPage)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.reviewsSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.reviewsSvcs.GetStats(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetEntityReviews(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	entityType := chi.URLParam(r, "entityType")
	refId := chi.URLParam(r, "refId")
	reviewsFilter := h.parseReviewsFilter(r)

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := h.reviewsSvcs.GetEntityReviews(ctx, entityType, refId, reviewsFilter, page, perPage)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetEntityStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	entityType := chi.URLParam(r, "entityType")
	refId := chi.URLParam(r, "refId")

	result, err := h.reviewsSvcs.GetEntityStats(ctx, entityType, refId)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.ReviewDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate required fields
	if data.Type == "" {
		return helpers.BadRequest("Review type is required (e.g., 'hotel', 'program', etc.)")
	}
	if data.Ref.IsZero() {
		return helpers.BadRequest("Reference ID is required")
	}
	if data.Description == "" {
		return helpers.BadRequest("Description is required")
	}
	if data.Value < 1 || data.Value > 5 {
		return helpers.BadRequest("Rating value must be between 1 and 5")
	}

	// Check if user is authenticated for additional validation
	user, _ := r.Context().Value(util.ReqUser).(*types.User)
	if user == nil {
		// For unauthenticated users, check type-specific requirements
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
	}

	result, err := h.reviewsSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *ReviewsHandler) AddEntityReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	entityType := chi.URLParam(r, "entityType")
	refId := chi.URLParam(r, "refId")

	var data models.ReviewDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Set entity type and ref from URL parameters
	data.Type = entityType
	refObjectId, err := primitive.ObjectIDFromHex(refId)
	if err != nil {
		return helpers.BadRequest("Invalid reference ID")
	}
	data.Ref = refObjectId

	// Validate required fields
	if data.Description == "" {
		return helpers.BadRequest("Description is required")
	}
	if data.Value < 1 || data.Value > 5 {
		return helpers.BadRequest("Rating value must be between 1 and 5")
	}

	// Check if user is authenticated for additional validation
	user, _ := r.Context().Value(util.ReqUser).(*types.User)
	if user == nil {
		// For unauthenticated users, check type-specific requirements
		if data.Type == "hotel" && (data.FirstName == "" || data.LastName == "" || data.Email == "") {
			return helpers.BadRequest("firstName, lastName, and email are required for unauthenticated hotel reviews")
		}
	}

	result, err := h.reviewsSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *ReviewsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.ReviewDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.reviewsSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review updated successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.reviewsSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review updated successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
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
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) AddReply(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var reply models.ReviewReply
	if err := json.NewDecoder(r.Body).Decode(&reply); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate reply text
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
	return helpers.WriteJson(w, http.StatusCreated, response)
}

func (h *ReviewsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var statusUpdate models.ReviewStatusDto
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.reviewsSvcs.UpdateReviewStatus(ctx, id, statusUpdate.Status)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review status updated successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) ApproveReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.reviewsSvcs.ApproveReview(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review approved successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) RejectReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.reviewsSvcs.RejectReview(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review rejected successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) GetByStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	status := chi.URLParam(r, "status")

	reviewsFilter := h.parseReviewsFilter(r)
	reviewsFilter.Status = status // Override status from URL

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := h.reviewsSvcs.GetAll(ctx, reviewsFilter, page, perPage)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetByUserId(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	userId := chi.URLParam(r, "userId")

	reviewsFilter := h.parseReviewsFilter(r)
	reviewsFilter.UserId = userId // Set userId filter

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := h.reviewsSvcs.GetAll(ctx, reviewsFilter, page, perPage)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}
