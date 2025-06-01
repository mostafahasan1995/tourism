package handler

// type OurAgentsHandler struct {
// 	ourAgentssvcs home.OurAgentsSvcs
// }

// func NewOurAgentsHandler(i *do.Injector, r *chi.Mux) {
// 	h := &OurAgentsHandler{
// 		ourAgentssvcs: do.MustInvoke[home.OurAgentsSvcs](i),
// 	}

// 	r.Route("/ourAgents", func(r chi.Router) {

// 		r.Get("/{id}", helpers.Make(h.GetOne))

// 		r.Get("/", helpers.Make(h.GetAll))
// 		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
// 		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddMany))

// 		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))

// 		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))

// 	})

// }

// func (l *OurAgentsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	id := chi.URLParam(r, "id")

// 	result, err := l.ourAgentssvcs.GetOne(ctx, id)
// 	if err != nil {
// 		return err

// 	}
// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (l *OurAgentsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)
// 	filterParam := r.URL.Query().Get("query")
// 	var filter filter.OurAgentsFilter
// 	if filterParam != "" {
// 		err := json.Unmarshal([]byte(filterParam), &filter)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return err
// 		}
// 	}
// 	result, err := l.ourAgentssvcs.GetAll(ctx, filter)
// 	if err != nil {
// 		return err
// 	}
// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (l *OurAgentsHandler) Add(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data models.OurAgentsDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	//and validations go here

// 	err := l.ourAgentssvcs.Add(ctx, &data)
// 	if err != nil {
// 		return err
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	return nil
// }

// func (l *OurAgentsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)
// 	id := chi.URLParam(r, "id")
// 	var err = l.ourAgentssvcs.Delete(ctx, id)

// 	if err != nil {
// 		return err
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	return nil
// }

// func (l *OurAgentsHandler) Update(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data models.OurAgentsDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	//and validations go here
// 	id := chi.URLParam(r, "id")
// 	err := l.ourAgentssvcs.Update(ctx, id, &data)
// 	if err != nil {
// 		return err
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	return nil
// }

// func (l *OurAgentsHandler) AddMany(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data []models.OurAgentsDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	//and validations go here

// 	err := l.ourAgentssvcs.AddMany(ctx, data)
// 	if err != nil {
// 		return err
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	return nil
// }
