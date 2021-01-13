package main

import (
	"encoding/json"
	"fmt"
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
	log.Fatal(http.ListenAndServe(":8000", router))

}

func profile(w http.ResponseWriter, r *http.Request) {

	var (
		data   user
		err    error
		isAuth bool
	)

	username := r.Header.Get("Username")

	isAuth, err = auth(username)

	if err != nil {
		sendResponse(w, "FAILURE", err.Error(), data)
		return
	}

	if !isAuth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := getUserDetails(username)

	if err != nil {
		sendResponse(w, "FAILURE", err.Error(), data)
		return
	}

	sendResponse(w, "SUCCESS", "", res.Data)

}

func serviceName(w http.ResponseWriter, r *http.Request) {
	var (
		name string
		err  error
	)
	name, err = getServiceName()
	if err == nil {
		w.Write([]byte(name))
	}
	return
}

func getUserDetails(username string) (response, error) {

	var res response

	url := "http://localhost:8080/user/profile"

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	rqst, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
	}
	rqst.Header.Set("Username", username)

	resp, err := client.Do(rqst)

	if err != nil {
		return res, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	bodyString := string(body)
	fmt.Println(bodyString)

	err = json.Unmarshal([]byte(bodyString), &res)

	return res, err

}

func getServiceName() (string, error) {

	var (
		name string
		err  error
	)

	url := "http://localhost:8080/microservice/name"

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	rqst, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return name, err
	}

	resp, err := client.Do(rqst)

	if err != nil {
		return name, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return name, err
	}

	defer resp.Body.Close()

	name = string(body)
	return name, nil
}

func auth(username string) (bool, error) {

	url := "http://localhost:9000/auth"
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	rqst, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return false, err
	}

	rqst.Header.Set("Username", username)

	resp, err := client.Do(rqst)

	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}

func sendResponse(w http.ResponseWriter, status string, errMsg string, data user) {

	var res response = response{status, errMsg, data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return
}
