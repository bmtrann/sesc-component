package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseRepository struct {
	db *mongo.Collection
}

func NewCourseRepository(db *mongo.Database, collection string) *CourseRepository {
	courseCol := db.Collection(collection)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := courseCol.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}

	return &CourseRepository{
		db: courseCol,
	}
}

func (r *CourseRepository) GetCourses(ctx context.Context, studentCourses []Course) ([]CourseView, error) {
	query := bson.M{}

	if studentCourses != nil {
		var courseNames []string
		for _, course := range studentCourses {
			courseNames = append(courseNames, course.Name)
		}
		query["name"] = bson.M{"$in": courseNames}
	}

	cursor, err := r.db.Find(ctx, query)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []MongoCourse
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var courses []CourseView
	for _, result := range results {
		res := toModel(&result)
		courses = append(courses, *res)
	}
	return courses, nil
}

func (r *CourseRepository) FindCourse(ctx context.Context, courseName string) (*Course, *CourseView, error) {
	course := new(MongoCourse)
	err := r.db.FindOne(ctx, bson.M{
		"name": courseName,
	}).Decode(course)

	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return toCourse(course), toModel(course), nil
}

func toModel(c *MongoCourse) *CourseView {
	return &CourseView{
		Name:        c.Name,
		Description: c.Description,
		Fees:        c.Fees,
	}
}

func toCourse(c *MongoCourse) *Course {
	return &Course{
		Id:   c.ID.Hex(),
		Name: c.Name,
	}
}
