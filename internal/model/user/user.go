package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       string
	Username string
	Password string
}

type MongoUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username" validate:"required"`
	Password string             `bson:"password" validate:"required"`
}
