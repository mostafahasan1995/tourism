package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsHandler struct {
	reviewsSvcs home.ReviewsSvcs
}

func NewReviewsHandler(i *do.Injector, r *chi.Mux) {
	h := &ReviewsHandler{
		reviewsSvcs: do.MustInvoke[home.ReviewsSvcs](i),
	}

	r.Route("/reviews", func(r chi.Router) {
		// Public endpoints
		r.Get("/", helpers.Make(h.GetAll))
		r.Get("/stats", helpers.Make(h.GetStats))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Post("/", helpers.Make(h.Add)) // Public endpoint for submitting reviews

		// Program-specific endpoints
		r.Get("/program/{programId}", helpers.Make(h.GetProgramReviews))
		r.Get("/program/{programId}/stats", helpers.Make(h.GetProgramStats))

		// Protected endpoints (require authentication)
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddMany))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		// Reply management endpoints
		r.With(middleware.Auth("authenticate")).Post("/{id}/replies", helpers.Make(h.AddReply))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
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

	// Handle direct parameters
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			reviewsFilter.Page = page
		}
	}
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			reviewsFilter.Size = size
		}
	}
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		reviewsFilter.Status = statusParam
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

	// Program filters
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

	result, err := h.reviewsSvcs.GetAll(ctx, reviewsFilter)
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

func (h *ReviewsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var data models.ReviewDto

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate required fields
	if data.Customer == "" {
		return helpers.BadRequest("Customer type is required (customer or agent)")
	}
	if data.Customer != "customer" && data.Customer != "agent" {
		return helpers.BadRequest("Customer type must be either 'customer' or 'agent'")
	}
	if data.Username == "" {
		return helpers.BadRequest("Username is required")
	}
	if data.ProgramId.IsZero() {
		return helpers.BadRequest("Program ID is required")
	}
	if data.Description == "" {
		return helpers.BadRequest("Description is required")
	}
	if data.Value < 1 || data.Value > 5 {
		return helpers.BadRequest("Rating value must be between 1 and 5")
	}

	// Set default status for new reviews
	if data.Status == "" {
		data.Status = "pending"
	}

	// Set date if not provided
	if data.Date.IsZero() {
		data.Date = time.Now()
	}

	// Initialize replies if nil
	if data.Replies == nil {
		data.Replies = []models.ReviewReply{}
	}

	// Initialize images if nil
	if data.Images == nil {
		data.Images = []models.ReviewImage{}
	}

	err := h.reviewsSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"message": "Review submitted successfully",
		"status":  "pending",
	}
	return helpers.WriteJson(w, http.StatusCreated, response)
}

func (h *ReviewsHandler) AddMany(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.ReviewDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.reviewsSvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Reviews created successfully",
	}
	return helpers.WriteJson(w, http.StatusCreated, response)
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

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.reviewsSvcs.UpdateReplyStatus(ctx, id, statusUpdate.Status)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Review status updated successfully",
	}
	return helpers.WriteJson(w, http.StatusOK, response)
}

func (h *ReviewsHandler) GetProgramReviews(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	programId := chi.URLParam(r, "programId")
	reviewsFilter := h.parseReviewsFilter(r)

	result, err := h.reviewsSvcs.GetProgramReviews(ctx, programId, reviewsFilter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *ReviewsHandler) GetProgramStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	programId := chi.URLParam(r, "programId")

	result, err := h.reviewsSvcs.GetProgramStats(ctx, programId)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}
