package storage

import (
	"database/sql"
	"os"

	"github.com/BattleBas/go-surprise/pkg/email"
	"github.com/BattleBas/go-surprise/pkg/matching"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Database exposes the supported database functionality
type Database interface {
	CreatePeopleTable() error
	CreateMatchesTable() error
	SavePeople(*matching.Group) error
	GetPeople() (matching.Group, error)
	SaveMatches(*matching.Matches) error
	GetMatched() ([]email.Match, error)
}

// DB stores the open database connection handle
type DB struct {
	*sql.DB
}

// NewDB initialize database
func NewDB() (*DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// CreatePeopleTable creates a table about people participating in Surprise
func (db *DB) CreatePeopleTable() error {

	qry := `
	CREATE TABLE IF NOT EXISTS people (
		name VARCHAR (50) NOT NULL PRIMARY KEY,
		email VARCHAR (350) NOT NULL,
		invalid TEXT[]
	)`

	if _, err := db.Exec(qry); err != nil {
		return err
	}

	return nil
}

// CreateMatchesTable creates a table for the matches given to people
func (db *DB) CreateMatchesTable() error {

	qry := `
	CREATE TABLE IF NOT EXISTS matches (
		giver VARCHAR (50) NOT NULL PRIMARY KEY,
		reciever VARCHAR (50) NOT NULL
	)`

	if _, err := db.Exec(qry); err != nil {
		return err
	}

	return nil
}

func (db *DB) deleteMatches() error {
	qry := `DELETE FROM matches`
	_, err := db.Exec(qry)
	return err
}

func (db *DB) deletePeople() error {
	qry := `DELETE FROM people`
	_, err := db.Exec(qry)
	return err
}

// SavePeople does a bulk import to save all the people into the database
func (db *DB) SavePeople(g *matching.Group) error {
	if err := db.deleteMatches(); err != nil {
		return err
	}
	if err := db.deletePeople(); err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, _ := txn.Prepare(pq.CopyIn("people", "name", "email", "invalid"))

	for _, person := range (*g).People {
		_, err = stmt.Exec(person.Name, person.Email, pq.Array(person.Invalid))
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetPeople returns the registered group of people from the database
func (db *DB) GetPeople() (matching.Group, error) {
	qry := `SELECT * FROM people`
	rows, err := db.Query(qry)
	if err != nil {
		return matching.Group{}, err
	}

	defer rows.Close()
	var group matching.Group
	for rows.Next() {
		var p matching.Person
		err = rows.Scan(&p.Name, &p.Email, pq.Array(&p.Invalid))
		if err != nil {
			return matching.Group{}, err
		}
		group.People = append(group.People, p)
	}

	return group, nil
}

// GetMatched gets the information needed to email the participants who they've been matched with
func (db *DB) GetMatched() ([]email.Match, error) {
	qry := `SELECT name, email, reciever FROM people INNER JOIN matches on people.name = matches.giver`
	rows, err := db.Query(qry)
	if err != nil {
		return []email.Match{}, err
	}

	defer rows.Close()
	matched := []email.Match{}
	for rows.Next() {
		var m email.Match
		err = rows.Scan(&m.Name, &m.Email, &m.Reciever)
		if err != nil {
			return []email.Match{}, err
		}

		matched = append(matched, m)
	}

	return matched, nil
}

// SaveMatches does a bulk import to save all the matches into the database
func (db *DB) SaveMatches(m *matching.Matches) error {
	if err := db.deleteMatches(); err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, _ := txn.Prepare(pq.CopyIn("matches", "giver", "reciever"))

	for _, pair := range (*m).Pairs {
		_, err = stmt.Exec(pair.Giver.Name, pair.Reciever.Name)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
