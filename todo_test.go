package todo

import (
	"errors"
	"testing"
)

type DBMock struct {
	DBInterface
	mockOpenDB      func() error
	mockClose       func()
	mockCreateTable func() error
	mockAdd         func(desc string) (*ToDo, error)
	mockEdit        func(id int, desc string) (*ToDo, error)
	mockList        func() ([]*ToDo, error)
	mockDone        func(id int) error
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

func (db *DBMock) edit(id int, desc string) (*ToDo, error) {
	return db.mockEdit(id, desc)
}

func (db *DBMock) list() ([]*ToDo, error) {
	return db.mockList()
}

func (db *DBMock) done(id int) error {
	return db.mockDone(id)
}

// default mock implementation
func genDefaultDBMock() *DBMock {
	return &DBMock{
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
		mockEdit: func(id int, desc string) (*ToDo, error) {
			return &ToDo{desc: "edited"}, nil
		},
		mockList: func() ([]*ToDo, error) {
			t := make([]*ToDo, 0, 0)
			t = append(t, &ToDo{desc: "test0"})
			t = append(t, &ToDo{desc: "test1"})
			t = append(t, &ToDo{desc: "test2"})
			return t, nil
		},
		mockDone: func(id int) error {
			return nil
		},
	}
}

func TestInitializeNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
}

func TestInitializeError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}

	err := Initialize(dbmock)
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
}

func TestAddNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
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

func TestAddError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := Initialize(dbmock)
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
	todo, err := Add("test")
	if todo != nil {
		t.Error("todo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is not nil but this is not expected")
	}

	// openDB goes success but Add goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockAdd = func(desc string) (*ToDo, error) {
		return nil, errors.New("error")
	}
	err = Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
	todo, err = Add("test")
	if todo != nil {
		t.Error("todo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is nil but this is not expected")
	}
}

func TestEditNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
	todo, err := Edit(0, "edited")
	if err != nil {
		t.Error("Add returned error")
	}
	if todo == nil {
		t.Error("Add returned nil")
	}
	if todo.desc != "edited" {
		t.Error("edit returned unexpected value")
	}
}

func TestEditError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := Initialize(dbmock)
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
	todo, err := Edit(0, "edited")
	if todo != nil {
		t.Error("todo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is not nil but this is not expected")
	}

	// openDB goes success but Add goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockEdit = func(id int, desc string) (*ToDo, error) {
		return nil, errors.New("error")
	}
	err = Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
	todo, err = Edit(0, "edited")
	if todo != nil {
		t.Error("todo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is nil but this is not expected")
	}
}

func TestListNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	Initialize(dbmock)
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

func TestListError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	Initialize(dbmock)
	todos, err := List()
	if todos != nil {
		t.Error("list of ToDo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("err is nil but this is not expected")
	}

	// openDB goes success but List goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockList = func() ([]*ToDo, error) {
		return nil, errors.New("error")
	}
	err = Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
	todos, err = List()
	if todos != nil {
		t.Error("list of ToDo is not nil but this is not expected")
	}
	if err == nil {
		t.Error("err is nil but this is not expected")
	}
}

func TestDoneNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	Initialize(dbmock)
	err := Done(0)
	if err != nil {
		t.Error("Done returned error")
	}
}

func TestDoneError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	Initialize(dbmock)
	err := Done(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}

	// openDB goes success but Done goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockDone = func(id int) error {
		return errors.New("error")
	}
	err = Initialize(dbmock)
	if err != nil {
		t.Error("Initialize failed")
	}
	err = Done(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}
}
