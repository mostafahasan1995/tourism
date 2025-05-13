package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	customform "larsa-tourism-microservices/pkg/services/custom-form"
	"larsa-tourism-microservices/pkg/services/custom-form/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

// deprecated - moved to travel-req service
type CustomFormHandler struct {
	cfsvcs customform.CustomFormSvcs
}

func NewCustomFormHandler(i *do.Injector, r *chi.Mux) {
	h := &CustomFormHandler{
		cfsvcs: do.MustInvoke[customform.CustomFormSvcs](i),
	}

	r.Route("/custom-forms", func(r chi.Router) {
		r.Post("/", helpers.Make(h.Add))
	})
}

func (h *CustomFormHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CustomFormDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.cfsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
