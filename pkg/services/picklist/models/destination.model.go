package models

import (
	// "larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DestinationDto struct {
	Name        string                     `bson:"name" json:"name"` //country name
	Images      []types.FileField          `bson:"images" json:"images"`
	Icon        types.FileField            `bson:"icon" json:"icon"`
	Description transl.Localizable[string] `bson:"description" json:"description"`
}

type Destination struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DestinationDto `bson:",inline"`
	// IsFav          bool               `bson:"isFav" json:"isFav"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type DestinationRes struct {
	Destination `bson:",inline"`
	IsFav       bool `bson:"isFav" json:"isFav"`
}

type DestinationPagination struct {
	Destinations []Destination    `bson:"destinations" json:"destinations"`
	Pagination   types.Pagination `bson:"pagination" json:"pagination"`
}

type DestinationPaginationRes struct {
	Destinations []DestinationRes `bson:"destinations" json:"destinations"`
	Pagination   types.Pagination `bson:"pagination" json:"pagination"`
}

type DestinationCountry struct {
	Country string `bson:"country" json:"country" validate:"required,country"`
}

// ValidateCountry validates that the country is in the allowed list
func ValidateCountry(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	for _, validCountry := range util.Countries {
		if country == validCountry {
			return true
		}
	}
	return false
}

// RegisterCustomValidations registers custom validation functions
func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("country", ValidateCountry)
}

func (d *DestinationCountry) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, d)
}
