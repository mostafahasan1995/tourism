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
