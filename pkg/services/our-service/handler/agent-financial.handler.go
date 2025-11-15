package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	ourservice "larsa-tourism-microservices/pkg/services/our-service"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type AgentFinancialHandler struct {
	svcs               ourservice.AgentFinancialSvcs
	validationInstance *validator.Validate
}

func NewAgentFinancialHandler(i *do.Injector, r *chi.Mux) {
	h := &AgentFinancialHandler{
		svcs:               do.MustInvoke[ourservice.AgentFinancialSvcs](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route("/agent-financial", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/{agentId}", helpers.Make(h.GetAccount))
		r.With(middleware.Auth("authenticate")).Post("/{agentId}", helpers.Make(h.CreateOrUpdateAccount))
		r.With(middleware.Auth("authenticate")).Put("/{agentId}", helpers.Make(h.CreateOrUpdateAccount))
		r.With(middleware.Auth("authenticate")).Post("/{agentId}/withdraw", helpers.Make(h.Withdraw))
		r.With(middleware.Auth("authenticate")).Patch("/{agentId}/withdrawals/{withdrawalId}/approve", helpers.Make(h.ApproveWithdrawal))
	})
}

func (h *AgentFinancialHandler) GetAccount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "agentId")

	// Parse limit from query parameter (default: 10, max: 100)
	limit := 10
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			} else if parsedLimit > 100 {
				limit = 100 // Cap at 100
			}
		}
	}

	account, err := h.svcs.GetAccount(ctx, agentId, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, account)
}

func (h *AgentFinancialHandler) Withdraw(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "agentId")

	var req models.AgentWithdrawRequest
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); err != nil {
		return err
	}

	// Parse limit from query parameter (default: 10, max: 100)
	limit := 10
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			if parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			} else if parsedLimit > 100 {
				limit = 100 // Cap at 100
			}
		}
	}

	account, err := h.svcs.Withdraw(ctx, agentId, &req, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, account)
}

func (h *AgentFinancialHandler) CreateOrUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "agentId")

	var dto models.AgentFinancialAccountDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &dto); err != nil {
		return err
	}

	if err := dto.Validate(h.validationInstance); err != nil {
		return err
	}

	account, err := h.svcs.CreateOrUpdateAccount(ctx, agentId, &dto)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, account)
}

func (h *AgentFinancialHandler) ApproveWithdrawal(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	agentId := chi.URLParam(r, "agentId")
	withdrawalId := chi.URLParam(r, "withdrawalId")

	account, err := h.svcs.ApproveWithdrawal(ctx, agentId, withdrawalId)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, account)
}
