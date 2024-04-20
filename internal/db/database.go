package db

import (
	"context"

	"github.com/bmtrann/sesc-component/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(dbConfig *config.DBConfig) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConfig.URI))

	if err != nil {
		panic(err)
	}

	return client.Database(dbConfig.Name)
}
