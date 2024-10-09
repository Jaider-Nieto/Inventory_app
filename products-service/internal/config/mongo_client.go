package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitMongoDB(url string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB!")
	Client = client
	return Client
}

func GetMongoCollection(client *mongo.Client, database, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}
