package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/liteapi"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type HotelHandler struct {
	hotelSvcs liteapi.HotelSvcs
}

func NewHotelHandler(i *do.Injector, r *chi.Mux) {
	h := &HotelHandler{
		hotelSvcs: do.MustInvoke[liteapi.HotelSvcs](i),
	}

	r.Route("/liteapi", func(r chi.Router) {
		r.Get("/data/hotels", helpers.Make(h.GetHotels))
	})

}

func (h *HotelHandler) GetHotels(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	params := r.URL.Query()

	query := make(map[string]string)
	for key, value := range params {
		query[key] = value[0]
	}

	result, err := h.hotelSvcs.SearchHotels(ctx, query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
