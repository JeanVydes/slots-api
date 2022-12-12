package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByID(accountID string) (User, bool) {
	var user User

	filter := bson.M{"_id": accountID}
	_ = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func GetUserByEmailAddress(emailAddress string) (User, bool) {
	var user User

	filter := bson.M{"email": emailAddress}
	err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return user, false
	}

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func GetUserByUsername(username string) (User, bool) {
	var user User

	filter := bson.M{"username": username}
	err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return user, false
	}

	if user.ID == "" {
		return user, false
	}

	return user, true
}

func UpdateUserBalance(accountID string, newBalance Balances) error {
	filter := bson.M{"_id": accountID}
	update := bson.M{"$set": bson.M{"balances": newBalance}}

	_, err = MongoClient.Database(os.Getenv("MAIN_DATABASE_NAME")).Collection(UsersCollection).UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}
