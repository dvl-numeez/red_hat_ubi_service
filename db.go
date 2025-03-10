package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO_URI string
var DB_NAME string
var COLLECTION_NAME string

type RedHatDBService struct {
	collection *mongo.Collection
}

func NewRedHatDBService(collectionName string) *RedHatDBService {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Unable to connect to DB err: ", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Unable to connect to DB err: ", err)
	}
	collection := client.Database(DB_NAME).Collection(collectionName)
	return &RedHatDBService{
		collection: collection,
	}

}

func (rs *RedHatDBService) fetchUsers() ([]User, error) {
	var result []User
	cur, err := rs.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var element User
		err := cur.Decode(&element)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, element)
	}

	return result, nil
}

func (rs *RedHatDBService) addUser(user User) error {
	_, err := rs.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
