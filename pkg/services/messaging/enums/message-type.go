package enums

type MsgTyps int

const (
	None MsgTyps = iota
	INVITATION
	ACCOUNTUPDATED
	WELCOME_NEWSLETTER
)
