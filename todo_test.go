package todo

import "testing"

type DBMock struct {
	DBInterface
	mockOpenDB      func() error
	mockClose       func()
	mockCreateTable func() error
	mockAdd         func(desc string) (*ToDo, error)
}

func (db *DBMock) openDB() error {
	return db.mockOpenDB()
}

func (db *DBMock) close() {
	db.mockClose()
}

func (db *DBMock) createTable() error {
	return db.mockCreateTable()
}

func (db *DBMock) add(desc string) (*ToDo, error) {
	return db.mockAdd(desc)
}

func TestAdd(t *testing.T) {
	dbmock := &DBMock{
		mockOpenDB: func() error {
			return nil
		},
		mockClose: func() {
		},
		mockCreateTable: func() error {
			return nil
		},
		mockAdd: func(desc string) (*ToDo, error) {
			return &ToDo{desc: "test"}, nil
		},
	}
	SetDB(dbmock)
	Initialize()
	todo, err := Add("test")
	if err != nil {
		t.Error("Add returned error")
	}
	if todo == nil {
		t.Error("Add returned nil")
	}
	if todo.desc != "test" {
		t.Error("Add returned unexpected value")
	}
}
