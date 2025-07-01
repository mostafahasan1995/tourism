package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	transtest "larsa-tourism-microservices/pkg/services/trans-test"
	"larsa-tourism-microservices/pkg/services/trans-test/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type TransTestHandler struct {
	transtestsvcs transtest.TransTestSvcs
}

func NewTransTestHandler(i *do.Injector, r *chi.Mux) {
	h := &TransTestHandler{
		transtestsvcs: do.MustInvoke[transtest.TransTestSvcs](i),
	}

	r.Route("/trans-test", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/", helpers.Make(h.Get))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
	})
}

func (h *TransTestHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.transtestsvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *TransTestHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := h.transtestsvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).EncodeContext(ctx, result)
}

func (h *TransTestHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TransTestDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.transtestsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
