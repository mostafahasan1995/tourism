package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HeroSection struct {
	HotelName    string            `bson:"hotelName" json:"hotelName" validate:"required"`
	Rating       float64           `bson:"rating" json:"rating" validate:"min=0,max=5"`
	PropertyType string            `bson:"propertyType" json:"propertyType" validate:"required"`
	Overview     string            `bson:"overview" json:"overview" validate:"required"`
	Logo         *types.FileField  `bson:"logo,omitempty" json:"logo,omitempty"`
	Images       []types.FileField `bson:"images" json:"images"`
}

type FacilityOption struct {
	Id          string `bson:"id" json:"id"`
	Label       string `bson:"label" json:"label"`
	Type        string `bson:"type" json:"type"` // "checkbox" or "radio"
	Category    string `bson:"category" json:"category"`
	IsSelected  bool   `bson:"isSelected" json:"isSelected"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Icon        string `bson:"icon,omitempty" json:"icon,omitempty"`
}

type FacilitiesSection struct {
	Restaurants    []FacilityOption `bson:"restaurants" json:"restaurants"`
	PoolsBeaches   []FacilityOption `bson:"poolsBeaches" json:"poolsBeaches"`
	SpaGym         []FacilityOption `bson:"spaGym" json:"spaGym"`
	HotelServices  []FacilityOption `bson:"hotelServices" json:"hotelServices"`
	Business       []FacilityOption `bson:"business" json:"business"`
	KidsFacilities []FacilityOption `bson:"kidsFacilities" json:"kidsFacilities"`
	Recreational   []FacilityOption `bson:"recreational" json:"recreational"`
}

type DynamicSection struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title    string             `bson:"title" json:"title" validate:"required"`
	Overview string             `bson:"overview" json:"overview" validate:"required"`
	Images   []types.FileField  `bson:"images" json:"images"`
	Order    int                `bson:"order" json:"order"`
}

type ExhibitorProfileDto struct {
	HotelId           primitive.ObjectID `bson:"hotelId" json:"hotelId" validate:"required"`
	HeroSection       HeroSection        `bson:"heroSection" json:"heroSection" validate:"required"`
	FacilitiesSection FacilitiesSection  `bson:"facilitiesSection" json:"facilitiesSection"`
	DynamicSections   []DynamicSection   `bson:"dynamicSections" json:"dynamicSections"`
	IsActive          bool               `bson:"isActive" json:"isActive"`
	IsPublished       bool               `bson:"isPublished" json:"isPublished"`
}

func (e *ExhibitorProfileDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, e)
}

type ExhibitorProfile struct {
	ExhibitorProfileDto `bson:",inline"`
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash               bool               `bson:"trash" json:"trash"`
	CreatedBy           primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt           time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy           primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt           time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ExhibitorProfilePagination struct {
	Profiles   []ExhibitorProfile `bson:"profiles" json:"profiles"`
	Pagination types.Pagination   `bson:"pagination" json:"pagination"`
}

type AddDynamicSectionDto struct {
	Title    string            `bson:"title" json:"title" validate:"required"`
	Overview string            `bson:"overview" json:"overview" validate:"required"`
	Images   []types.FileField `bson:"images" json:"images"`
}

func (a *AddDynamicSectionDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type UpdateDynamicSectionDto struct {
	Title    *string            `bson:"title,omitempty" json:"title,omitempty"`
	Overview *string            `bson:"overview,omitempty" json:"overview,omitempty"`
	Images   *[]types.FileField `bson:"images,omitempty" json:"images,omitempty"`
	Order    *int               `bson:"order,omitempty" json:"order,omitempty"`
}

type FacilityUpdateDto struct {
	FacilityId string `bson:"facilityId" json:"facilityId" validate:"required"`
	Category   string `bson:"category" json:"category" validate:"required"`
	IsSelected bool   `bson:"isSelected" json:"isSelected"`
}

func (f *FacilityUpdateDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, f)
}
