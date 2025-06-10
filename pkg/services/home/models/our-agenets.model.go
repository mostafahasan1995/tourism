package models

// deprecated
// type OurAgentsDto struct {
// 	FullName          string          `bson:"fullName" json:"fullName"`
// 	Email             string          `bson:"email" json:"email"`
// 	PhoneNumber       string          `bson:"phoneNumber" json:"phoneNumber"`
// 	Bio               string          `bson:"bio" json:"bio"`
// 	Nationality       string          `bson:"nationality" json:"nationality"`
// 	LanguagesSpoken   []string        `bson:"languagesSpoken" json:"languagesSpoken"`
// 	CountriesYouServe []string        `bson:"countriesYouServe" json:"countriesYouServe"`
// 	CompanyName       string          `bson:"companyName" json:"companyName"`
// 	CompanyLogo       types.FileField `bson:"companyLogo" json:"companyLogo"`
// }

// type OurAgents struct {
// 	OurAgentsDto `bson:",inline"`

// 	Id primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

// 	Trash bool `bson:"trash" json:"trash"`

// 	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
// 	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
// 	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
// 	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
// }

// type OurAgentsPagination struct {
// 	OurAgents []OurAgents `bson:"ourAgents" json:"ourAgents"`

// 	Pagination common.Pagination `bson:"pagination" json:"pagination"`
// }
