package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFailedLogin(t *testing.T) {

	user := map[string]interface{}{
		"name":     "tejaa",
		"password": "0",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP request with the JSON request body
	req, err := http.NewRequest("GET", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP recorder (mock response writer)
	rr := httptest.NewRecorder()

	fmt.Println(req, rr)

	// Call the handler function with the mock request and response recorder
	Login(rr, req)
	fmt.Print("hello vishal")
	// Check if the status code is 200 Created
	if status := rr.Code; status != 200 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	//fmt.Println(len(rr.Body.String()))

	// Check if the response body is as expected
	expected := `{"invalid":true,"token":"","type":""}`
	if len(rr.Body.String()) != 38 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSuccessLogin(t *testing.T) {

	user := map[string]interface{}{
		"name":     "teja",
		"password": "0",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP request with the JSON request body
	req, err := http.NewRequest("GET", "localhost:8000/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP recorder (a.k.a. mock response writer)
	rr := httptest.NewRecorder()

	// Call the handler function with the mock request and response recorder
	Login(rr, req)

	// Check if the status code is 200 Created
	if status := rr.Code; status != 200 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	//fmt.Println(len(rr.Body.String()))

	// Check if the response body is as expected
	fmt.Println(len(rr.Body.String()))
	if len(rr.Body.String()) != 54 {
		t.Errorf("handler returned unexpected body of unexpected length: got %d want a string of length 54",
			len(rr.Body.String()))
	}
}
