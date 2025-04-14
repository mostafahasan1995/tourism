package handler

import (
	"encoding/json"
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type VipCarRequestHandler struct {
	vipCarRequestsvcs home.VipCarRequestSvcs
}

func NewVipCarRequestHandler(i *do.Injector, r *chi.Mux) {
	h := &VipCarRequestHandler{
		vipCarRequestsvcs: do.MustInvoke[home.VipCarRequestSvcs](i),
	}

	r.Route("/vipCarRequest", func(r chi.Router) {

		r.Get("/{id}", helpers.Make(h.GetOne))

		r.Get("/", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddMany))
		
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

	})

}

func (l *VipCarRequestHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.vipCarRequestsvcs.GetOne(ctx, id)
	if err != nil {
		return err

	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *VipCarRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	filterParam := r.URL.Query().Get("query")
	var filter filter.VipCarRequestFilter
	if filterParam != "" {
		err := json.Unmarshal([]byte(filterParam), &filter)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}
	result, err := l.vipCarRequestsvcs.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *VipCarRequestHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.VipCarRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	err := l.vipCarRequestsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *VipCarRequestHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.vipCarRequestsvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *VipCarRequestHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.VipCarRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.vipCarRequestsvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *VipCarRequestHandler) AddMany(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.VipCarRequestDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	err := l.vipCarRequestsvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}