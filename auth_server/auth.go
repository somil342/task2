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

var session map[string]time.Time = map[string]time.Time{}

type valid struct {
	IsValid bool `json:"isvalid,omitempty"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth", auth).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func auth(w http.ResponseWriter, r *http.Request) {

	auth := false
	user := r.Header.Get("Username")

	_, ok := session[user]

	if ok {
		if time.Now().Local().Sub(session[user]).Seconds() <= 10 {
			auth = true
		}
	}

	if auth {
		w.WriteHeader(http.StatusOK)
		return
	}

	valid, err := isvalid(user)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if valid {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	return
}

func isvalid(username string) (bool, error) {

	url := "http://localhost:3000/isvalid"
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	rqst, err := http.NewRequest("GET", url, nil)
	rqst.Header.Set("Username", username)

	if err != nil {
		return false, err
	}

	resp, err := client.Do(rqst)

	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	bodyString := string(body)

	var data valid
	err = json.Unmarshal([]byte(bodyString), &data)

	if err != nil {
		return false, err
	}

	if data.IsValid {
		return true, nil
	}

	return false, nil
}
