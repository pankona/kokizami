package kokizami

import (
	"fmt"
	"os"
	"path/filepath"
)

// Kokizami provides APIs to manage tasks
type Kokizami struct {
	DB DBInterface
}

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

	k.DB = dbi
	if k.DB == nil {
		k.DB = newDB(dbpath)
	}

	return k.DB.createTable()
}

// Start starts a specified Kizami to DB
func (k *Kokizami) Start(desc string) (Kizamier, error) {
	t, err := k.DB.start(desc)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Get returns a Kizami by specified id
func (k *Kokizami) Get(id int) (Kizamier, error) {
	t, err := k.DB.get(id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Edit edits a specified Kizami item
func (k *Kokizami) Edit(id int, field, newValue string) (Kizamier, error) {
	t, err := k.DB.edit(id, field, newValue)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// List returns list of Kizami
func (k *Kokizami) List() ([]Kizamier, error) {
	c, err := k.DB.count()
	if err != nil {
		return nil, err
	}

	l, err := k.DB.list(0, c)
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
	return k.DB.stop(id)
}

// StopAll updates specified task's stopped_at
func (k *Kokizami) StopAll() error {
	return k.DB.stopall()
}

// Delete delets specified task
func (k *Kokizami) Delete(id int) error {
	return k.DB.delete(id)
}
