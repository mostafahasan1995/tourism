package enums

type TravelReqStatus string

const (
	TravelReqStatusPending   TravelReqStatus = "pending"
	TravelReqStatusApproved  TravelReqStatus = "approved"
	TravelReqStatusRejected  TravelReqStatus = "rejected"
	TravelReqStatusCompleted TravelReqStatus = "completed"
)
