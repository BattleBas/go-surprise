package storage

import (
	"database/sql"
	"fmt"
	"os"

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
		id SERIAL PRIMARY KEY,
		name VARCHAR (50) NOT NULL,
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
		id INTEGER PRIMARY KEY,
		giver VARCHAR (50) NOT NULL,
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
	qry := `ALTER SEQUENCE people_id_seq RESTART WITH 1`
	if _, err := db.Exec(qry); err != nil {
		return err
	}

	qry = `DELETE FROM people`
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
		err = rows.Scan(&p.ID, &p.Name, &p.Email, pq.Array(&p.Invalid))
		if err != nil {
			return matching.Group{}, err
		}
		group.People = append(group.People, p)
		fmt.Println(p.Name)
	}

	return group, nil
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

	stmt, _ := txn.Prepare(pq.CopyIn("matches", "id", "giver", "reciever"))

	for _, pair := range (*m).Pairs {
		_, err = stmt.Exec(pair.Giver.ID, pair.Giver.Name, pair.Reciever.Name)
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
