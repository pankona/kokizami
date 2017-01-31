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
	id        int
	desc      string
	done      int
	createdAt string
}

func (t *ToDo) Error() string {
	return strconv.Itoa(t.id) + " " + t.desc + " " + t.createdAt
}

// IsDone returns ToDo's state.
// if ToDo is marked as done, return true, otherwise false.
func (t *ToDo) IsDone() bool {
	return t.done == 1
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

// Add adds a specified ToDo item to DB
func Add(desc string) (*ToDo, error) {
	// TODO: support multi platform
	homeDir := os.Getenv("HOME")
	db, err := sql.Open("sqlite3", homeDir+"/.todo.db")
	if err != nil {
		panic("failed to open database")
	}
	defer db.Close()

	t, err := execAdd(db, desc)
	if err != nil {
		return nil, errors.New("failed to add ToDo")
	}

	return t, nil
}

func execAdd(db *sql.DB, desc string) (*ToDo, error) {
	q := "INSERT INTO todo "
	q += " (desc)"
	q += " VALUES"
	q += " ('" + desc + "')"
	result, err := db.Exec(q)
	if err != nil {
		return nil, errors.New("failed to create table")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("failed to determine last inserted item")
	}

	t := &ToDo{}
	err = db.QueryRow("SELECT id, desc, done, created_at FROM todo WHERE id=?", id).Scan(&t.id, &t.desc, &t.done, &t.createdAt)
	switch {
	case err == sql.ErrNoRows:
		log.Println("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	}
	return t, nil
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
