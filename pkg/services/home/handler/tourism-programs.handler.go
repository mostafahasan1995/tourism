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

type TourismProgramHandler struct {
	tourismProgramsvcs home.TourismProgramSvcs
}

func NewTourismProgramHandler(i *do.Injector, r *chi.Mux) {
	h := &TourismProgramHandler{
		tourismProgramsvcs: do.MustInvoke[home.TourismProgramSvcs](i),
	}

	r.Route("/tourism-program", func(r chi.Router) {

		r.Get("/{id}", helpers.Make(h.GetOne))

		r.Get("/", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

	})

}

func (l *TourismProgramHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.tourismProgramsvcs.GetOne(ctx, id)
	if err != nil {
		return err

	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *TourismProgramHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	filterParam := r.URL.Query().Get("query")
	var filter filter.TourismProgramFilter
	if filterParam != "" {
		err := json.Unmarshal([]byte(filterParam), &filter)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}

	result, err := l.tourismProgramsvcs.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, result)
}

func (l *TourismProgramHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TourismProgramDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	err := l.tourismProgramsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
func (l *TourismProgramHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.tourismProgramsvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *TourismProgramHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.TourismProgramDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.tourismProgramsvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
