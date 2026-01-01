package handler

// type DiaryHandler struct {
// 	diarysvcs interactions.DiarySvcs
// }

// func NewDiaryHandler(i *do.Injector, r *chi.Mux) {

// 	h := &DiaryHandler{
// 		diarysvcs: do.MustInvoke[interactions.DiarySvcs](i),
// 	}

// 	r.Route("/diaries", func(r chi.Router) {
// 		r.Get("/all", helpers.Make(h.GetAll))
// 		r.Get("/", helpers.Make(h.Get))
// 		r.Get("/{id}", helpers.Make(h.GetOne))
// 		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
// 	})

// }

// func (h *DiaryHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	query := r.URL.Query().Get("query")

// 	result, err := h.diarysvcs.GetAll(ctx, query)
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (h *DiaryHandler) Get(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	skip, limit, errGetPaginate := util.Paginate(r)
// 	if errGetPaginate != nil {
// 		return errGetPaginate
// 	}

// 	query := r.URL.Query().Get("query")

// 	result, err := h.diarysvcs.Get(ctx, skip, limit, query)
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (h *DiaryHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	id := chi.URLParam(r, "id")

// 	result, err := h.diarysvcs.GetOne(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

// func (h *DiaryHandler) Add(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var data models.DiaryDto
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return helpers.InvalidJSON()
// 	}

// 	result, err := h.diarysvcs.Add(ctx, &data)
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }
