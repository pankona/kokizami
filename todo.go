package todo

import (
	"database/sql"
	"errors"
	"log"
	"os"

	// go-sqlite3 only imported here
	_ "github.com/mattn/go-sqlite3"
)

// ToDo represents a struct of ToDo item
type ToDo struct {
	desc string
}

// NewToDo allocates a ToDo item
func NewToDo(desc string) *ToDo {
	return &ToDo{desc: desc}
}

// Add adds a specified ToDo item to DB
func Add(t *ToDo) error {
	log.Println("Add. t =", t.desc)
	// TODO: store desc to DB

	// TODO: support multi platform
	homeDir := os.Getenv("HOME")

	// TODO: error check
	_, err := os.Create(homeDir + "/.todo.db")
	if err != nil {
		return errors.New("failed to create db file")
	}

	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		return errors.New("failed to open database")
	}

	err = execCreateTable(db)
	if err != nil {
		return errors.New("failed to create table")
	}

	err = execAddToDo(db, t)
	if err != nil {
		return errors.New("failed to add ToDo")
	}

	return nil
}

// TODO: this should be done at initialization
func execCreateTable(db *sql.DB) error {
	q := "CREATE TABLE todo ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", desc VARCHAR(255) NOT NULL"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"
	_, err := db.Exec(q)
	if err != nil {
		return errors.New("failed to create table")
	}
	return nil
}

func execAddToDo(db *sql.DB, t *ToDo) error {
	q := "INSERT INTO todo "
	q += " (desc)"
	q += " VALUES"
	q += " ('" + t.desc + "')"
	_, err := db.Exec(q)
	if err != nil {
		return errors.New("failed to create table")
	}
	return nil
}
