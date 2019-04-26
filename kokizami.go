package kokizami

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// go-sqlite3 is imported only here
	_ "github.com/mattn/go-sqlite3"
	"github.com/pankona/kokizami/models"
	"github.com/xo/xoutil"
)

// Kokizami represents a instance of kokizami
// Kokizami provides most APIs of kokizami library
type Kokizami struct {
	DBPath string
}

// initialTime is used to insert a time value that indicates initial value of time.
var initialTime = func() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	if err != nil {
		panic(fmt.Sprintf("failed to parse time for initial value for time: %v", err))
	}
	return t
}()

type database models.XODB

func (k *Kokizami) execWithDB(f func(db database) error) error {
	conn, err := sql.Open("sqlite3", k.DBPath)
	if err != nil {
		return err
	}
	defer func() {
		e := conn.Close()
		if e != nil {
			log.Printf("failed to close DB connection: %v", e)
		}
	}()

	return f(conn)
}

// Initialize initializes Kokizami
// Kokizami's member field must be fulfilled in advance of calling this function
func (k *Kokizami) Initialize() error {
	return k.execWithDB(func(db database) error {
		if err := models.CreateKizamiTable(db); err != nil {
			return fmt.Errorf("failed to create kizami table: %v", err)
		}
		if err := models.CreateTagTable(db); err != nil {
			return fmt.Errorf("failed to create tag table: %v", err)
		}
		if err := models.CreateRelationTable(db); err != nil {
			return fmt.Errorf("failed to create relation table: %v", err)
		}
		return nil
	})
}

// Start starts a new kizami with specified desc
func (k *Kokizami) Start(desc string) (*Kizami, error) {
	var ki *Kizami
	return ki, k.execWithDB(func(db database) error {
		entry := &models.Kizami{
			Desc:      desc,
			StartedAt: sqTime(time.Now()),
			StoppedAt: sqTime(initialTime),
		}
		err := entry.Insert(db)
		if err != nil {
			return err
		}
		m, err := models.KizamiByID(db, entry.ID)
		ki = toKizami(m)
		return err
	})
}

// Get returns a Kizami by specified ID
func (k *Kokizami) Get(id int) (*Kizami, error) {
	var ki *Kizami
	return ki, k.execWithDB(func(db database) error {
		m, err := models.KizamiByID(db, id)
		ki = toKizami(m)
		return err
	})
}

// Edit edits a specified kizami and update its model
func (k *Kokizami) Edit(ki *Kizami) (*Kizami, error) {
	return ki, k.execWithDB(func(db database) error {
		m, err := models.KizamiByID(db, ki.ID)
		if err != nil {
			return err
		}

		*m = *(ki.toModel())

		err = m.Update(db)
		if err != nil {
			return err
		}
		m, err = models.KizamiByID(db, ki.ID)
		ki = toKizami(m)
		return err
	})
}

// Stop stops a on-going kizami by specified ID
func (k *Kokizami) Stop(id int) error {
	return k.execWithDB(func(db database) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		ki.StoppedAt = sqTime(time.Now())
		return ki.Update(db)
	})
}

// StopAll stops all on-going kizamis
func (k *Kokizami) StopAll() error {
	return k.execWithDB(func(db database) error {
		ks, err := models.KizamisByStoppedAt(db, sqTime(initialTime))
		if err != nil {
			return err
		}
		now := time.Now()
		for i := range ks {
			ks[i].StoppedAt = sqTime(now)
			if err := ks[i].Update(db); err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete deletes a kizami by specified ID
func (k *Kokizami) Delete(id int) error {
	return k.execWithDB(func(db database) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		return ki.Delete(db)
	})
}

// List returns all Kizamis
func (k *Kokizami) List() ([]Kizami, error) {
	var ks []Kizami
	return ks, k.execWithDB(func(db database) error {
		ms, err := models.AllKizami(db)
		if err != nil {
			return err
		}

		ks = make([]Kizami, len(ms))
		for i := range ms {
			ks[i].ID = ms[i].ID
			ks[i].Desc = ms[i].Desc
			ks[i].StartedAt = ms[i].StartedAt.Time
			ks[i].StoppedAt = ms[i].StoppedAt.Time
		}
		return nil
	})
}

// Summary returns total elapsed time of Kizamis in specified month
func (k *Kokizami) Summary(yyyymm string) ([]Elapsed, error) {
	var s []Elapsed
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}
	return s, k.execWithDB(func(db database) error {
		ms, err := models.ElapsedOfMonth(db, yyyymm)
		if err != nil {
			return err
		}

		s = make([]Elapsed, len(ms))
		for i := range ms {
			s[i].Desc = ms[i].Desc
			s[i].Count = ms[i].Count
			s[i].Elapsed = ms[i].Elapsed
		}
		return nil
	})
}

func sqTime(t time.Time) xoutil.SqTime {
	return xoutil.SqTime{Time: t}
}
