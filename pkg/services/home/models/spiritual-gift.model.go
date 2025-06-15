package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiritualGiftDto struct {
	Enabled      bool               `bson:"enabled" json:"enabled"`
	PackageId    primitive.ObjectID `bson:"packageId,omitempty" json:"packageId,omitempty" validate:"required_if=Enabled true"`
	ProgramTitle string             `bson:"programTitle,omitempty" json:"programTitle,omitempty" validate:"required_if=Enabled true"`
	ProgramId    string             `bson:"programId,omitempty" json:"programId,omitempty" validate:"required_if=Enabled true"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
}

func (s *SpiritualGiftDto) Validate(v *validator.Validate) error {
	return helpers.GenericValidation(v, s)
}

type SpiritualGift struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	SpiritualGiftDto `bson:",inline"`
	Trash            bool               `bson:"trash" json:"trash"`
	CreatedBy        primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy        primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
type SpiritualGiftWithPagination struct {
	SpiritualGifts []SpiritualGift  `bson:"spiritualGifts" json:"spiritualGifts"`
	Pagination     types.Pagination `bson:"pagination" json:"pagination"`
}
