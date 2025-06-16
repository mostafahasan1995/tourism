package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type SettingsHandler struct {
	settingssvcs       ourservice.SettingsSvcs
	validationInstance *validator.Validate
}

func NewSettingsHandler(i *do.Injector, r *chi.Mux) {
	h := &SettingsHandler{
		settingssvcs:       do.MustInvoke[ourservice.SettingsSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/settings", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Put("/", helpers.Make(h.Update))
	})
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.settingssvcs.Get(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.SettingsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.settingssvcs.Update(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
