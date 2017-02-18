package kokizami

import (
	"strconv"
	"time"

	// go-sqlite3 is only imported here
	_ "github.com/mattn/go-sqlite3"
)

// Kizami represents a struct of task item
type Kizami struct {
	id        int
	desc      string
	startedAt time.Time
	stoppedAt time.Time
	elapsed   time.Time
}

// String returns string representation of a Kizami.
// note that the timestamps is not considered time zone.
func (k *Kizami) String() string {
	return strconv.Itoa(k.id) + "\t" +
		k.desc + "\t" +
		k.startedAt.Format("2006-01-02 15:04:05") + "\t" +
		k.stoppedAt.Format("2006-01-02 15:04:05") + "\t" +
		k.Elapsed().String()
}

// ID returns Kizami's id
func (k *Kizami) ID() int {
	return k.id
}

// Desc returns Kizami's description
func (k *Kizami) Desc() string {
	return k.desc
}

// StartedAt returns Kizami's startedAt
func (k *Kizami) StartedAt() time.Time {
	return k.startedAt
}

// StoppedAt returns Kizami's stoppedAt
func (k *Kizami) StoppedAt() time.Time {
	return k.stoppedAt
}

// Elapsed returns Kizami's elapsed time
func (k *Kizami) Elapsed() time.Duration {
	var elapsed time.Duration
	if k.stoppedAt.Unix() == 0 {
		// this Kizami is on going. show elapsed time until now.
		now := time.Now().UTC()
		elapsed = now.Sub(k.startedAt)
	} else {
		elapsed = k.stoppedAt.Sub(k.startedAt)
		if elapsed < 0 {
			elapsed = 0
		}
	}
	return elapsed
}

var dbinterface DBInterface

// Initialize initializes Kizami library.
// this function will create DB file and prepare tables.
func Initialize(dbpath string) error {
	return initialize(nil, dbpath)
}

func initialize(dbi DBInterface, dbpath string) error {
	dbinterface = dbi
	if dbinterface == nil {
		dbinterface = newDB(dbpath)
	}

	err := dbinterface.openDB()
	if err != nil {
		return err
	}
	defer dbinterface.close()

	_ = dbinterface.createTable()
	return nil
}

// Start starts a specified Kizami to DB
func Start(desc string) (*Kizami, error) {
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
func Get(id int) (*Kizami, error) {
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
func Edit(id int, field, newValue string) (*Kizami, error) {
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
func List() ([]*Kizami, error) {
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
