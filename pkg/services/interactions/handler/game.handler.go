package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type GameHandler struct {
	gameSvcs interactions.GameSvcs
}

func NewGameHandler(i *do.Injector, r *chi.Mux) {
	h := &GameHandler{
		gameSvcs: do.MustInvoke[interactions.GameSvcs](i),
	}

	r.Route("/game", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetGame))
		r.Patch("/box/{boxId}", helpers.Make(h.UpdateBox))
		r.Patch("/settings", helpers.Make(h.UpdateSettings))
	})
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.gameSvcs.GetGame(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *GameHandler) UpdateBox(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "boxId")

	var data models.MysteryBox
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.gameSvcs.UpdateBox(ctx, id, data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *GameHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.Attempts
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.gameSvcs.UpdateSettings(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
