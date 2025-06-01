package models

// import (
// 	// "larsa-tourism-microservices/pkg/types"
// 	"larsa-tourism-microservices/pkg/types"
// 	"time"

// 	"git.larsa.io/mahdawi/microservices-commons.git/common"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// // deprecated
// type TourismProgramDto struct {
// 	Name        string   `bson:"name" json:"name"`
// 	Destination string   `bson:"destination" json:"destination"`
// 	Description string   `bson:"description" json:"description"`
// 	TravelType  string   `bson:"travelType" json:"travelType"`
// 	Duration    string   `bson:"duration" json:"duration"`
// 	Price       int      `bson:"price" json:"price"`
// 	Interests   []string `bson:"interests" json:"interests"`

// 	GroupSize   string            `bson:"groupSize" json:"groupSize"`
// 	Image       types.FileField   `bson:"image" json:"image"`
// 	Gallery     []types.FileField `bson:"gallery" json:"gallery"`
// 	Activities  []Activity        `bson:"activities" json:"activities"`
// 	DaysProgram []DayProgram      `bson:"daysProgram" json:"daysProgram"`
// }

// type Activity struct {
// 	Name        string          `bson:"name" json:"name"`
// 	Description string          `bson:"description" json:"description"`
// 	Image       types.FileField `bson:"image" json:"image"`
// }

// type DayProgram struct {
// 	Title      string            `bson:"title" json:"title"`
// 	Advantages []string          `bson:"advantages" json:"advantages"`
// 	Gallery    []types.FileField `bson:"gallery" json:"gallery"`
// }

// type TourismProgram struct {
// 	TourismProgramDto `bson:",inline"`

// 	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

// 	Trash bool `bson:"trash" json:"trash"`

// 	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
// 	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
// 	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
// 	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
// }

// type TourismProgramRes struct {
// 	TourismProgram `bson:",inline"`
// 	IsFav          bool `bson:"isFav" json:"isFav"`
// }

// type TourismProgramPagination struct {
// 	TourismProgram []TourismProgram `bson:"tourismProgram" json:"tourismProgram"`

// 	Pagination common.Pagination `bson:"pagination" json:"pagination"`
// }
