package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type FaveHandler struct {
	favesvcs interactions.FaveSvcs
}

func NewFaveHandler(i *do.Injector, r *chi.Mux) {
	h := &FaveHandler{
		favesvcs: do.MustInvoke[interactions.FaveSvcs](i),
	}

	r.Route("/favorites", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
		r.With(middleware.Auth("authenticate")).Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.Get))
		// r.With(middleware.Auth("authenticate")).Get("/{type}/all", helpers.Make(h.GetAllByType))
	})
}

// func (h *FaveHandler) GetAllByType(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	t := chi.URLParam(r, "type")

// 	result, err := h.favesvcs.GetAllByType(ctx, models.FaveType(t))
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

func (h *FaveHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	var data models.FaveDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.favesvcs.Patch(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *FaveHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	err := h.favesvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Favorite deleted successfully",
	})
}

func (h *FaveHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.favesvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *FaveHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")
	skipStr := r.URL.Query().Get("skip")
	limitStr := r.URL.Query().Get("limit")

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		skip = 0
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}

	result, err := h.favesvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *FaveHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FaveDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.favesvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
