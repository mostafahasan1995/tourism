package types

import (
	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pagination struct {
	TotalPages float64 `json:"totalPages"`
	PerPage    int64   `json:"perPage"`
	TotalCount int64   `json:"totalCount"`
}

type User struct {
	//for now we only need the id we can add other fields as need
	Id       primitive.ObjectID
	UserData common.User
}

type FileField struct {
	Id           primitive.ObjectID `bson:"_id" json:"_id"`
	OriginalName string             `bson:"originalName" json:"originalName"`
	Path         string             `bson:"path" json:"path"`
	Service      string             `bson:"service" json:"service"`
	Expire       string             `bson:"expire" json:"expire"`
	Variants     []string           `bson:"variants" json:"variants"`
}

type ServiceToken struct {
	Token string `bson:"token" json:"token"`
}

type UserList struct {
	Users []common.User `bson:"users" json:"users"`
}

type PhoneNumber struct {
	Pre     string `bson:"pre" json:"pre"`
	Content string `bson:"content" json:"content"`
}

type CapabilityCheck struct {
	Capability string
	IsAllowed  bool
}

type Role struct {
	Id   primitive.ObjectID `json:"_id"`
	Name string             `json:"name"`
}

type RoleList struct {
	Roles []Role `json:"roles"`
}
