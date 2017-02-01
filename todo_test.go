package todo

import "testing"

type DBMock struct {
	DBInterface
}

func (db *DBMock) openDB() error {
	return nil
}

func (db *DBMock) close() {
}

func (db *DBMock) createTable() error {
	return nil
}

func (db *DBMock) add(desc string) (*ToDo, error) {
	return &ToDo{desc: desc}, nil
}

func TestAdd(t *testing.T) {
	SetDB(&DBMock{})
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
