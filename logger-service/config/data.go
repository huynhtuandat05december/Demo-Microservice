package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURL = "mongodb://mongo:27017"
)

func ConnectDB() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return client, nil
}
