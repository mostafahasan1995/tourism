package handler

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	travelreq "larsa-tourism-microservices/pkg/services/travel-req"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type TravelReqHandler struct {
	travelreqsvcs travelreq.TravelReqSvcs
}

func NewTravelReqHandler(i *do.Injector, r *chi.Mux) {
	h := &TravelReqHandler{
		travelreqsvcs: do.MustInvoke[travelreq.TravelReqSvcs](i),
	}

	r.Route("/travel-requests", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/{id}", helpers.Make(h.GetRelatedReq))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.GetTravelReqs))
		r.With(middleware.Auth("authenticate")).Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/{reqtype}", helpers.Make(h.Add))

	})
}

func (h *TravelReqHandler) GetRelatedReq(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id") //travel request id

	result, err := h.travelreqsvcs.GetRelatedReq(ctx, id)

	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelReqHandler) GetTravelReqs(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, err := h.travelreqsvcs.Get(ctx, skip, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelReqHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.travelreqsvcs.GetAll(ctx)

	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TravelReqHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	reqType := chi.URLParam(r, "reqtype") // request type

	var data json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	result, err := h.travelreqsvcs.Add(ctx, reqType, data)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
