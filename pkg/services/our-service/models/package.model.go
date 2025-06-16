package models

import (
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageDto struct {
	Name                string              `bson:"name" json:"name"`
	Thumbnail           []types.FileField   `bson:"thumbnail" json:"thumbnail"`
	Status              enums.PackageStatus `bson:"status" json:"status"`
	IsSystemPkg         bool                `bson:"isSystemPkg" json:"isSystemPkg"`
	AllowGeneralProgram bool                `bson:"allowGeneralProgram" json:"allowGeneralProgram"`
	AllowCustomProgram  bool                `bson:"allowCustomProgram" json:"allowCustomProgram"`
}

type Package struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	PackageDto `bson:",inline"`
	Trash      bool               `bson:"trash" json:"trash"`
	CreatedBy  primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type PackageWithPagination struct {
	Packages   []Package        `bson:"packages" json:"packages"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
