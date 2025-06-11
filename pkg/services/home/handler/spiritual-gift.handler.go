package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type SpiritualGiftHandler struct {
	spiritualGiftSvcs  home.SpiritualGiftSvcs
	validationInstance *validator.Validate
}

func NewSpiritualGiftHandler(i *do.Injector, r *chi.Mux) {
	h := &SpiritualGiftHandler{
		spiritualGiftSvcs:  do.MustInvoke[home.SpiritualGiftSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/spiritual-gifts", func(r chi.Router) {
		r.Get("/", helpers.Make(h.Get))
		r.Get("/current", helpers.Make(h.GetCurrent))
		r.Get("/{id}", helpers.Make(h.GetById))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Save))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/toggle", helpers.Make(h.Toggle))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (h *SpiritualGiftHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")
	result, err := h.spiritualGiftSvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *SpiritualGiftHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.spiritualGiftSvcs.GetById(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *SpiritualGiftHandler) GetCurrent(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.spiritualGiftSvcs.GetCurrent(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
func (h *SpiritualGiftHandler) Save(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.SpiritualGiftDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	// Validate the data
	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	// If enabled is true, make sure required fields are present
	if data.Enabled {
		if data.PackageId.IsZero() || data.ProgramTitle == "" || data.ProgramId == "" {
			return helpers.BadRequest("When enabled, package ID, program title, and program ID are required")
		}
	}

	result, err := h.spiritualGiftSvcs.Save(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusCreated, result)
}

func (h *SpiritualGiftHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.SpiritualGiftDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	if data.Enabled {
		if data.PackageId.IsZero() || data.ProgramTitle == "" || data.ProgramId == "" {
			return helpers.BadRequest("When enabled, package ID, program title, and program ID are required")
		}
	}

	result, err := h.spiritualGiftSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *SpiritualGiftHandler) Toggle(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.BadRequest("Invalid JSON format")
	}

	if data.Enabled {
		gift, err := h.spiritualGiftSvcs.GetById(ctx, id)
		if err != nil {
			return err
		}

		if gift.PackageId.IsZero() || gift.ProgramTitle == "" || gift.ProgramId == "" {
			return helpers.BadRequest("Cannot enable: package ID, program title, and program ID must be set")
		}
	}

	result, err := h.spiritualGiftSvcs.Toggle(ctx, id, data.Enabled)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *SpiritualGiftHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	if err := h.spiritualGiftSvcs.Delete(ctx, id); err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Spiritual gift configuration deleted successfully",
	})
}
