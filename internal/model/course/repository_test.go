package model_test

import (
	"context"
	"testing"

	model "github.com/bmtrann/sesc-component/internal/model/course"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCourseRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("test list all courses", func(mt *mtest.T) {
		courseRepo := model.CreateTestRepo(mt.Coll)
		//mt.AddMockResponses(mtest.CreateSuccessResponse())

		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		course1 := mtest.CreateCursorResponse(1, "db.courses", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: id1},
			{Key: "name", Value: "SESC"},
			{Key: "description", Value: "Programming"},
			{Key: "fees", Value: 16},
		})

		course2 := mtest.CreateCursorResponse(1, "db.courses", mtest.NextBatch, bson.D{
			{Key: "_id", Value: id2},
			{Key: "name", Value: "OOP"},
			{Key: "description", Value: "Object-Oriented Programming"},
			{Key: "fees", Value: 32},
		})

		course0 := mtest.CreateCursorResponse(0, "db.courses", mtest.NextBatch)
		mt.AddMockResponses(course1, course2, course0)

		courses, err := courseRepo.GetCourses(context.TODO(), nil)

		assert.Nil(t, err)
		assert.Equal(t, len(courses), 2)
	})

	mt.Run("test find course", func(mt *mtest.T) {
		courseRepo := model.CreateTestRepo(mt.Coll)
		//mt.AddMockResponses(mtest.CreateSuccessResponse())

		id := primitive.NewObjectID()

		viewRes := model.CourseView{
			Name:        "SESC",
			Description: "Programming",
			Fees:        16,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.courses", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: id},
			{Key: "name", Value: "SESC"},
			{Key: "description", Value: "Programming"},
			{Key: "fees", Value: 16},
		}))

		model, view, err := courseRepo.FindCourse(context.TODO(), "SESC")

		assert.Nil(t, err)
		assert.Equal(t, len(model.Id), 24)
		assert.Equal(t, view, &viewRes)
	})
}
