package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type user struct {
	Username    string    `json:"username,omitempty"`
	Dob         time.Time `json:"dob,omitempty"`
	Age         int       `json:"age,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
}

type valid struct {
	IsValid bool `json:"isvalid,omitempty"`
}

var user_list = map[string]user{
	"abhay": user{"abhay", getAge("1996-11-18"), 25, "abhay@yahoo.com", "1111111111"},
	"ajay":  user{"ajay", getAge("1997-10-18"), 24, "ajay@yahoo.com", "222222222"},
	"james": user{"james", getAge("1996-01-01"), 25, "james@yahoo.com", "3333333"},
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/isvalid", isValid).Methods("GET")
	router.HandleFunc("/getdata", getData).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func isValid(w http.ResponseWriter, r *http.Request) {

	userName := r.Header.Get("Username")
	_, ok := getDetails(userName)

	var Data valid = valid{ok}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Data)

	return
}

func getData(w http.ResponseWriter, r *http.Request) {

	userName := r.Header.Get("Username")

	user_data, _ := getDetails(userName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user_data)
	return
}

func getDetails(userName string) (user, bool) {
	data, ok := user_list[userName]
	return data, ok
}

func getAge(str string) time.Time {
	t, _ := time.Parse("2006-01-02", str)
	return t
}
