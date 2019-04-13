package kokizami

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

type DBMock struct {
	DBInterface
	mockCreateTable func() error
	mockStart       func(desc string) (*kizami, error)
	mockEdit        func(id int, field, newValue string) (*kizami, error)
	mockCount       func() (int, error)
	mockList        func() ([]*kizami, error)
	mockStop        func(id int) error
	mockDelete      func(id int) error
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

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
}

func TestInitializeError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but failure is expected")
	}
}

func TestStartNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := kkzm.Start("test")
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

// nolint: gocyclo
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

	kkzm := &Kokizami{}
	err = kkzm.initialize(nil, fp.Name())
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := kkzm.Start("test")
	if err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Start returned unexpected kizami instance")
	}
	if k.Desc() != "test" {
		t.Fatalf("Start returned unexpected kizami instance")
	}
	l, err := kkzm.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(l) != 1 {
		t.Fatalf("unexpected list length")
	}

	k, err = kkzm.Edit(1, "desc", "edited")
	if err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if k.ID() != 1 {
		t.Fatalf("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Fatalf("Edit returned unexpected kizami")
	}

	k, err = kkzm.Edit(1, "started_at", "2010-01-02 03:04:05")
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

	k, err = kkzm.Edit(1, "stopped_at", "2011-01-02 03:04:05")
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

	err = kkzm.Stop(1)
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

	k, err = kkzm.Get(1)
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

	err = kkzm.Delete(1)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	l, err = kkzm.List()
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("List returned unexpected result")
	}

	_, err = kkzm.Start("test")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	err = kkzm.StopAll()
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	k, err = kkzm.Edit(2, "started_at", "2010-01-02 03:04:05")
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

	k, err = kkzm.Edit(2, "stopped_at", "2010-01-02 04:04:05")
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
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}
	dbmock.mockStart = func(desc string) (*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but failure is expected")
	}
	k, err := kkzm.Start("test")
	if k != nil {
		t.Fatalf("k is not nil but nil is expected")
	}
	if err == nil {
		t.Fatalf("error is nil but non-nil is expected")
	}

	dbmock = genDefaultDBMock()
	dbmock.mockStart = func(desc string) (*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm = &Kokizami{}
	err = kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err = kkzm.Start("test")
	if k != nil {
		t.Fatalf("k is not nil but nil is expected")
	}
	if err == nil {
		t.Fatalf("error is nil but this non-nil is expected")
	}
}

func TestEditNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err := kkzm.Edit(0, "desc", "edited")
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
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}
	dbmock.mockEdit = func(id int, field, newValue string) (*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("Initialize succeeded but failure is expected")
	}
	k, err := kkzm.Edit(0, "desc", "edited")
	if k != nil {
		t.Fatalf("k is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("error is not nil but this is not expected")
	}

	dbmock = genDefaultDBMock()
	dbmock.mockEdit = func(id int, field, newValue string) (*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm = &Kokizami{}
	err = kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	k, err = kkzm.Edit(0, "desc", "edited")
	if k != nil {
		t.Fatalf("k is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("error is nil but this is not expected")
	}
}

func TestListNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("initialize failed")
	}
	ks, err := kkzm.List()
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
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}
	dbmock.mockList = func() ([]*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but this failure is expected")
	}
	ks, err := kkzm.List()
	if ks != nil {
		t.Fatalf("list of kizami is not nil but failure is expected")
	}
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	dbmock = genDefaultDBMock()
	dbmock.mockList = func() ([]*kizami, error) {
		return nil, errors.New("error")
	}

	kkzm = &Kokizami{}
	err = kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	ks, err = kkzm.List()
	if ks != nil {
		t.Fatalf("list of kizami is not nil but this is not expected")
	}
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}

func TestStopNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	err = kkzm.Stop(0)
	if err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
}

func TestStopError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}
	dbmock.mockStop = func(id int) error {
		return errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but failure is expected")
	}
	err = kkzm.Stop(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	dbmock = genDefaultDBMock()
	dbmock.mockStop = func(id int) error {
		return errors.New("error")
	}

	kkzm = &Kokizami{}
	err = kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed")
	}
	err = kkzm.Stop(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}

func TestDeleteNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}
	err = kkzm.Delete(0)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockCreateTable = func() error {
		return errors.New("error")
	}
	dbmock.mockDelete = func(id int) error {
		return errors.New("error")
	}

	kkzm := &Kokizami{}
	err := kkzm.initialize(dbmock, "")
	if err == nil {
		t.Fatalf("initialize succeeded but failure is expected")
	}
	err = kkzm.Delete(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}

	dbmock = genDefaultDBMock()
	dbmock.mockDelete = func(id int) error {
		return errors.New("error")
	}

	kkzm = &Kokizami{}
	err = kkzm.initialize(dbmock, "")
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	err = kkzm.Delete(0)
	if err == nil {
		t.Fatalf("err is nil but this is not expected")
	}
}
