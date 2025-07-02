package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type CarsHandler struct {
	carssvcs picklist.CarsSvcs
}

func NewCarsHandler(i *do.Injector, r *chi.Mux) {
	h := &CarsHandler{
		carssvcs: do.MustInvoke[picklist.CarsSvcs](i),
	}

	r.Route("/cars", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/", helpers.Make(h.GetPaginated))
		r.Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})
}

func (l *CarsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.carssvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *CarsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	result, err := l.carssvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)

}

func (l *CarsHandler) GetPaginated(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	result, err := l.carssvcs.GetPaginated(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *CarsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CarsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	// TODO: Add validation if needed

	car, err := l.carssvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, car)
}

func (l *CarsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := l.carssvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{"message": "Car deleted successfully"})
}

func (l *CarsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.CarsDto
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return helpers.InvalidJSON()
	}

	// TODO: Add validation if needed
	id := chi.URLParam(r, "id")

	updatedCar, err := l.carssvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, updatedCar)
}
