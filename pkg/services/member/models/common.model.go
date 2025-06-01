package models

type MemberContact struct {
	Mobile   string `bson:"mobile" json:"mobile"`
	Whatsapp string `bson:"whatsapp" json:"whatsapp"`
	Website  string `bson:"website" json:"website"`
	Email    string `bson:"email" json:"email"`
}

type MemberSecurity struct {
	Email       string `bson:"email" json:"email"`
	NewPassword string `bson:"newPassword" json:"newPassword"`
}

// Agent-specific contact structure
type AgentPhone struct {
	Pre     string `bson:"pre" json:"pre"`
	Content string `bson:"content" json:"content"`
}

type AgentContact struct {
	Phone AgentPhone `bson:"phone" json:"phone"`
	Email string     `bson:"email" json:"email"`
	Web   string     `bson:"web" json:"web"`
}
