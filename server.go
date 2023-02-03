package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Userb struct {
	Name     string `bson:"name,omitempty"`
	Password string `bson:"password,omitempty"`
	Type     string `bson:"type,omitempty"`
}

type Tokenb struct {
	Name  string `bson:"name,omitempty"`
	Token string `bson:"token,omitempty"`
}

type LoginRequestj struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponsej struct {
	Invalid bool   `json:"invalid"`
	Token   string `json:"token"`
	Type    string `json:"type"`
}

type LeaveRequestj struct {
	Token  string `json:"token,omitempty"`
	From   string `json:"from,omitempty"`
	Days   int    `json:"days,omitempty"`
	Id     string `json:"id,omitempty"`
	Reason string `json:"reason,omitempty"`
}
type LeaveRequestB struct {
	From   string `json:"from,omitempty"`
	Days   int    `json:"days,omitempty"`
	Id     string `json:"id,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status bool   `json:"status,omitempty"`
	Name   string `json:"name,omitempty"`
}

type StatusRequestj struct {
	Token string `json:"token,omitempty"`
}

type StatusResponsej struct {
	Status string `json:"status,omitempty"`
}

type LogoutRequestj struct {
	Token string `json:"token,omitempty"`
}

type LogoutResponsej struct {
	Response string `json:"response,omitempty"`
}

type AcceptLeaveRequestj struct {
	Token string `json:"token,omitempty"`
	Name  string `json:"name,omitempty"`
}

type AddStudentRequestj struct {
	Token    string `json:"token,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type AddStudentResponsej struct {
	Response string `json:"response,omitempty"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	n := 10
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	//Init Router
	r := mux.NewRouter()

	// arrange our route
	r.HandleFunc("/login", Login).Methods("GET")
	r.HandleFunc("/logout", Logout).Methods("GET")
	r.HandleFunc("/applyLeave", ApplyLeave).Methods("GET")
	r.HandleFunc("/checkLeaveStatus", CheckLeaveStatus).Methods("GET")
	r.HandleFunc("/listLeaves", ListLeaves).Methods("GET")
	r.HandleFunc("/acceptLeave", AcceptLeave).Methods("GET")

	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))
}

func Login(w http.ResponseWriter, req *http.Request) {
	var user LoginRequestj
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	username := user.Name
	password := user.Password
	log.Println("name is", username)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://saiteja:saiteja@cluster0.ugdvlxb.mongodb.net/?retryWrites=true&w=majority"))
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

	coll := client.Database("lms").Collection("user")
	loginresonse := &LoginResponsej{
		Invalid: true,
		Token:   "",
		Type:    "",
	}
	filter := bson.M{"name": username}
	dbuser := &Userb{}
	err = coll.FindOne(context.TODO(), filter).Decode(dbuser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		log.Println(err)
		json.NewEncoder(w).Encode(loginresonse)
		return

	}
	if password != dbuser.Password {
		w.Header().Set("Content-Type", "application/json")
		log.Println("here at password")

		json.NewEncoder(w).Encode(loginresonse)
		return
	}

	randomToken := RandomString()
	tokentodb := &Tokenb{
		Name:  username,
		Token: randomToken,
	}

	if dbuser.Type == "student" {
		coll2 := client.Database("lms").Collection("tokens")
		_, err = coll2.InsertOne(context.TODO(), tokentodb)
		if err != nil {
			panic(err)
		}
		loginresonse = &LoginResponsej{
			Invalid: false,
			Token:   randomToken,
			Type:    "student",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginresonse)
	} else if dbuser.Type == "admin" {
		coll2 := client.Database("lms").Collection("tokena")
		_, err = coll2.InsertOne(context.TODO(), tokentodb)
		if err != nil {
			panic(err)
		}
		loginresonse = &LoginResponsej{
			Invalid: false,
			Token:   randomToken,
			Type:    "admin",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginresonse)
	}
}

func Logout(w http.ResponseWriter, req *http.Request) {
	var logout LogoutRequestj
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&logout)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	token := logout.Token

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://saiteja:saiteja@cluster0.ugdvlxb.mongodb.net/?retryWrites=true&w=majority"))
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

	coll := client.Database("lms").Collection("tokens")
	logoutresponse := &LogoutResponsej{
		Response: "Not logged in",
	}
	filter := bson.M{"token": token}
	//if token exists in db as session
	tokenb := &Tokenb{}
	db := "tokens"
	err = coll.FindOne(context.TODO(), filter).Decode(tokenb)
	if err != nil {
		coll1 := client.Database("lms").Collection("tokena")

		err = coll1.FindOne(context.TODO(), filter).Decode(tokenb)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			log.Println(err)
			json.NewEncoder(w).Encode(logoutresponse)
			return
		}
		db = "tokena"
	}

	coll = client.Database("lms").Collection(db)
	_, err = coll.DeleteMany(context.TODO(), bson.M{"name": tokenb.Name})
	if err != nil {
		log.Fatal(err)
	}

	logoutresponse = &LogoutResponsej{
		Response: "successfully logged out!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logoutresponse)
}

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

func CheckLeaveStatus(w http.ResponseWriter, req *http.Request) {
	var tokenfromrequest LogoutRequestj
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&tokenfromrequest)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	token := tokenfromrequest.Token

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

		json.NewEncoder(w).Encode(H{"alert": "invalid request. Please login before the request"})
		return
	}
	leaveb := &LeaveRequestB{}
	coll = client.Database("lms").Collection("leave")
	err = coll.FindOne(context.TODO(), bson.M{"name": tokenb.Name}).Decode(leaveb)
	if err != nil {
		json.NewEncoder(w).Encode(H{"message": "no leaves found on your name!!"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaveb)
	if leaveb.Status {
		json.NewEncoder(w).Encode(H{"message": "leave is accepted"})
	} else {
		json.NewEncoder(w).Encode(H{"message": "leave request is still pending"})
	}
}

func ListLeaves(w http.ResponseWriter, req *http.Request) {
	var tokenfromrequest LogoutRequestj
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&tokenfromrequest)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	token := tokenfromrequest.Token

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

	coll := client.Database("lms").Collection("tokena")

	filter := bson.M{"token": token}
	tokenb := &Tokenb{}
	err = coll.FindOne(context.TODO(), filter).Decode(tokenb)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(H{"alert": "invalid request. Please login before the request"})
		return
	}
	collection := client.Database("lms").Collection("leave")

	// Find all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	// Marshal the documents into a JSON array
	var results []bson.M
	for cur.Next(context.TODO()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON array to the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func AcceptLeave(w http.ResponseWriter, req *http.Request) {
	var tokenfromrequest AcceptLeaveRequestj
	log.Println(tokenfromrequest.Name)
	// use json.NewDecoder to read the JSON payload from the request body
	err := json.NewDecoder(req.Body).Decode(&tokenfromrequest)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	token := tokenfromrequest.Token

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

	coll := client.Database("lms").Collection("tokena")

	filter := bson.M{"token": token}
	tokenb := &Tokenb{}
	err = coll.FindOne(context.TODO(), filter).Decode(tokenb)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(H{"alert": "invalid request. Please login before the request"})
		return
	}
	coll = client.Database("lms").Collection("leave")

	filter = bson.M{"name": tokenfromrequest.Name}
	update := bson.M{"$set": bson.M{"status": true}}
	fmt.Println("the query is", filter)
	fmt.Println("update is ", update)
	// Find all documents in the collection
	leaveres := &LeaveRequestB{}
	result := coll.FindOneAndUpdate(context.TODO(), filter, update).Decode(leaveres)
	log.Println(leaveres)
	if err != nil {
		log.Println(err)
	}
	log.Println(result)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(H{"alert": "status is changed!"})

	//defer cur.Close(context.TODO())

}
