package models

import (
	"encoding/json"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HeroSection struct {
	HotelName    transl.Localizable[string] `bson:"hotelName" json:"hotelName" validate:"required"`
	Rating       float64                    `bson:"rating" json:"rating" validate:"min=0,max=5"`
	PropertyType transl.Localizable[string] `bson:"propertyType" json:"propertyType" validate:"required"`
	Overview     transl.Localizable[string] `bson:"overview" json:"overview" validate:"required"`
	Logo         *types.FileField           `bson:"logo,omitempty" json:"logo,omitempty"`
	Images       []types.FileField          `bson:"images" json:"images"`
}

type FacilityOption struct {
	Id          string                     `bson:"id" json:"id"`
	Label       transl.Localizable[string] `bson:"label" json:"label"`
	Type        string                     `bson:"type" json:"type"` // "checkbox" or "radio"
	Category    string                     `bson:"category" json:"category"`
	IsSelected  bool                       `bson:"isSelected" json:"isSelected"`
	Description transl.Localizable[string] `bson:"description,omitempty" json:"description,omitempty"`
	Icon        string                     `bson:"icon,omitempty" json:"icon,omitempty"`
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
	Id       primitive.ObjectID         `bson:"_id,omitempty" json:"_id,omitempty"`
	Title    transl.Localizable[string] `bson:"title" json:"title" validate:"required"`
	Overview transl.Localizable[string] `bson:"overview" json:"overview" validate:"required"`
	Images   []types.FileField          `bson:"images" json:"images"`
	Order    int                        `bson:"order" json:"order"`
}

type ExhibitorProfileDto struct {
	ExhibitionId      primitive.ObjectID     `bson:"exhibitionId" json:"exhibitionId" `
	HotelId           primitive.ObjectID     `bson:"hotelId" json:"hotelId" validate:"required"`
	HeroSection       HeroSection            `bson:"heroSection" json:"heroSection" validate:"required"`
	FacilitiesSection FacilitiesSection      `bson:"facilitiesSection" json:"facilitiesSection"`
	DynamicSections   []DynamicSection       `bson:"dynamicSections" json:"dynamicSections"`
	ContactInfo       map[string]interface{} `bson:"contactInfo" json:"contactInfo"`
	IsActive          bool                   `bson:"isActive" json:"isActive"`
	IsPublished       bool                   `bson:"isPublished" json:"isPublished"`
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
	Status              string             `bson:"status,omitempty" json:"status,omitempty"`
}

// Custom JSON marshalling for ExhibitorProfile to handle zero ObjectIDs
func (ep ExhibitorProfile) MarshalJSON() ([]byte, error) {
	type Alias ExhibitorProfile
	aux := &struct {
		*Alias
		ExhibitionId *primitive.ObjectID `json:"exhibitionId,omitempty"`
		HotelId      *primitive.ObjectID `json:"hotelId,omitempty"`
		CreatedBy    *primitive.ObjectID `json:"createdBy,omitempty"`
		UpdatedBy    *primitive.ObjectID `json:"updatedBy,omitempty"`
	}{
		Alias: (*Alias)(&ep),
	}

	// Handle ExhibitionId
	if !ep.ExhibitionId.IsZero() {
		aux.ExhibitionId = &ep.ExhibitionId
	}

	// Handle HotelId
	if !ep.HotelId.IsZero() {
		aux.HotelId = &ep.HotelId
	}

	// Handle CreatedBy
	if !ep.CreatedBy.IsZero() {
		aux.CreatedBy = &ep.CreatedBy
	}

	// Handle UpdatedBy
	if !ep.UpdatedBy.IsZero() {
		aux.UpdatedBy = &ep.UpdatedBy
	}

	return json.Marshal(aux)
}

type ExhibitorProfilePagination struct {
	Profiles   []ExhibitorProfile `bson:"profiles" json:"profiles"`
	Pagination types.Pagination   `bson:"pagination" json:"pagination"`
}

type AddDynamicSectionDto struct {
	Title    transl.Localizable[string] `bson:"title" json:"title" validate:"required"`
	Overview transl.Localizable[string] `bson:"overview" json:"overview" validate:"required"`
	Images   []types.FileField          `bson:"images" json:"images"`
}

func (a *AddDynamicSectionDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, a)
}

type UpdateDynamicSectionDto struct {
	Title    *transl.Localizable[string] `bson:"title,omitempty" json:"title,omitempty"`
	Overview *transl.Localizable[string] `bson:"overview,omitempty" json:"overview,omitempty"`
	Images   *[]types.FileField          `bson:"images,omitempty" json:"images,omitempty"`
	Order    *int                        `bson:"order,omitempty" json:"order,omitempty"`
}

type FacilityUpdateDto struct {
	FacilityId string `bson:"facilityId" json:"facilityId" validate:"required"`
	Category   string `bson:"category" json:"category" validate:"required"`
	IsSelected bool   `bson:"isSelected" json:"isSelected"`
}

func (f *FacilityUpdateDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, f)
}

type PhoneDto struct {
	Pre     string `bson:"pre" json:"pre"`
	Content string `bson:"content" json:"content"`
}

// FOUND
type ExhibitorRequestDto struct {
	ExhibitionId primitive.ObjectID `bson:"exhibitionId" json:"exhibitionId"`
	HotelId      primitive.ObjectID `bson:"hotelId" json:"hotelId" validate:"required"`

	HotelWebsite string                     `bson:"hotelWebsite" json:"hotelWebsite"`
	Phone        PhoneDto                   `bson:"phone" json:"phone"`
	Email        string                     `bson:"email" json:"email"`
	Location     transl.Localizable[string] `bson:"location" json:"location"`

	Overview transl.Localizable[string] `bson:"overview" json:"overview" validate:"required"`
	Status   string                     `bson:"status" json:"status" validate:"omitempty,oneof=Pending Replied Closed"`
}

func (e *ExhibitorRequestDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, e)
}

// ExhibitorRequestPagination represents paginated exhibitor request results
type ExhibitorRequestPagination struct {
	Profiles   []ExhibitorProfile `bson:"profiles" json:"profiles"`
	Pagination types.Pagination   `bson:"pagination" json:"pagination"`
}
