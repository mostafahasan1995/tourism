package enums

type TravelReqStatus string

const (
	TravelReqStatusPending                   TravelReqStatus = "pending"
	TravelReqStatusWaitingForCustomerService TravelReqStatus = "waiting-customer-service"
	TravelReqStatusWaiting                   TravelReqStatus = "waiting"
	TravelReqStatusApproved                  TravelReqStatus = "approved"
	TravelReqStatusRejected                  TravelReqStatus = "rejected"
	TravelReqStatusCompleted                 TravelReqStatus = "completed"
)
