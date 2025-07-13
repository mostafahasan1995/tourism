package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/member"
	"larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type AgentHandler struct {
	agentsvcs          member.AgentSvcs
	validationInstance *validator.Validate
}

func NewAgentHandler(i *do.Injector, r *chi.Mux) {
	h := &AgentHandler{
		agentsvcs:          do.MustInvoke[member.AgentSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/agents", func(r chi.Router) {
		r.Get("/all", helpers.Make(h.GetAll))
		r.Get("/", helpers.Make(h.Get))
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/destinations", helpers.Make(h.GetDestinationAgents))
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/status", helpers.Make(h.UpdateStatus))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
	})

	r.Route("/agents/v2", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Post("/", helpers.Make(h.GetV2))
		r.With(middleware.OptionalAuth()).Post("/all", helpers.Make(h.GetAllV2))
	})

	r.Route("/agent-joins", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/{id}", helpers.Make(h.GetOneAgentJoin))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.GetJoinRequests))
		r.Post("/", helpers.Make(h.Join))
		r.With(middleware.Auth("authenticate")).Post("/{id}/convert", helpers.Make(h.ConvertToAgent))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/reject", helpers.Make(h.RejectJoin))
		r.With(middleware.Auth("authenticate")).Patch("/{id}/pending", helpers.Make(h.SetAsPending))
	})

	r.Route("/agent-joins/v2", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.GetJoinRequestsV2))
	})
}

func (h *AgentHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id")

	result, err := h.agentsvcs.GetOne(ctx, agentId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.agentsvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.agentsvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.AgentDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.agentsvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id")

	var data models.AgentDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.agentsvcs.Update(ctx, agentId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id")

	var data models.UpdateStatusDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	if err := data.Validate(h.validationInstance); err != nil {
		return err
	}

	result, err := h.agentsvcs.UpdateStatus(ctx, agentId, data.Status)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id")

	if err := h.agentsvcs.Delete(ctx, agentId); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

// agent join

func (h *AgentHandler) GetOneAgentJoin(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id") //agent join id

	result, err := h.agentsvcs.GetOneAgentJoin(ctx, agentId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) GetJoinRequests(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.agentsvcs.GetJoinRequests(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) Join(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.AgentJoinDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.agentsvcs.Join(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) ConvertToAgent(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id") // agent join id

	var data models.AgentJoinDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	result, err := h.agentsvcs.ConvertToAgent(ctx, agentId, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) RejectJoin(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id") // agent join id

	if err := h.agentsvcs.RejectJoin(ctx, agentId); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

func (h *AgentHandler) SetAsPending(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "id") // agent join id

	if err := h.agentsvcs.SetAsPending(ctx, agentId); err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, "ok")
}

// test
func (h *AgentHandler) GetDestinationAgents(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.agentsvcs.GetDestinationAgents(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

// v2
func (h *AgentHandler) GetV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.agentsvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) GetAllV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.agentsvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *AgentHandler) GetJoinRequestsV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := h.agentsvcs.GetJoinRequestsV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
