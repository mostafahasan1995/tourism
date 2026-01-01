package models

import "larsa-tourism-microservices/pkg/types"

type MemberContact struct {
	Mobile   types.PhoneNumber `bson:"mobile" json:"mobile"`
	Whatsapp types.PhoneNumber `bson:"whatsapp" json:"whatsapp"`
	Website  string            `bson:"website" json:"website"`
	Email    string            `bson:"email" json:"email"`
}

type MemberSecurity struct {
	Email       string `bson:"email" json:"email"`
	NewPassword string `bson:"newPassword" json:"newPassword"`
}

// Agent-specific contact structure
// type AgentPhone struct {
// 	Pre     string `bson:"pre" json:"pre"`
// 	Content string `bson:"content" json:"content"`
// }

type AgentContact struct {
	Phone types.PhoneNumber `bson:"phone" json:"phone"`
	Email string            `bson:"email" json:"email"`
	Web   string            `bson:"web" json:"web"`
}
