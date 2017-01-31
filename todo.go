package todo

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	// go-sqlite3 only imported here
	_ "github.com/mattn/go-sqlite3"
)

// ToDo represents a struct of ToDo item
type ToDo struct {
	desc string
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
	q += ", done INTEGER DEFAULT 0"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"
	_, err := db.Exec(q)
	if err != nil {
		return errors.New("failed to create table")
	}
	return nil
}

// NewToDo allocates a ToDo item
func NewToDo(desc string) *ToDo {
	return &ToDo{desc: desc}
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
	q := "SELECT id, desc, done, created_at FROM todo"
	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.New("failed to select rows")
	}

	var id int
	var desc string
	var done int
	var createdAt string
	for rows.Next() {
		err = rows.Scan(&id, &desc, &done, &createdAt)
		if err != nil {
			panic(err.Error())
		}
		log.Println(id, desc, done, createdAt)
	}
	return nil, nil
}

// Done will mark specified ToDo item as "done"
func Done(id int) error {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		panic("failed to open database")
	}
	defer db.Close()

	err = execDone(db, id)
	if err != nil {
		return errors.New("failed to select database")
	}

	return nil
}

func execDone(db *sql.DB, id int) error {
	q := "UPDATE todo SET done = 1 WHERE id = " + strconv.Itoa(id)
	_, err := db.Exec(q)
	if err != nil {
		return errors.New("failed to select rows")
	}
	return nil
}
