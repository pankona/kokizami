package todo

import (
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

var dbinterface DBInterface

// Initialize initializes ToDo library.
// this function will create DB file and prepare tables.
func Initialize(dbi DBInterface) error {
	dbinterface = dbi
	if dbinterface == nil {
		dbinterface = newDB()
	}

	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	_ = dbinterface.createTable()
	return nil
}

// Add adds a specified ToDo item to DB
func Add(desc string) (*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	t, err := dbinterface.add(desc)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Edit edits a specified ToDo item
func Edit(id int, desc string) (*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	t, err := dbinterface.edit(id, desc)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// List returns list of ToDo
func List() ([]*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	l, err := dbinterface.list()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Done will mark specified ToDo as "done"
func Done(id int) error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.done(id)
	if err != nil {
		return err
	}
	return nil
}
