package handler

import (
	"encoding/json"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourService "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type HotelsHandler struct {
	hotelssvcs ourService.HotelsSvcs
}

func NewHotelsHandler(i *do.Injector, r *chi.Mux) {
	h := &HotelsHandler{
		hotelssvcs: do.MustInvoke[ourService.HotelsSvcs](i),
	}

	r.Route("/hotels", func(r chi.Router) {

		r.Get("/{id}", helpers.Make(h.GetOne))

		r.Get("/", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

	})

}

func (l *HotelsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.hotelssvcs.GetOne(ctx, id)
	if err != nil {
		return err

	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	filterParam := r.URL.Query().Get("query")
	var filter filter.HotelsFilter
	if filterParam != "" {
		err := json.Unmarshal([]byte(filterParam), &filter)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}
	result, err := l.hotelssvcs.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *HotelsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	err := l.hotelssvcs.Add(ctx, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.hotelssvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *HotelsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.HotelsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.hotelssvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
