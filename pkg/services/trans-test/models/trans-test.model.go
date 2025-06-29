package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Localizable[T any] map[string]T

// func (l Localizable[T]) MarshalJSON(ctx context.Context) ([]byte, error) {
// 	fmt.Println("calling custom marshaler")
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lang := cfg.Lang

// 	content, exists := l[lang]
// 	if exists {
// 		return json.Marshal(content)
// 	} else {
// 		return json.Marshal(l["en"])
// 	}
// }

// TransTestDto with multi-language support
type TransTestDto struct {
	Name        transl.Localizable[string]   `bson:"name" json:"name"`
	Description transl.Localizable[string]   `bson:"description" json:"description"`
	About       About                        `bson:"about" json:"about"`
	Notes       []transl.Localizable[string] `bson:"notes" json:"notes"`
	Date        time.Time                    `bson:"date" json:"date"`
	Images      []types.FileField            `bson:"images" json:"images"`
}

type About struct {
	Bio          transl.Localizable[string] `bson:"bio" json:"bio"`
	Testimonials []Testimonial              `bson:"testimonials" json:"testimonials"`
	Project      Project                    `bson:"project" json:"project"`
}

type Testimonial struct {
	Name        transl.Localizable[string] `bson:"name" json:"name"`
	Testimonial transl.Localizable[string] `bson:"testimonial" json:"testimonial"`
}

type Project struct {
	ProjectName        transl.Localizable[string] `bson:"project_name" json:"project_name"`
	ProjectDescription transl.Localizable[string] `bson:"project_description" json:"project_description"`
	Details            ProjectDetails             `bson:"details" json:"details"`
}

type ProjectDetails struct {
	Description transl.Localizable[string] `bson:"description" json:"description"`
	Location    transl.Localizable[string] `bson:"location" json:"location"`
}

type TransTest struct {
	Id           primitive.ObjectID `bson:"_id" json:"_id"`
	TransTestDto `bson:",inline"`
	Status       string             `bson:"status" json:"status"`
	Trash        bool               `bson:"trash" json:"trash"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy    primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}

type TransTestPagination struct {
	Trans      []TransTest      `bson:"trans" json:"trans"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}
