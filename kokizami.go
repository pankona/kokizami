package kokizami

import (
	"fmt"
	"os"
	"path/filepath"
)

var dbinterface DBInterface

// Kokizami provides APIs to manage tasks
type Kokizami struct{}

// Initialize initializes Kizami library.
// this function will create DB file and prepare tables.
func (k *Kokizami) Initialize(dbpath string) error {
	return k.initialize(nil, dbpath)
}

func (k *Kokizami) initialize(dbi DBInterface, dbpath string) error {
	err := os.MkdirAll(filepath.Dir(dbpath), 0755) // #nosec
	if err != nil {
		return fmt.Errorf("failed to create a directory to store DB: %v", err)
	}

	dbinterface = dbi
	if dbinterface == nil {
		dbinterface = newDB(dbpath)
	}

	err = dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	return dbinterface.createTable()
}

// Start starts a specified Kizami to DB
func (k *Kokizami) Start(desc string) (Kizamier, error) {
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

// Get returns a Kizami by specified id
func (k *Kokizami) Get(id int) (Kizamier, error) {
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

// Edit edits a specified Kizami item
func (k *Kokizami) Edit(id int, field, newValue string) (Kizamier, error) {
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

// List returns list of Kizami
func (k *Kokizami) List() ([]Kizamier, error) {
	err := dbinterface.openDB()
	if err != nil {
		return nil, err
	}
	defer dbinterface.close()

	c, err := dbinterface.count()
	if err != nil {
		return nil, err
	}

	l, err := dbinterface.list(0, c)
	if err != nil {
		return nil, err
	}
	kizamiers := make([]Kizamier, 0)
	for _, v := range l {
		kizamiers = append(kizamiers, v)
	}
	return kizamiers, nil
}

// Stop updates specified task's stopped_at
func (k *Kokizami) Stop(id int) error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.stop(id)
	return err
}

// StopAll updates specified task's stopped_at
func (k *Kokizami) StopAll() error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.stopall()
	return err
}

// Delete delets specified task
func (k *Kokizami) Delete(id int) error {
	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	err = dbinterface.delete(id)
	return err
}
