package database

import (
	"context"
	"go-jwt/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = ConnectDatabase()

func ConnectDatabase() *mongo.Client {
	_, uri, _ := utils.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func OpenColletion(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("go-jwt").Collection(collectionName)
}
