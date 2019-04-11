package kokizami

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

type DBMock struct {
	DBInterface
	mockOpenDB      func() error
	mockClose       func()
	mockCreateTable func() error
	mockStart       func(desc string) (*kizami, error)
	mockEdit        func(id int, field, newValue string) (*kizami, error)
	mockCount       func() (int, error)
	mockList        func() ([]*kizami, error)
	mockStop        func(id int) error
	mockDelete      func(id int) error
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

func (db *DBMock) start(desc string) (*kizami, error) {
	return db.mockStart(desc)
}

func (db *DBMock) edit(id int, field, newValue string) (*kizami, error) {
	return db.mockEdit(id, field, newValue)
}

func (db *DBMock) count() (int, error) {
	return db.mockCount()
}

func (db *DBMock) list(start, end int) ([]*kizami, error) {
	return db.mockList()
}

func (db *DBMock) stop(id int) error {
	return db.mockStop(id)
}

func (db *DBMock) delete(id int) error {
	return db.mockDelete(id)
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
		mockStart: func(desc string) (*kizami, error) {
			return &kizami{desc: "test"}, nil
		},
		mockEdit: func(id int, field, newValue string) (*kizami, error) {
			return &kizami{desc: "edited"}, nil
		},
		mockCount: func() (int, error) {
			return 3, nil
		},
		mockList: func() ([]*kizami, error) {
			t := make([]*kizami, 0)
			t = append(t, &kizami{desc: "test0"})
			t = append(t, &kizami{desc: "test1"})
			t = append(t, &kizami{desc: "test2"})
			return t, nil
		},
		mockStop: func(id int) error {
			return nil
		},
		mockDelete: func(id int) error {
			return nil
		},
	}
}

func TestInitializeNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
}

func TestInitializeError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}

	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but failure is expected")
	}
}

func TestStartNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := Start("test")
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if k == nil {
		t.Fatalf("Start returned nil")
	}
	if k.Desc() != "test" {
		t.Fatalf("Start returned unexpected value")
	}
}

func TestNormalWithDB(t *testing.T) {
	fp, err := ioutil.TempFile("", "tmp_")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		err = os.Remove(fp.Name())
		if err != nil {
			t.Fatalf("failed to remove temporary file: %v", err)
		}
	}()

	err = initialize(nil, fp.Name())
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := Start("test")
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Start returned unexpected kizami instance")
	}
	if k.Desc() != "test" {
		t.Fatalf("Start returned unexpected kizami instance")
	}
	l, err := List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(l) != 1 {
		t.Fatalf("unexpected list length")
	}

	k, err = Edit(1, "desc", "edited")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Edit returned unexpected kizami")
	}

	k, err = Edit(1, "started_at", "2010-01-02 03:04:05")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "1970-01-01 00:00:00" {
		t.Fatalf("Edit returned unexpected kizami")
	}

	k, err = Edit(1, "stopped_at", "2011-01-02 03:04:05")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "2011-01-02 03:04:05" {
		t.Fatalf("Edit returned unexpected kizami")
	}

	err = Stop(1)
	if err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Fatalf("Stop returned unexpected kizami")
	}

	k, err = Get(1)
	if err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("Stop returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Fatalf("Stop returned unexpected kizami")
	}

	err = Delete(1)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	l, err = List()
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("List returned unexpected result")
	}

	_, err = Start("test")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	err = StopAll()
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	k, err = Edit(2, "started_at", "2010-01-02 03:04:05")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if k.ID() != 2 {
		t.Fatalf("Get returned unexpected kizami")
	}
	if k.Desc() != "test" {
		t.Fatalf("Get returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("Get returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Fatalf("Get returned unexpected kizami")
	}

	k, err = Edit(2, "stopped_at", "2010-01-02 04:04:05")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if k.ID() != 2 {
		t.Fatalf("ID returned unexpected value")
	}
	if k.Desc() != "test" {
		t.Fatalf("Desc returned unexpected value")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Fatalf("StartedAt returned unexpected value")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "2010-01-02 04:04:05" {
		t.Fatalf("StoppedAt returned unexpected value")
	}
	if k.String() != "2\ttest\t2010-01-02 03:04:05\t2010-01-02 04:04:05\t1h0m0s" {
		t.Fatalf("Error returned unexpected value. actual = %v", k.String())
	}
}

func TestStartError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but failure is expected")
	}
	k, err := Start("test")
	if k != nil {
		t.Fatalf("k is not nil but nil is expected")
	}
	if err == nil {
		t.Fatalf("error is nil but non-nil is expected")
	}

	// openDB goes success but Start goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockStart = func(desc string) (*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err = Start("test")
	if k != nil {
		t.Fatalf("k is not nil but nil is expected")
	}
	if err == nil {
		t.Fatalf("error is nil but this non-nil is expected")
	}
}

func TestEditNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := Edit(0, "desc", "edited")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if k == nil {
		t.Fatalf("Edit returned nil")
	}
	if k.Desc() != "edited" {
		t.Fatalf("edit returned unexpected value")
	}
}

func TestEditError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but this is not expected: %v", err)
	}
	k, err := Edit(0, "desc", "edited")
	if k != nil {
		t.Fatalf("k is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("error is not nil but this is not expected")
	}

	// openDB goes success but Edit goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockEdit = func(id int, field, newValue string) (*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err = Edit(0, "desc", "edited")
	if k != nil {
		t.Fatalf("k is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("error is nil but this is not expected")
	}
}

func TestListNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("initialize failed")
	}
	ks, err := List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if ks == nil {
		t.Fatalf("List returned nil")
	}
	if ks[0].Desc() != "test0" {
		t.Fatalf("List returned unexpected value")
	}
	if ks[1].Desc() != "test1" {
		t.Fatalf("List returned unexpected value")
	}
	if ks[2].Desc() != "test2" {
		t.Fatalf("List returned unexpected value")
	}
}

func TestListError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but this failure is expected")
	}
	ks, err := List()
	if ks != nil {
		t.Fatalf("list of kizami is not nil but failure is expected")
	}
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	// openDB goes success but List goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockList = func() ([]*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	ks, err = List()
	if ks != nil {
		t.Fatalf("list of kizami is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}

func TestStopNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	err = Stop(0)
	if err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
}

func TestStopError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but failure is expected")
	}
	err = Stop(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	// openDB goes success but stop goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockStop = func(id int) error {
		return errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed")
	}
	err = Stop(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}

func TestDeleteNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	err = Delete(0)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but failure is expected")
	}
	err = Delete(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	// openDB goes success but stop goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockDelete = func(id int) error {
		return errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	err = Delete(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}
