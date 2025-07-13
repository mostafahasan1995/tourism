package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type DestinationHandler struct {
	destinationSvcs    picklist.DestinationSvcs
	validationInstance *validator.Validate
}

func NewDestinationHandler(i *do.Injector, r *chi.Mux) {
	h := &DestinationHandler{
		destinationSvcs:    do.MustInvoke[picklist.DestinationSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/destinations", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
		r.With(middleware.Auth("authenticate")).Post("/country", helpers.Make(h.GetDistinationByCountry))
	})

	r.Route("/destinations/v2", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Post("/all", helpers.Make(h.GetAllV2))
	})
}

func (h *DestinationHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	destinationId := chi.URLParam(r, "id")

	result, err := h.destinationSvcs.GetOne(ctx, destinationId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *DestinationHandler) GetDistinationByCountry(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.DestinationCountry
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.destinationSvcs.GetDestinationByCountry(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *DestinationHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.destinationSvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *DestinationHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.DestinationDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.destinationSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *DestinationHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	destinationId := chi.URLParam(r, "id")

	var data models.DestinationDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.destinationSvcs.Update(ctx, destinationId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}

func (h *DestinationHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	destinationId := chi.URLParam(r, "id")

	if err := h.destinationSvcs.Delete(ctx, destinationId); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, "ok")
}

// v2
func (h *DestinationHandler) GetAllV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.destinationSvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, result)
}
