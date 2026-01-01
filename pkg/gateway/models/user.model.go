package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostUserData struct {
	UserId       primitive.ObjectID
	FirstName    string
	LastName     string
	Email        string
	Password     string
	Roles        []primitive.ObjectID
	Capabilities []primitive.ObjectID
}
