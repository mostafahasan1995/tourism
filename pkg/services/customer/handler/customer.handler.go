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
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

	})
}

func (h *CustomerHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CustomerDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	customer, err := h.customersvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	customerId := chi.URLParam(r, "id")

	var data models.CustomerDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	customer, err := h.customersvcs.Update(ctx, customerId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, customer)
}
