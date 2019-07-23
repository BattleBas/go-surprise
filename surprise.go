package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Register parses a group of people and sorts them
// making each of them a giver and assigining them a reciever
func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	givers := Group{}
	err = json.Unmarshal(body, &givers)
	if err != nil {
		panic(err)
	}

	recievers := Group{}
	err = json.Unmarshal(body, &recievers)
	if err != nil {
		panic(err)
	}

	matches, err := CreateMatches(givers, recievers)
	response, err := json.Marshal(matches)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", Register).Methods("POST")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
