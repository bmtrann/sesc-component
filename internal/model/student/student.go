package model

import (
	model "github.com/bmtrann/sesc-component/internal/model/course"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	FirstName string
	Surname   string
	StudentId string
	Courses   []model.Course
}

type MongoStudent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AccountId string             `bson:"account_id" validate:"required"`
	FirstName string             `bson:"first_name"`
	Surname   string             `bson:"surname"`
	StudentId string             `bson:"student_id" validate:"required"`
	Courses   []model.Course     `bson:"courses"`
}
