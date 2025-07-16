package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/marketing"
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type NewsletterHandler struct {
	newsletterSvcs marketing.NewsletterSvcs
}

func NewNewsletterHandler(i *do.Injector, r *chi.Mux) {
	h := &NewsletterHandler{
		newsletterSvcs: do.MustInvoke[marketing.NewsletterSvcs](i),
	}

	r.Route("/newsletters", func(r chi.Router) {
		r.With(middleware.OptionalAuth()).Post("/subscribe", helpers.Make(h.Subscribe))
		r.With(middleware.OptionalAuth()).Get("/unsubscribe", helpers.Make(h.Unsubscribe))
	})
}

func (h *NewsletterHandler) Subscribe(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.NewsletterDTO
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	if err := validator.New().Struct(data); err != nil {
		return helpers.InvalidJSON()
	}
	err := h.newsletterSvcs.Subscribe(ctx, data.Email)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{"message": "Newsletter subscribed successfully"})
}

// We can send unsubscribe link to the user in every email
// This is the endpoint that can be used to unsubscribe
// example: GET https://imkan.com/newsletters/unsubscribe?email=test@test.com
func (h *NewsletterHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	email := r.URL.Query().Get("email")
	if email == "" {
		return helpers.BadRequest("Email is required")
	}

	err := h.newsletterSvcs.Unsubscribe(ctx, email)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{"message": "Newsletter unsubscribed successfully"})
}
