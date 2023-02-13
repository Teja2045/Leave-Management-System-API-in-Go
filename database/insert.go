package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type H map[string]interface{}

func ApplyLeave(w http.ResponseWriter, req *http.Request) {
	var leave LeaveRequestj
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&leave)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	log.Println(leave)
	token := leave.Token

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://saiteja:saiteja@cluster0.ugdvlxb.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	dbctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(dbctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(dbctx)

	err = client.Ping(dbctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("lms").Collection("tokens")

	filter := bson.M{"token": token}
	tokenb := &Tokenb{}
	err = coll.FindOne(context.TODO(), filter).Decode(tokenb)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(H{"alert": "invalid request. Please login before request"})
		return
	}
	leaveb := &LeaveRequestB{
		Name:   tokenb.Name,
		From:   leave.From,
		Days:   leave.Days,
		Reason: leave.Reason,
		Id:     leave.Id,
		Status: false,
	}
	coll = client.Database("lms").Collection("leave")
	_, err = coll.DeleteMany(context.TODO(), bson.M{"name": tokenb.Name})
	if err != nil {
		log.Fatal(err)
	}
	result, err := coll.InsertOne(context.TODO(), leaveb)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	json.NewEncoder(w).Encode(H{"message": "successfully applied for leave !!"})

}
