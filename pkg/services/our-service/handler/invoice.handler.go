package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourService "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type InvoiceHandler struct {
	invoicesvcs        ourService.InvoiceSvcs
	validationInstance *validator.Validate
}

func NewInvoiceHandler(i *do.Injector, r *chi.Mux) {
	h := &InvoiceHandler{
		invoicesvcs:        do.MustInvoke[ourService.InvoiceSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/invoices", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Post("/{id}/payments", helpers.Make(h.AddPayment))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/payments/{paymentId}", helpers.Make(h.UpdatePayment))
		r.With(middleware.Auth("authenticate")).Delete("/{id}/payments/{paymentId}", helpers.Make(h.DeletePayment))
		r.With(middleware.Auth("authenticate")).Post("/{id}/pay-order", helpers.Make(h.PayOrder))
	})
}

func (h *InvoiceHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.invoicesvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := h.invoicesvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.InvoiceDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.invoicesvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) AddPayment(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	invoiceId := chi.URLParam(r, "id")

	var data models.PaymentDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.invoicesvcs.AddPayment(ctx, invoiceId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	invoiceId := chi.URLParam(r, "id")
	paymentId := chi.URLParam(r, "paymentId")

	var data models.PaymentDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.invoicesvcs.UpdatePayment(ctx, invoiceId, paymentId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) DeletePayment(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	invoiceId := chi.URLParam(r, "id")
	paymentId := chi.URLParam(r, "paymentId")

	result, err := h.invoicesvcs.DeletePayment(ctx, invoiceId, paymentId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *InvoiceHandler) PayOrder(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	invoiceId := chi.URLParam(r, "id")

	var data models.PayOrder
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.invoicesvcs.PayOrder(ctx, invoiceId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
