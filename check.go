package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type LeaveRequest struct {
	From   string `json:"from,omitempty"`
	Days   int    `json:"days,omitempty"`
	Id     string `json:"id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status bool   `json:"status,omitempty"`
	Name   string `json:"name,omitempty"`
}
type User struct {
	Name     string `bson:"name,omitempty"`
	Password string `bson:"password,omitempty"`
	Type     string `bson:"type,omitempty"`
}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://ritvik:ritvik@cluster0.x1pdgs7.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	dbctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(dbctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(dbctx)

	err = client.Ping(dbctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	dbuser := &User{
		Name:     "saiteja",
		Password: "1234",
		Type:     "student",
	}
	coll := client.Database("lms").Collection("leave")

	_, _ = coll.InsertOne(context.TODO(), dbuser)
	//filter := bson.M{"name": "saiteja"}
	var result []LeaveRequest
	cursor, err := coll.Find(context.TODO(), bson.D{})
	//err = coll.FindOne(context.TODO(), filter).Decode(dbuser)
	if err != nil {

		log.Println(err)

	}
	cursor.All(context.Background(), &result)
	fmt.Println(result, len(result))

}
