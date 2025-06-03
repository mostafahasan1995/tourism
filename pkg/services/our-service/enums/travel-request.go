package enums

type TravelReqStatus string

const (
	TravelReqStatusPending   TravelReqStatus = "pending"
	TravelReqStatusWaiting   TravelReqStatus = "waiting"
	TravelReqStatusApproved  TravelReqStatus = "approved"
	TravelReqStatusRejected  TravelReqStatus = "rejected"
	TravelReqStatusCompleted TravelReqStatus = "completed"
)
