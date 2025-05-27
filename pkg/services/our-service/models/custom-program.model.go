package models

type CustomProgram struct {
	Delegation          Delegation          `bson:"delegation,omitempty" json:"delegation,omitempty"`
	BusinessMan         BusinessMan         `bson:"businessMan,omitempty" json:"businessMan,omitempty"`
	VipCar              VipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty"`
	CustomPlan          CustomPlan          `bson:"customPlan,omitempty" json:"customPlan,omitempty"`
	FlightTicketRequest FlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty"`
	PartnerRequest      PartnerRequest      `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty"`
	//
	Destinations []Destination `bson:"destinations,omitempty" json:"destinations,omitempty"`
}
