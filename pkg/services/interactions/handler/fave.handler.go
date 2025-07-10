package handler

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/services/interactions"
	"larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"

	"github.com/goccy/go-json"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		r.With(middleware.Auth("authenticate")).Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
		r.With(middleware.Auth("authenticate")).Get("/all", helpers.Make(h.GetAll))
		r.With(middleware.Auth("authenticate")).Get("/", helpers.Make(h.Get))
		// r.With(middleware.Auth("authenticate")).Get("/{type}/all", helpers.Make(h.GetAllByType))
	})
}

// func (h *FaveHandler) GetAllByType(w http.ResponseWriter, r *http.Request) error {
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	t := chi.URLParam(r, "type")

// 	result, err := h.favesvcs.GetAllByType(ctx, models.FaveType(t))
// 	if err != nil {
// 		return err
// 	}

// 	return helpers.WriteJson(w, http.StatusOK, result)
// }

func (h *FaveHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	// Create a temporary struct to handle string refId from JSON
	var rawData struct {
		Type  string `json:"type"`
		RefId string `json:"refId"`
		IsFav bool   `json:"isFav"`
	}

	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &rawData); err != nil {
		return err
	}

	// Create the proper FaveDto with converted ObjectID if refId is provided
	data := &models.FaveDto{
		IsFav: rawData.IsFav,
	}

	if rawData.Type != "" {
		data.Type = models.FaveType(rawData.Type)
	}

	if rawData.RefId != "" {
		objId, err := primitive.ObjectIDFromHex(rawData.RefId)
		if err != nil {
			return helpers.BadRequest("Invalid refId format")
		}
		data.RefId = objId
	}

	result, err := h.favesvcs.Patch(ctx, id, data)
	if err != nil {
		return err
	}

	// // Update the entity's isFav field if both type and refId are provided
	// if data.Type != "" && !data.RefId.IsZero() {
	// 	h.updateEntityIsFav(ctx, data.Type, data.RefId, data.IsFav)
	// }

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaveHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	err := h.favesvcs.Delete(ctx, id)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{
		"message": "Favorite deleted successfully",
	})
}

func (h *FaveHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, err := h.favesvcs.GetAll(ctx, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (h *FaveHandler) Get(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	query := r.URL.Query().Get("query")

	result, err := h.favesvcs.Get(ctx, skip, limit, query)
	if err != nil {
		return err
	}

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)

}

// // Helper method to update entity isFav field based on favorite type
// func (h *FaveHandler) updateEntityIsFav(ctx context.Context, faveType models.FaveType, refId primitive.ObjectID, isFav bool) {
// 	// Update entity asynchronously to avoid blocking the response
// 	// If it fails, it's not critical as the favorite is still saved
// 	go func() {
// 		cfg, err := util.GetReqAppCfg(ctx)
// 		if err != nil {
// 			return
// 		}

// 		// Get the collection name based on favorite type
// 		var collectionName string
// 		switch faveType {
// 		case models.FaveTypeProgram:
// 			collectionName = "tourismPrograms"
// 		case models.FaveTypeHotel:
// 			collectionName = "tourismHotels"
// 		case models.FaveTypeDestination:
// 			collectionName = "tourismDestinations"
// 		case models.FaveTypeExhibition:
// 			collectionName = "tourismExhibitions"
// 		case models.FaveTypeDiary:
// 			collectionName = "tourismDiaries"
// 		case models.FaveTypeAgent:
// 			collectionName = "tourismAgents"
// 		default:
// 			return // Unknown type, skip update
// 		}

// 		// Update the entity's isFav field directly in the database
// 		collection := h.db.Database(cfg.Db).Collection(collectionName)
// 		filter := bson.M{"_id": refId}
// 		update := bson.M{
// 			"$set": bson.M{
// 				"isFav":     isFav,
// 				"updatedAt": time.Now(),
// 			},
// 		}

// 		// Perform the update - don't worry about errors since this is supplementary
// 		_, _ = collection.UpdateOne(ctx, filter, update)
// 	}()
// }

func (h *FaveHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Create a temporary struct to handle string refId from JSON
	var rawData struct {
		Type  string `json:"type"`
		RefId string `json:"refId"`
		IsFav bool   `json:"isFav"`
	}

	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &rawData); err != nil {
		return err
	}

	// Convert string refId to ObjectID
	objId, err := primitive.ObjectIDFromHex(rawData.RefId)
	if err != nil {
		return helpers.BadRequest("Invalid refId format")
	}

	// Create the proper FaveDto with converted ObjectID
	data := &models.FaveDto{
		Type:  models.FaveType(rawData.Type),
		RefId: objId,
		IsFav: rawData.IsFav,
	}

	result, err := h.favesvcs.Add(ctx, data)
	if err != nil {
		return err
	}

	// Update the entity's isFav field
	// h.updateEntityIsFav(ctx, data.Type, data.RefId, data.IsFav)

	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}
