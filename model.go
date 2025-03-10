package main

type User struct {
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	City    string `json:"city" bson:"city"`
}
