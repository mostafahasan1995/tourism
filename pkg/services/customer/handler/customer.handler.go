package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/customer"
	"larsa-tourism-microservices/pkg/services/customer/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type CustomerHandler struct {
	customersvcs customer.CustomerSvcs
}

func NewCustomerHandler(i *do.Injector, r *chi.Mux) {
	h := &CustomerHandler{
		customersvcs: do.MustInvoke[customer.CustomerSvcs](i),
	}

	r.Route("/customers", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/{id}", helpers.Make(h.GetOne))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.Get))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (h *CustomerHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	customerId := chi.URLParam(r, "id")

	result, err := h.customersvcs.GetOne(ctx, customerId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.customersvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *CustomerHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CustomerDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.customersvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	customerId := chi.URLParam(r, "id")

	var data models.CustomerDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.customersvcs.Update(ctx, customerId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	customerId := chi.URLParam(r, "id")

	if err := h.customersvcs.Delete(ctx, customerId); err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, "ok")
}
