package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/request"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type TravelReqHandler struct {
	travelreqsvcs request.TravelReqSvcs
}

func NewTravelReqHandler(i *do.Injector, r *chi.Mux) {
	h := &TravelReqHandler{
		travelreqsvcs: do.MustInvoke[request.TravelReqSvcs](i),
	}

	r.Route("/requests", func(r chi.Router) {
		r.Post("/{svcstype}", helpers.Make(h.Add))
	})
}

func (h *TravelReqHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	svcstype := chi.URLParam(r, "svcstype")

	var data json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Add(ctx, svcstype, data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)

}
