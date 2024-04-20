package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	Id   string
	Name string
}

type CourseView struct {
	Name        string
	Description string
	Fees        float32
}

type MongoCourse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Fees        float32            `bson:"fees"`
}
