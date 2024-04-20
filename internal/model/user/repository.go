package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepository {
	userCol := db.Collection(collection)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCol.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}

	return &UserRepository{
		db: userCol,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, record *MongoUser) error {
	_, err := r.db.InsertOne(ctx, record)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, username, password string) (*User, error) {
	user := new(MongoUser)
	err := r.db.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	}).Decode(user)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return toModel(user), nil
}

func toModel(u *MongoUser) *User {
	return &User{
		Id:       u.ID.Hex(),
		Username: u.Username,
		Password: u.Password,
	}
}
