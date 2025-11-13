package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type TravelRequestHandler struct {
	travelreqsvcs      ourservice.TravelRequestSvcs
	validationInstance *validator.Validate
}

func NewTravelRequestHandler(i *do.Injector, r *chi.Mux) {
	h := &TravelRequestHandler{
		travelreqsvcs:      do.MustInvoke[ourservice.TravelRequestSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/travel-requests", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.With(
			middleware.Auth("authenticate", "tourismGetTravelRequests"),
			middleware.CapabilityCheck("tourismGetOtherTravelRequests"),
		).Get("/", helpers.Make(h.Get))
		r.With(
			middleware.Auth("authenticate", "tourismGetTravelRequests"),
			middleware.CapabilityCheck("tourismGetOtherTravelRequests"),
		).Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Get("/customer/{customerId}", helpers.Make(h.GetCustomerRequests))
		r.With(middleware.Auth("authenticate")).Get("/agent/{agentId}", helpers.Make(h.GetAgentTransactions))
		r.With(middleware.Auth("authenticate")).Get("/profit", helpers.Make(h.GetCompanyTransactions))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Get("/my-requests/{status}", helpers.Make(h.MyRequests))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/approve", helpers.Make(h.Approve))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/reject", helpers.Make(h.Reject))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/complete", helpers.Make(h.SetAsCompleted))
	})

	r.Route("/travel-requests/v2", func(r chi.Router) {
		r.With(
			middleware.Auth("authenticate", "tourismGetTravelRequests"),
			middleware.CapabilityCheck("tourismGetOtherTravelRequests")).Post("/", helpers.Make(h.GetV2))
		r.With(
			middleware.Auth("authenticate", "tourismGetTravelRequests"),
			middleware.CapabilityCheck("tourismGetOtherTravelRequests")).Post("/all", helpers.Make(h.GetAllV2))
	})
}

func (h *TravelRequestHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := h.travelreqsvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) GetCustomerRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	customerId := chi.URLParam(r, "customerId")
	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.GetCustomerRequests(ctx, customerId, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) GetAgentTransactions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	agentId := chi.URLParam(r, "agentId")
	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.GetAgentTransactions(ctx, agentId, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) GetCompanyTransactions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.travelreqsvcs.GetCompanyTransactions(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TravelRequestDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	var data models.TravelRequestDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) MyRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	status := chi.URLParam(r, "status")

	result, err := h.travelreqsvcs.MyRequests(ctx, status)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Approve(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // travel request id

	result, err := h.travelreqsvcs.Approve(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) Reject(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // travel request id

	var data models.RejectMyReq
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Reject(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) SetAsCompleted(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") // travel request id

	result, err := h.travelreqsvcs.SetAsCompleted(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// v2
func (h *TravelRequestHandler) GetV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TravelRequestHandler) GetAllV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
