package models

import (
	"larsa-tourism-microservices/pkg/types"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TravelPerkDto struct {
	Title       string `bson:"title" json:"title" validate:"required"`
	Description string `bson:"description" json:"description"`
	IsActive    bool   `bson:"isActive" json:"isActive"`
}

type CustomerPersonaDto struct {
	Title                  string             `bson:"title" json:"title" validate:"required"`
	ExplorationPreferences string             `bson:"explorationPreferences" json:"explorationPreferences"` // What's your ideal way to explore a new city?
	RelaxationPreferences  string             `bson:"relaxationPreferences" json:"relaxationPreferences"`   // How do you prefer to relax during a trip?
	TravelMustHaves        string             `bson:"travelMustHaves" json:"travelMustHaves"`               // What's your biggest travel must-have?
	TravelPerks            []TravelPerkDto    `bson:"travelPerks" json:"travelPerks"`                       // Special perks for this persona
	PersonaImage           types.FileField    `bson:"personaImage" json:"personaImage" validate:"required"` // Visual representation of the persona
	Description            string             `bson:"description" json:"description"`                       // Detailed description
	ProgramId              primitive.ObjectID `bson:"programId,omitempty" json:"programId,omitempty"`       // Associated program if any
	DisplayOrder           int                `bson:"displayOrder" json:"displayOrder"`                     // Order for display
	IsActive               bool               `bson:"isActive" json:"isActive"`                             // Active status
}

type CustomerPersona struct {
	CustomerPersonaDto `bson:",inline"`

	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type CustomerPersonaPagination struct {
	Personas   []CustomerPersona `bson:"personas" json:"personas"`
	Pagination common.Pagination `bson:"pagination" json:"pagination"`
}
