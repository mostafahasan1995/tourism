package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type FaveHandler struct {
	favesvcs interactions.FaveSvcs
	db       *mongo.Client
}

func NewFaveHandler(i *do.Injector, r *chi.Mux) {
	h := &FaveHandler{
		favesvcs: do.MustInvoke[interactions.FaveSvcs](i),
		db:       do.MustInvoke[*mongo.Client](i),
	}

	r.Route("/favorites", func(r chi.Router) {
		r.With(middleware.Auth("authenticate")).Get("/{type}/all", helpers.Make(h.GetAllByType))
		//V2
		r.With(middleware.Auth("authenticate")).Post("/v2", helpers.Make(h.Fav))
	})
}

func (h *FaveHandler) Fav(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data models.FaveDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return err
	}

	if err := validator.New().Struct(data); err != nil {
		return err
	}

	message, err := h.favesvcs.Fav(ctx, &data)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{
		"message": message,
	})
}

func (h *FaveHandler) GetAllByType(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	typeStr := chi.URLParam(r, "type")
	if typeStr == "" {
		return helpers.BadRequest("Missing 'type' parameter")
	}

	switch typeStr {
	case string(models.FaveTypeProgram),
		string(models.FaveTypeHotel),
		string(models.FaveTypeTravelerStory),
		string(models.FaveTypeClientStory),
		string(models.FaveTypeExhibition),
		string(models.FaveTypeAgent),
		string(models.FaveTypeDestination):
		// valid
	default:
		return helpers.BadRequest("Invalid 'type' value")
	}

	faveType := models.FaveType(typeStr)
	result, err := h.favesvcs.GetAllByType(ctx, faveType)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
