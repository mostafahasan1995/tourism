package handler

import (
	"encoding/json"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourService "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type HotelsHandler struct {
	hotelssvcs ourService.HotelsSvcs
}

func NewHotelsHandler(i *do.Injector, r *chi.Mux) {
	h := &HotelsHandler{
		hotelssvcs: do.MustInvoke[ourService.HotelsSvcs](i),
	}

	r.Route("/hotels", func(r chi.Router) {

		r.Get("/{id}", helpers.Make(h.GetOne))

		r.Get("/", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

		// Hotel review endpoints
		r.With(middleware.OptionalAuth()).Post("/{id}/reviews", helpers.Make(h.AddReview))
		r.Get("/{id}/reviews", helpers.Make(h.GetHotelReviews))

		// Debug endpoint to see all reviews
		r.Get("/debug/reviews", helpers.Make(h.GetAllReviews))

		// Review management endpoints (require authentication)
		r.With(middleware.Auth("authenticate")).Patch("/reviews/{reviewId}/approve", helpers.Make(h.ApproveReview))
		r.With(middleware.Auth("authenticate")).Patch("/reviews/{reviewId}/reject", helpers.Make(h.RejectReview))
		r.With(middleware.Auth("authenticate")).Patch("/reviews/{reviewId}/status", helpers.Make(h.UpdateReviewStatus))

	})

}

func (l *HotelsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.hotelssvcs.GetOne(ctx, id)
	if err != nil {
		return err

	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Parse JSON filter parameter if present
	filterParam := r.URL.Query().Get("query")
	var filter filter.HotelsFilter
	if filterParam != "" {
		err := json.Unmarshal([]byte(filterParam), &filter)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}

	// Parse direct query parameters into filter
	if searchWord := r.URL.Query().Get("searchWord"); searchWord != "" {
		filter.SearchWord = searchWord
	}
	if location := r.URL.Query().Get("location"); location != "" {
		filter.Locations = []string{location}
	}

	// Use util.Paginate to get standardized pagination values
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := l.hotelssvcs.GetAll(ctx, filter, page, perPage)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	result, err := l.hotelssvcs.Add(ctx, &data)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.hotelssvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.hotelssvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) AddReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	hotelId := chi.URLParam(r, "id")

	var data models.HotelReviewDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	// Check if user is authenticated by getting user from context
	user, _ := r.Context().Value(util.ReqUser).(*types.User)
	if user == nil {
		// Unauthenticated user - use default userId and require form fields
		if data.FirstName == "" || data.LastName == "" || data.Email == "" {
			return helpers.BadRequest("firstName, lastName, and email are required for unauthenticated users")
		}
		data.UserId = "000000000000000000000000" // Default ObjectID for unauthenticated users
	} else {
		// Authenticated user - use their userId and clear personal info fields
		data.UserId = user.Id.Hex()
		data.FirstName = "" // Empty for authenticated users
		data.LastName = ""  // Empty for authenticated users
		data.Email = ""     // Empty for authenticated users
	}

	result, err := l.hotelssvcs.AddReview(ctx, hotelId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *HotelsHandler) GetHotelReviews(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	hotelId := chi.URLParam(r, "id")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	// Convert skip/limit to page/perPage for the repository
	page := int((skip / limit) + 1)
	perPage := int(limit)

	result, err := h.hotelssvcs.GetHotelReviews(ctx, hotelId, page, perPage)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := l.hotelssvcs.GetAllReviews(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) ApproveReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	reviewId := chi.URLParam(r, "reviewId")

	err := l.hotelssvcs.ApproveReview(ctx, reviewId)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) RejectReview(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	reviewId := chi.URLParam(r, "reviewId")

	err := l.hotelssvcs.RejectReview(ctx, reviewId)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) UpdateReviewStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	reviewId := chi.URLParam(r, "reviewId")

	var data models.ReviewStatusDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	err := l.hotelssvcs.UpdateReviewStatus(ctx, reviewId, data.Status)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
