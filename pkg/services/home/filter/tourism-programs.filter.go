package filter

type TourismProgramFilter struct {
	Destination string             `bson:"destination" json:"destination"`
	TravelType  string             `bson:"travelType" json:"travelType"`
	Duration    int                `bson:"duration" json:"duration"`
	GroupSize   string             `bson:"groupSize" json:"groupSize"`
	Page   int             `bson:"page" json:"page"`
	Size   int             `bson:"size" json:"size"`
}


