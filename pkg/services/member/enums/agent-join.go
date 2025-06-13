package enums

type AgentJoinStatus string

const (
	AgentJoinStatusPending   AgentJoinStatus = "pending"
	AgentJoinStatusRejected  AgentJoinStatus = "rejected"
	AgentJoinStatusConverted AgentJoinStatus = "converted"
)
