package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/newtoallofthis123/ranhash"

	_ "github.com/lib/pq"
)

func NewDBInstance() *sql.DB {
	db, err := sql.Open("postgres", GetDBUrl())
	if err != nil {
		log.Println("error is : ", err)
		panic("Unable to connect the db")
	}

	query := `
	CREATE TABLE IF NOT EXISTS subs(
		id TEXT PRIMARY KEY,
		name TEXT,
		password TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS notifications(
		id TEXT PRIMARY KEY,
		title TEXT,
		content TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Println("error is :", err)
		panic("Unable to Init the DB")
	}

	return db

}

type Subscriber struct {
	id         string
	name       string
	password   string
	created_at string
	valid      bool
}

func (s *Server) InsertSub(name string) (Subscriber, error) {
	query := `
	INSERT INTO subs(id, name, password)
	VALUES($1, $2, $3);
	`

	sub := Subscriber{
		id:         ranhash.RanHash(10),
		name:       name,
		password:   ranhash.RanHash(16),
		created_at: time.Now().String(),
		valid:      true,
	}
	_, err := s.db.Exec(query, sub.id, sub.name, sub.password)
	if err != nil {
		return Subscriber{}, err
	}

	return sub, nil
}

func (s *Server) GetSub(id string) (Subscriber, error) {
	query := `
	SELECT * from subs where id=$1;
	`

	var sub Subscriber
	rows := s.db.QueryRow(query, id)
	err := rows.Scan(&sub.id, &sub.name, &sub.password, &sub.created_at)
	if err != nil {
		return Subscriber{}, err
	}

	sub.valid = true
	return sub, nil
}
