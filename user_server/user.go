package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
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

type response struct {
	Status string `json:"status,omitempty"`
	ErrMsg string `json:"errmsg,omitempty"`
	Data   user   `json:"data,omitempty"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user/profile", profile).Methods("GET")
	router.HandleFunc("/microservice/name", serviceName).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func profile(w http.ResponseWriter, r *http.Request) {

	var (
		data user
		err  error
	)

	user := r.Header.Get("Username")

	data, err = getUserDetails(user)

	if err != nil {
		sendResponse(w, "FAILURE", err.Error(), data)
		return
	}

	sendResponse(w, "SUCCESS", "", data)
}

func serviceName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user-microservice"))
	return
}

func getUserDetails(username string) (user, error) {

	var data user

	url := "http://localhost:3000/getdata"
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	rqst, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	rqst.Header.Set("Username", username)

	resp, err := client.Do(rqst)

	if err != nil {
		return data, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	bodyString := string(body)

	err = json.Unmarshal([]byte(bodyString), &data)

	return data, err
}

func sendResponse(w http.ResponseWriter, status string, errMsg string, data user) {

	var res response = response{status, errMsg, data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return
}
