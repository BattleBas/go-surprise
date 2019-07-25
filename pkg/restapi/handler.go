package restapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/battlebas/go-surprise/pkg/matching"
	"github.com/gorilla/mux"
)

// Handler sets up all the HTTP routers
func Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/register", Register).Methods("POST")
	return r
}

// Register parses a group of people and sorts them
// making each of them a giver and assigining them a reciever
func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var givers matching.Group
	err = json.Unmarshal(body, &givers)
	if err != nil {
		panic(err)
	}

	var recievers matching.Group
	err = json.Unmarshal(body, &recievers)
	if err != nil {
		panic(err)
	}

	matches, err := matching.CreateMatches(givers, recievers)
	response, err := json.Marshal(matches)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
