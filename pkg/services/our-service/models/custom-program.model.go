package models

type CustomProgram struct {
	// Title       string             `bson:"title" json:"title"`
	// ServiceType string             `bson:"serviceType" json:"serviceType"`
	// TravelReqId primitive.ObjectID `bson:"travelReqId" json:"travelReqId"`
	// CustomerId  primitive.ObjectID `bson:"customerId" json:"customerId"`
	// Status      string             `bson:"status" json:"status"`
	// Package     primitive.ObjectID `bson:"package" json:"package"`
	// ProgramType string             `bson:"programType" json:"programType"` // general - custom
	//
	Delegation          *Delegation          `bson:"delegation,omitempty" json:"delegation,omitempty"`
	BusinessMan         *BusinessMan         `bson:"businessMan,omitempty" json:"businessMan,omitempty"`
	VipCar              *VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty"`
	CustomPlan          *CustomPlan          `bson:"customPlan,omitempty" json:"customPlan,omitempty"`
	FlightTicketRequest *FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty"`
	PartnerRequest      *PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty"`
	//
	Destination *Destination `bson:"destination,omitempty" json:"destination,omitempty"`
}

// type CustomProgram struct {
// 	Id               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
// 	CustomProgramDto `bson:",inline"`
// 	Trash            bool               `bson:"trash" json:"trash"`
// 	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
// 	CreatedBy        primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
// 	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
// 	UpdatedBy        primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
// }

// type CustomProgramPagination struct {
// 	Programs   []CustomProgram  `bson:"programs" json:"programs"`
// 	Pagination types.Pagination `bson:"pagination" json:"pagination"`
// }
