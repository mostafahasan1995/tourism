package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaqHandler struct {
	faqPageSvcs     home.FaqPageSvcs
	faqGroupSvcs    home.FaqGroupSvcs
	faqQuestionSvcs home.FaqQuestionSvcs
}

func NewFaqHandler(i *do.Injector, r *chi.Mux) {
	h := &FaqHandler{
		faqPageSvcs:     do.MustInvoke[home.FaqPageSvcs](i),
		faqGroupSvcs:    do.MustInvoke[home.FaqGroupSvcs](i),
		faqQuestionSvcs: do.MustInvoke[home.FaqQuestionSvcs](i),
	}

	// FAQ Pages routes
	r.Route("/faq-pages", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetAllPages))
		r.Get("/{id}", helpers.Make(h.GetOnePage))
		r.Get("/{id}/complete", helpers.Make(h.GetPageComplete))
		r.Post("/initialize", helpers.Make(h.InitializeStaticPages))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddPage))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddManyPages))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.UpdatePage))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.PatchPage))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.DeletePage))

		//V2
		r.Route("/v2", func(r chi.Router) {
			r.Post("/", helpers.Make(h.GetPagesV2))
			r.Post("/all", helpers.Make(h.GetAllPagesV2))
		})
	})

	// FAQ Groups routes - nested under FAQ pages
	r.Route("/faq-pages/{faqPageId}/groups", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetAllGroups))
		r.Get("/{id}", helpers.Make(h.GetOneGroup))
		r.Get("/{id}/with-questions", helpers.Make(h.GetGroupWithQuestions))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddGroup))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddManyGroups))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.UpdateGroup))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.PatchGroup))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.DeleteGroup))

		// V2
		r.Route("/v2", func(r chi.Router) {
			r.Post("/", helpers.Make(h.GetGroupsV2))
			r.Post("/all", helpers.Make(h.GetAllGroupsV2))
		})
	})

	// FAQ Groups routes - simplified without page ID requirement
	r.Route("/faq-page-groups", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetAllGroups))
		r.Get("/{id}", helpers.Make(h.GetOneGroup))
		r.Get("/{id}/with-questions", helpers.Make(h.GetGroupWithQuestions))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddGroup))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddManyGroups))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.UpdateGroup))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.PatchGroup))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.DeleteGroup))
	})

	// FAQ Questions routes - simplified to use only group ID
	r.Route("/faq-groups/{faqGroupId}/questions", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetAllQuestions))
		r.Get("/{id}", helpers.Make(h.GetOneQuestion))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddQuestion))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddManyQuestions))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.UpdateQuestion))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.PatchQuestion))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.DeleteQuestion))

		// V2
		r.Route("/v2", func(r chi.Router) {
			r.Post("/", helpers.Make(h.GetQuestionsV2))
			r.Post("/all", helpers.Make(h.GetAllQuestionsV2))
		})
	})

	// FAQ Questions routes - general questions (no group)
	r.Route("/faq-questions", func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetAllQuestions))
		r.Get("/{id}", helpers.Make(h.GetOneQuestion))

		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.AddQuestion))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddManyQuestions))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.UpdateQuestion))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.PatchQuestion))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.DeleteQuestion))
	})

	// Debug route - temporary for troubleshooting
	r.Get("/debug/questions", helpers.Make(h.DebugAllQuestions))
}

// Helper function to parse pagination and filter parameters
func parsePageFilter(r *http.Request) filter.FaqPageFilter {
	var pageFilter filter.FaqPageFilter

	// Handle JSON query parameter
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		json.Unmarshal([]byte(filterParam), &pageFilter)
	}

	// Handle direct parameters
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			pageFilter.Page = page
		}
	}
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			pageFilter.Size = size
		}
	}
	if nameParam := r.URL.Query().Get("name"); nameParam != "" {
		pageFilter.Name = nameParam
	}
	if activeParam := r.URL.Query().Get("isActive"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			pageFilter.IsActive = &active
		}
	}

	return pageFilter
}

