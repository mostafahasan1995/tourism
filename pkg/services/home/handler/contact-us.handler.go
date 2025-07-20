package handler

import (
	"fmt"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/middleware"
	"larsa-tourism-microservices/pkg/query"
	"larsa-tourism-microservices/pkg/services/home"
	"larsa-tourism-microservices/pkg/services/home/filter"
	"larsa-tourism-microservices/pkg/services/home/models"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type ContactUsHandler struct {
	contactUssvcs home.ContactUsSvcs
}

// parseContactUsData parses raw JSON data into ContactUsDto with flexible field handling
func (l *ContactUsHandler) parseContactUsData(rawData map[string]interface{}) models.ContactUsDto {
	data := models.ContactUsDto{
		AdditionalFields: make(map[string]interface{}),
	}

	// Extract known fields
	if val, ok := rawData["fullName"]; ok {
		if str, ok := val.(string); ok {
			data.FullName = str
		}
	}
	if val, ok := rawData["emailAddress"]; ok {
		if str, ok := val.(string); ok {
			data.EmailAddress = str
		}
	}
	if val, ok := rawData["phoneNumber"]; ok {
		if str, ok := val.(string); ok {
			data.PhoneNumber = str
		}
	}
	if val, ok := rawData["howDidYouFindUs"]; ok {
		if str, ok := val.(string); ok {
			data.HowDidYouFindUs = str
		}
	}
	if val, ok := rawData["message"]; ok {
		if str, ok := val.(string); ok {
			data.Message = str
		}
	}

	// Handle additionalFields if provided as a nested object
	if val, ok := rawData["additionalFields"]; ok {
		if additionalMap, ok := val.(map[string]interface{}); ok {
			for k, v := range additionalMap {
				data.AdditionalFields[k] = v
			}
		}
	}

	// Move any unknown fields to additionalFields
	knownFields := map[string]bool{
		"fullName":         true,
		"emailAddress":     true,
		"phoneNumber":      true,
		"howDidYouFindUs":  true,
		"message":          true,
		"additionalFields": true,
	}

	for key, value := range rawData {
		if !knownFields[key] {
			data.AdditionalFields[key] = value
		}
	}

	return data
}

func NewContactUsHandler(i *do.Injector, r *chi.Mux) {
	h := &ContactUsHandler{
		contactUssvcs: do.MustInvoke[home.ContactUsSvcs](i),
	}

	r.Route("/contactUs", func(r chi.Router) {
		r.Get("/{id}", helpers.Make(h.GetOne))
		r.Get("/", helpers.Make(h.GetAll))
		r.Post("/", helpers.Make(h.Add))
		r.With(middleware.Auth("authenticate")).Post("/many", helpers.Make(h.AddMany))
		r.With(middleware.Auth("authenticate")).Put("/{id}", helpers.Make(h.Update))
		r.With(middleware.Auth("authenticate")).Patch("/{id}", helpers.Make(h.Patch))
		r.With(middleware.Auth("authenticate")).Delete("/{id}", helpers.Make(h.Delete))
		r.Route("/settings", func(r chi.Router) {
			r.Get("/", helpers.Make(h.GetSettings))
			r.Post("/", helpers.Make(h.AddOrUpdateSettings))
		})
		// V2
		r.Route("/v2", func(r chi.Router) {
			r.Post("/", helpers.Make(h.GetV2))
			r.Post("/all", helpers.Make(h.GetAllV2))
		})
	})

}

func (l *ContactUsHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")

	result, err := l.contactUssvcs.GetOne(ctx, id)
	if err != nil {
		return err

	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *ContactUsHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var filter filter.ContactUsFilter

	// Handle JSON query parameter (existing format)
	filterParam := r.URL.Query().Get("query")
	if filterParam != "" {
		err := json.Unmarshal([]byte(filterParam), &filter)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}

	// Handle direct pagination parameters (new format)
	pageParam := r.URL.Query().Get("page")
	sizeParam := r.URL.Query().Get("size")
	statusParam := r.URL.Query().Get("status")

	if pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if sizeParam != "" {
		if size, err := strconv.Atoi(sizeParam); err == nil && size > 0 {
			filter.Size = size
		}
	}

	if statusParam != "" {
		filter.Status = statusParam
	}

	result, err := l.contactUssvcs.GetAll(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *ContactUsHandler) Add(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// First decode into a generic map to capture all fields
	var rawData map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &rawData); err != nil {
		return helpers.InvalidJSON()
	}

	// Parse the data using the helper function
	data := l.parseContactUsData(rawData)

	//and validations go here

	err := l.contactUssvcs.Add(ctx, &data)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Contact form submitted successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (l *ContactUsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	var err = l.contactUssvcs.Delete(ctx, id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *ContactUsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// First decode into a generic map to capture all fields
	var rawData map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &rawData); err != nil {
		return helpers.InvalidJSON()
	}

	// Parse the data using the helper function
	data := l.parseContactUsData(rawData)

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.contactUssvcs.Update(ctx, id, &data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (l *ContactUsHandler) Patch(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Decode the raw data for partial updates
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &updates); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here
	id := chi.URLParam(r, "id")
	err := l.contactUssvcs.Patch(ctx, id, updates)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "Contact updated successfully",
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, response)
}

func (l *ContactUsHandler) AddMany(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var data []models.ContactUsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &data); err != nil {
		return helpers.InvalidJSON()
	}

	//and validations go here

	err := l.contactUssvcs.AddMany(ctx, data)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// V2:
func (l *ContactUsHandler) GetV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}

	result, err := l.contactUssvcs.GetV2(ctx, skip, limit, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *ContactUsHandler) GetAllV2(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	var query query.Conditions
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		return err
	}
	result, err := l.contactUssvcs.GetAllV2(ctx, &query)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, result)
}

func (l *ContactUsHandler) GetSettings(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	settings, err := l.contactUssvcs.GetSettings(ctx)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, settings)
}

func (l *ContactUsHandler) AddOrUpdateSettings(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	var settings models.ContactUsSettingsDto
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &settings); err != nil {
		return err
	}

	if err := validator.New().Struct(settings); err != nil {
		return err
	}

	err := l.contactUssvcs.AddOrUpdateSettings(ctx, &settings)
	if err != nil {
		return err
	}
	return helpers.WriteJsonCtx(ctx, w, http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}
