package model

import (
	"context"
	"log"

	"github.com/bmtrann/sesc-component/internal/exception"
	model "github.com/bmtrann/sesc-component/internal/model/course"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentRepository struct {
	db *mongo.Collection
}

func NewStudentRepository(db *mongo.Database, collection string) *StudentRepository {
	studentCol := db.Collection(collection)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"accountId": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := studentCol.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}

	return &StudentRepository{
		db: studentCol,
	}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, record *MongoStudent) (*Student, error) {
	_, err := r.db.InsertOne(ctx, record)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	student := toModel(record)
	return student, nil
}

func (r *StudentRepository) GetStudent(ctx context.Context, id string) (*Student, error) {
	student := new(MongoStudent)

	err := r.db.FindOne(ctx, bson.M{
		"student_id": id,
	}).Decode(student)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return toModel(student), nil
}

func (r *StudentRepository) UpdateStudentProfile(ctx context.Context, studentId string, data map[string]string) error {
	filter := bson.M{"student_id": studentId}

	newData := bson.M{}

	if value, ok := data["firstName"]; ok {
		newData["first_name"] = value
	}

	if value, ok := data["surname"]; ok {
		newData["surname"] = value
	}

	update := bson.M{"$set": newData}

	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.MatchedCount == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func (r *StudentRepository) AddCourseToStudent(ctx context.Context, studentId string, course *model.Course) error {
	filter := bson.M{"student_id": studentId}
	update := bson.M{"$push": bson.M{"courses": course}}

	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return err
	}

	if result.MatchedCount == 0 {
		return exception.ErrUserNotFound
	}

	return nil
}

func toModel(s *MongoStudent) *Student {
	return &Student{
		FirstName: s.FirstName,
		Surname:   s.Surname,
		StudentId: s.StudentId,
		Courses:   s.Courses,
	}
}