func parseGroupFilter(r *http.Request) filter.FaqGroupFilter {
	var groupFilter filter.FaqGroupFilter

	// Handle JSON query parameter
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		json.Unmarshal([]byte(filterParam), &groupFilter)
	}

	// Handle direct parameters
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			groupFilter.Page = page
		}
	}
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			groupFilter.Size = size
		}
	}
	if faqPageIdParam := r.URL.Query().Get("faqPageId"); faqPageIdParam != "" {
		if id, err := primitive.ObjectIDFromHex(faqPageIdParam); err == nil {
			groupFilter.FaqPageId = id
		}
	}
	if nameParam := r.URL.Query().Get("name"); nameParam != "" {
		groupFilter.Name = nameParam
	}
	if activeParam := r.URL.Query().Get("isActive"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			groupFilter.IsActive = &active
		}
	}

	return groupFilter
}

func parseQuestionFilter(r *http.Request) filter.FaqQuestionFilter {
	var questionFilter filter.FaqQuestionFilter

	// Handle JSON query parameter
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		json.Unmarshal([]byte(filterParam), &questionFilter)
	}

	// Handle direct parameters
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			questionFilter.Page = page
		}
	}
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			questionFilter.Size = size
		}
	}
	if faqPageIdParam := r.URL.Query().Get("faqPageId"); faqPageIdParam != "" {
		if id, err := primitive.ObjectIDFromHex(faqPageIdParam); err == nil {
			questionFilter.FaqPageId = id
		}
	}
	if faqGroupIdParam := r.URL.Query().Get("faqGroupId"); faqGroupIdParam != "" {
		if faqGroupIdParam == "null" || faqGroupIdParam == "" {
			// For general questions (no group)
			nilId := primitive.NilObjectID
			questionFilter.FaqGroupId = &nilId
		} else if id, err := primitive.ObjectIDFromHex(faqGroupIdParam); err == nil {
			questionFilter.FaqGroupId = &id
		}
	}
	if questionParam := r.URL.Query().Get("question"); questionParam != "" {
		questionFilter.Question = questionParam
	}
	if answerParam := r.URL.Query().Get("answer"); answerParam != "" {
		questionFilter.Answer = answerParam
	}
	if searchParam := r.URL.Query().Get("search"); searchParam != "" {
		questionFilter.Search = searchParam
	}
	if activeParam := r.URL.Query().Get("isActive"); activeParam != "" {
		if active, err := strconv.ParseBool(activeParam); err == nil {
			questionFilter.IsActive = &active
		}
	}

	return questionFilter
}

// =============================================================================
// FAQ PAGE HANDLERS
// =============================================================================

