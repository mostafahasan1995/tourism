package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type FinancialSettingsHandler struct {
	financialsettingssvcs ourservice.FinancialSettingsSvcs
	validationInstance    *validator.Validate
}

func NewFinancialSettingsHandler(i *do.Injector, r *chi.Mux) {
	h := &FinancialSettingsHandler{
		financialsettingssvcs: do.MustInvoke[ourservice.FinancialSettingsSvcs](i),
		validationInstance:    do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/financial-settings", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Put("/", helpers.Make(h.Update))
	})
}

func (h *FinancialSettingsHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.financialsettingssvcs.Get(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *FinancialSettingsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FinancialSettingsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.financialsettingssvcs.Update(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
