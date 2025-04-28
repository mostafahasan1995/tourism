package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type MysteryBoxHandler struct {
	mysteryboxsvcs interactions.MysteryBoxSvcs
}

func NewMysteryBoxHandler(i *do.Injector, r *chi.Mux) {
	h := &MysteryBoxHandler{
		mysteryboxsvcs: do.MustInvoke[interactions.MysteryBoxSvcs](i),
	}

	r.Route("/mystery-boxes", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Get("/active", helpers.Make(h.GetActive))
		r.With(middleware.Auth("authenticate")).Post("/{mysteryBoxId}/box/{boxId}", helpers.Make(h.OpenBox))
	})
}

func (h *MysteryBoxHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.MysteryBoxDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	result, err := h.mysteryboxsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *MysteryBoxHandler) GetActive(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.mysteryboxsvcs.GetMysteryBox(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *MysteryBoxHandler) OpenBox(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	mysteryBoxId := chi.URLParam(r, "mysteryBoxId")
	boxId := chi.URLParam(r, "boxId")

	result, err := h.mysteryboxsvcs.OpenBox(ctx, mysteryBoxId, boxId)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