func (h *FaqHandler) GetAllPages(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	pageFilter := parsePageFilter(r)

	result, err := h.faqPageSvcs.GetAll(ctx, pageFilter)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetOnePage(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.faqPageSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetPageComplete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.faqPageSvcs.GetComplete(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) AddPage(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FaqPageDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqPageSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ page created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) AddManyPages(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.FaqPageDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqPageSvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ pages created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) UpdatePage(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.FaqPageDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqPageSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ page updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) PatchPage(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqPageSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ page updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) DeletePage(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.faqPageSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ page deleted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) InitializeStaticPages(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	err := h.faqPageSvcs.InitializeStaticPages(ctx)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Static FAQ pages initialized successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

// V2:

func (h *FaqHandler) GetPagesV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.faqPageSvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetAllPagesV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}
	result, err := h.faqPageSvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// =============================================================================
// FAQ GROUP HANDLERS
// =============================================================================

func (h *FaqHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	groupFilter := parseGroupFilter(r)

	// Check if faqPageId is provided in URL parameters
	faqPageIdParam := chi.URLParam(r, "faqPageId")
	includeQuestions := faqPageIdParam != ""

	// If faqPageId is provided in URL parameters, use it to filter
	if faqPageIdParam != "" {
		if faqPageId, err := primitive.ObjectIDFromHex(faqPageIdParam); err == nil {
			groupFilter.FaqPageId = faqPageId
		}
	}

	// If we need to include questions (when called from /faq-pages/{id}/groups route)
	if includeQuestions {
		// Get all groups first
		result, err := h.faqGroupSvcs.GetAll(ctx, groupFilter)
		if err != nil {
			return err
		}

		// Convert each group to include questions
		var groupsWithQuestions []models.FaqGroupWithQuestions
		for _, group := range result.FaqGroups {
			groupWithQuestions, err := h.faqGroupSvcs.GetWithQuestions(ctx, group.Id.Hex())
			if err != nil {
				// If error getting questions, include group without questions
				groupsWithQuestions = append(groupsWithQuestions, models.FaqGroupWithQuestions{
					FaqGroup:  group,
					Questions: []models.FaqQuestion{},
				})
			} else {
				groupsWithQuestions = append(groupsWithQuestions, *groupWithQuestions)
			}
		}

		// Return groups with questions
		response := map[string]interface{}{
			"faqGroups":  groupsWithQuestions,
			"pagination": result.Pagination,
		}
		return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
	}

	// Default behavior - return groups without questions
	result, err := h.faqGroupSvcs.GetAll(ctx, groupFilter)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetOneGroup(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.faqGroupSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetGroupWithQuestions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.faqGroupSvcs.GetWithQuestions(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) AddGroup(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FaqGroupDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	// Get faqPageId from URL parameters if provided (for nested routes)
	faqPageIdParam := chi.URLParam(r, "faqPageId")
	if faqPageIdParam != "" {
		faqPageId, err := primitive.ObjectIDFromHex(faqPageIdParam)
		if err != nil {
			return helpers.InvalidObjectId()
		}
		// Set the faqPageId from URL parameter
		data.FaqPageId = faqPageId
	}
	// If no faqPageId in URL, it should be provided in the request body

	err := h.faqGroupSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ group created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) AddManyGroups(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.FaqGroupDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	// Get faqPageId from URL parameters if provided (for nested routes)
	faqPageIdParam := chi.URLParam(r, "faqPageId")
	if faqPageIdParam != "" {
		faqPageId, err := primitive.ObjectIDFromHex(faqPageIdParam)
		if err != nil {
			return helpers.InvalidObjectId()
		}
		// Set the faqPageId from URL parameter for all groups
		for i := range data {
			data[i].FaqPageId = faqPageId
		}
	}
	// If no faqPageId in URL, it should be provided in the request body for each group

	err := h.faqGroupSvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ groups created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.FaqGroupDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqGroupSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ group updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) PatchGroup(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqGroupSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ group updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.faqGroupSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ group deleted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

// V2:

func (h *FaqHandler) GetGroupsV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.faqGroupSvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetAllGroupsV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}
	result, err := h.faqGroupSvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// =============================================================================
// FAQ QUESTION HANDLERS
// =============================================================================

func (h *FaqHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	questionFilter := parseQuestionFilter(r)

	// If faqGroupId is provided in URL parameters, use it to filter
	if faqGroupIdParam := chi.URLParam(r, "faqGroupId"); faqGroupIdParam != "" {
		if faqGroupId, err := primitive.ObjectIDFromHex(faqGroupIdParam); err == nil {
			questionFilter.FaqGroupId = &faqGroupId
		}
	}

	result, err := h.faqQuestionSvcs.GetAll(ctx, questionFilter)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetOneQuestion(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	result, err := h.faqQuestionSvcs.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) AddQuestion(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FaqQuestionDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	// Get faqGroupId from URL parameters if provided
	faqGroupIdParam := chi.URLParam(r, "faqGroupId")
	if faqGroupIdParam != "" {
		faqGroupId, err := primitive.ObjectIDFromHex(faqGroupIdParam)
		if err != nil {
			return helpers.InvalidObjectId()
		}
		data.FaqGroupId = &faqGroupId
	}

	err := h.faqQuestionSvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ question created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) AddManyQuestions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.FaqQuestionDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	// Get faqGroupId from URL parameters if provided
	faqGroupIdParam := chi.URLParam(r, "faqGroupId")
	if faqGroupIdParam != "" {
		faqGroupId, err := primitive.ObjectIDFromHex(faqGroupIdParam)
		if err != nil {
			return helpers.InvalidObjectId()
		}

		for i := range data {
			data[i].FaqGroupId = &faqGroupId
		}
	}

	err := h.faqQuestionSvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ questions created successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusCreated, response)
}

func (h *FaqHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var data models.FaqQuestionDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqQuestionSvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ question updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) PatchQuestion(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	err := h.faqQuestionSvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ question updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (h *FaqHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	err := h.faqQuestionSvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "FAQ question deleted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

// V2:

func (h *FaqHandler) GetQuestionsV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.faqQuestionSvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaqHandler) GetAllQuestionsV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}
	result, err := h.faqQuestionSvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// Debug route - temporary for troubleshooting
func (h *FaqHandler) DebugAllQuestions(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.faqQuestionSvcs.GetAll(ctx, filter.FaqQuestionFilter{})
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
