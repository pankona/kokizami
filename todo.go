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

func init() {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		panic("failed to open database")
	}
	defer db.Close()

	_ = execCreateTable(db)
}

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

// Add adds a specified ToDo item to DB
func Add(t *ToDo) error {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		panic("failed to open database")
	}
	defer db.Close()

	err = execAddToDo(db, t)
	if err != nil {
		return errors.New("failed to add ToDo")
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

// List returns list of ToDo
func List() ([]*ToDo, error) {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		panic("failed to open database")
	}
	defer db.Close()

	l, err := execList(db)
	if err != nil {
		return nil, errors.New("failed to select database")
	}

	return l, nil
}

func execList(db *sql.DB) ([]*ToDo, error) {
	q := "SELECT id, desc FROM todo"
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.New("failed to create table")
	}

	var id int
	var desc string
	for rows.Next() {
		rows.Scan(&id, &desc)
		log.Println(id, desc)
	}

	return nil, nil
}
