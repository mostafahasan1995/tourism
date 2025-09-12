package models

import "time"

// Facility represents a hotel facility with translations
type Facility struct {
	FacilityID  int64           `json:"facility_id" bson:"facility_id"`
	Facility    string          `json:"facility" bson:"facility"`
	Sort        int             `json:"sort" bson:"sort"`
	Translation []FacilityTrans `json:"translation" bson:"translation"`
	ExpiresAt   time.Time       `json:"-" bson:"expiresAt"`
}

// FacilityTrans represents a translation for a facility
type FacilityTrans struct {
	Lang     string `json:"lang" bson:"lang"`
	Facility string `json:"facility" bson:"facility"`
}

type FacilityList struct {
	Data []Facility `json:"data"`
}
