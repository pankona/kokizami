package todo

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	// go-sqlite3 is only imported here
	_ "github.com/mattn/go-sqlite3"
)

// DBInterface represents interface of DB
type DBInterface interface {
	openDB() error
	close()
	createTable() error
	add(desc string) (*ToDo, error)
	list() ([]*ToDo, error)
	done(id int) error
}

// DB represents DB instance
type DB struct {
	DBInterface
	conn *sql.DB
}

func newDB() *DB {
	return &DB{}
}

func (db *DB) openDB() error {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	conn, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		log.Println("failed to open DB")
		return errors.New("failed to open DB")
	}
	db.conn = conn
	return nil
}

func (db *DB) close() {
	db.conn.Close()
}

func (db *DB) createTable() error {
	q := "CREATE TABLE todo ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", desc VARCHAR(255) NOT NULL"
	q += ", done INTEGER DEFAULT 0"
	q += ", created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ")"

	_, err := db.conn.Exec(q)
	if err != nil {
		return errors.New("failed to create table")
	}
	return nil
}

func (db *DB) add(desc string) (*ToDo, error) {
	q := "INSERT INTO todo "
	q += " (desc)"
	q += " VALUES"
	q += " ('" + desc + "')"

	result, err := db.conn.Exec(q)
	if err != nil {
		return nil, errors.New("failed to create table")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("failed to determine last inserted item")
	}

	t := &ToDo{}
	err = db.conn.QueryRow("SELECT id, desc, done, created_at FROM todo WHERE id=?", id).Scan(&t.id, &t.desc, &t.done, &t.createdAt)
	switch {
	case err == sql.ErrNoRows:
		log.Println("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	}
	return t, nil
}

func (db *DB) list() ([]*ToDo, error) {
	q := "SELECT id, desc, done, created_at FROM todo"

	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, errors.New("failed to select rows")
	}

	todos := make([]*ToDo, 0, 0)
	var id int
	var desc string
	var done int
	var createdAt string
	for rows.Next() {
		err = rows.Scan(&id, &desc, &done, &createdAt)
		if err != nil {
			panic(err.Error())
		}
		todos = append(todos,
			&ToDo{
				id:        id,
				desc:      desc,
				done:      done,
				createdAt: createdAt,
			})
	}
	return todos, nil
}

func (db *DB) done(id int) error {
	q := "UPDATE todo SET done = 1 WHERE id = " + strconv.Itoa(id)

	_, err := db.conn.Exec(q)
	if err != nil {
		return errors.New("failed to select rows")
	}
	return nil
}
