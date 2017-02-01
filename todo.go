package todo

import (
	"errors"
	"strconv"

	// go-sqlite3 is only imported here
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

// Initialize initializes ToDo library.
// this function will create DB file and prepare tables.
func Initialize() {
	db, err := openDB()
	if err != nil {
		panic("failed to open database")
	}
	defer db.close()
	_ = db.createTable()
}

// Add adds a specified ToDo item to DB
func Add(desc string) (*ToDo, error) {
	db, err := openDB()
	if err != nil {
		panic("failed to open database")
	}
	defer db.close()

	t, err := db.add(desc)
	if err != nil {
		return nil, errors.New("failed to add ToDo")
	}

	return t, nil
}

// List returns list of ToDo
func List() ([]*ToDo, error) {
	db, err := openDB()
	if err != nil {
		panic("failed to open database")
	}
	defer db.close()

	l, err := db.list()
	if err != nil {
		return nil, errors.New("failed to select database")
	}

	return l, nil
}

// Done will mark specified ToDo as "done"
func Done(id int) error {
	db, err := openDB()
	if err != nil {
		panic("failed to open database")
	}
	defer db.close()

	err = db.done(id)
	if err != nil {
		return errors.New("failed to select database")
	}

	return nil
}
