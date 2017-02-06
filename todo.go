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
	startedAt string
	stoppedAt string
}

func (t *ToDo) Error() string {
	return strconv.Itoa(t.id) + "\t" + t.desc + "\t" + t.startedAt + "\t" + t.stoppedAt
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

// Start starts a specified ToDo item to DB
func Start(desc string) (*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	t, err := dbinterface.start(desc)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Get returns a ToDo item by specified id
func Get(id int) (*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	t, err := dbinterface.get(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Edit edits a specified ToDo item
func Edit(id int, field, newValue string) (*ToDo, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	t, err := dbinterface.edit(id, field, newValue)
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

// Stop updates specified task's stopped_at
func Stop(id int) error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.stop(id)
	if err != nil {
		return err
	}
	return nil
}

// StopAll updates specified task's stopped_at
func StopAll() error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.stopall()
	if err != nil {
		return err
	}
	return nil
}

// Delete delets specified task
func Delete(id int) error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.delete(id)
	if err != nil {
		return err
	}
	return nil
}
