package main

import (
	"database/sql"
	"log"
)

type Repository interface {
	InitTables()
	Upsert(a Artist) (*Artist, error)
}

type SQLiteRespository struct {
	db *sql.DB
}

func NewSQLiteRespository(db *sql.DB) *SQLiteRespository {
	return &SQLiteRespository{
		db: db,
	}
}

func (r *SQLiteRespository) InitTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS artists(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		path TEXT NOT NULL UNIQUE
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRespository) Insert(a Artist) (*Artist, error) {
	res, err := r.db.Exec("INSERT INTO artists(name, url, path) values(?,?,?)", a.Name, a.Url, a.Path)
	if err != nil {
		log.Fatal("Inserting", err)
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Getting insertion id", err)
		return nil, err
	}
	a.ID = id
	return &a, nil
}
