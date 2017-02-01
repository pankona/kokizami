package todo

import "testing"

type DBMock struct {
	DBInterface
	mockOpenDB      func() error
	mockClose       func()
	mockCreateTable func() error
	mockAdd         func(desc string) (*ToDo, error)
	mockList        func() ([]*ToDo, error)
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

func (db *DBMock) list() ([]*ToDo, error) {
	return db.mockList()
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

func TestList(t *testing.T) {
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
			return nil, nil
		},
		mockList: func() ([]*ToDo, error) {
			t := make([]*ToDo, 0, 0)
			t = append(t, &ToDo{desc: "test0"})
			t = append(t, &ToDo{desc: "test1"})
			t = append(t, &ToDo{desc: "test2"})
			return t, nil
		},
	}
	SetDB(dbmock)
	Initialize()
	todos, err := List()
	if err != nil {
		t.Error("List returned error")
	}
	if todos == nil {
		t.Error("List returned nil")
	}
	if todos[0].desc != "test0" {
		t.Error("List returned unexpected value")
	}
	if todos[1].desc != "test1" {
		t.Error("List returned unexpected value")
	}
	if todos[2].desc != "test2" {
		t.Error("List returned unexpected value")
	}
}
