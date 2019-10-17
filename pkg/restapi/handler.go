package restapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BattleBas/go-surprise/pkg/email"
	"github.com/BattleBas/go-surprise/pkg/matching"
	"github.com/BattleBas/go-surprise/pkg/storage"
	"github.com/gorilla/mux"
)

type env struct {
	db storage.Database
}

// Handler sets up all the HTTP routers
func Handler() http.Handler {
	db, err := storage.NewDB()
	if err != nil {
		log.Panic(err)
	}
	e := &env{db}
	err = e.db.CreatePeopleTable()
	if err != nil {
		log.Panic(err)
	}
	e.db.CreateMatchesTable()
	if err != nil {
		log.Panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/email", e.Email).Methods("POST")
	r.HandleFunc("/register", e.Register).Methods("POST")
	r.HandleFunc("/matches", e.Matches).Methods("POST")
	return r
}

// Register parses a group of people and sorts them
// making each of them a giver and assigining them a reciever
func (e *env) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var g matching.Group
	err = json.Unmarshal(body, &g)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong unmarshalling!"))
		return
	}

	err = e.db.SavePeople(&g)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong saving to database!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Succesfully registered"))
}

// Matches retreives the registered people and pairs each of them together
func (e *env) Matches(w http.ResponseWriter, r *http.Request) {

	givers, err := e.db.GetPeople()
	if err != nil {
		log.Panic(err)
	}
	recievers := givers
	recievers.People = make([]matching.Person, len(givers.People))
	copy(recievers.People, givers.People)

	matches, err := matching.CreateMatches(givers, recievers)
	err = e.db.SaveMatches(&matches)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong saving to database!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Matching successful"))
}

//Email will send an email to every participant listing the match they were given
func (e *env) Email(w http.ResponseWriter, r *http.Request) {
	m, err := e.db.GetMatched()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong getting the matches!"))
		return
	}

	for _, p := range m {
		err = email.Send(p)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something went wrong sending the emails!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - Matching successful"))
}
