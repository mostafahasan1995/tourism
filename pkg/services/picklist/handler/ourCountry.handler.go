package handler

// type OurCountryHandler struct {
// 	ourCountrysvcs picklist.OurCountrySvcs
// }

// func NewOurCountryHandler(i *do.Injector, r *chi.Mux) {
// 	h := &OurCountryHandler{
// 		ourCountrysvcs: do.MustInvoke[picklist.OurCountrySvcs](i),
// 	}

// 	r.Route("/ourCountry", func(r chi.Router) {
// 		r.Get("/{id}", helpers.Make(h.GetOne))
// 		r.Get("/", helpers.Make(h.GetAll))
// 		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
// 		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
// 		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
// 	})
// }

// func (l *OurCountryHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	id := chi.URLParam(r, "id")

// 	result, err := l.ourCountrysvcs.GetOne(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (l *OurCountryHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)
// 	filterParam := r.URL.Query().Get("query")
// 	var f filter.OurCountryFilter
// 	if filterParam != "" {
// 		err := json.Unmarshal([]byte(filterParam), &f)
// 		if err != nil {
// 			return helpers.InvalidJSON()
// 		}
// 	}
// 	result, err := l.ourCountrysvcs.GetAll(ctx, f)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (l *OurCountryHandler) Add(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data models.OurCountryDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	country, err := l.ourCountrysvcs.Add(ctx, &data)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusCreated, country)
// }

// func (l *OurCountryHandler) Delete(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)
// 	id := chi.URLParam(r, "id")
// 	err := l.ourCountrysvcs.Delete(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, map[string]string{"message": "Country deleted successfully"})
// }

// func (l *OurCountryHandler) Update(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data models.OurCountryDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	id := chi.URLParam(r, "id")
// 	updatedCountry, err := l.ourCountrysvcs.Update(ctx, id, &data)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, updatedCountry)
// }
