package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoClient *mongo.Client

	UsersCollection = "users"
)

func InitializeMongoConnection() {
	MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}

	err = MongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB!")
}

func InsertDocument(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(collectionName).InsertOne(context.TODO(), document)

	return result, err
}
