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

func (k *Kokizami) execWithDB(f func(db models.XODB) error) error {
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

func (k *Kokizami) Initialize() error {
	return k.execWithDB(func(db models.XODB) error {
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

func (k *Kokizami) Start(desc string) (*models.Kizami, error) {
	var ki *models.Kizami
	return ki, k.execWithDB(func(db models.XODB) error {
		entry := &models.Kizami{
			Desc:      desc,
			StartedAt: xoutil.SqTime{time.Now()},
			StoppedAt: xoutil.SqTime{initialTime},
		}
		err := entry.Insert(db)
		if err != nil {
			return err
		}
		ki, err = models.KizamiByID(db, entry.ID)
		return err
	})
}

func (k *Kokizami) Get(id int) (*models.Kizami, error) {
	var ki *models.Kizami
	var err error
	return ki, k.execWithDB(func(db models.XODB) error {
		ki, err = models.KizamiByID(db, id)
		return err
	})
}

func (k *Kokizami) Edit(ki *models.Kizami) (*models.Kizami, error) {
	return ki, k.execWithDB(func(db models.XODB) error {
		a, err := models.KizamiByID(db, ki.ID)
		if err != nil {
			return err
		}

		*a = *ki

		err = a.Update(db)
		if err != nil {
			return err
		}
		ki, err = models.KizamiByID(db, ki.ID)
		return err
	})
}

func (k *Kokizami) Stop(id int) error {
	return k.execWithDB(func(db models.XODB) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		ki.StoppedAt = xoutil.SqTime{time.Now()}
		return ki.Update(db)
	})
}

func (k *Kokizami) StopAll() error {
	return k.execWithDB(func(db models.XODB) error {
		ks, err := models.KizamisByStoppedAt(db, xoutil.SqTime{initialTime})
		if err != nil {
			return err
		}
		now := time.Now()
		for i := range ks {
			ks[i].StoppedAt = xoutil.SqTime{now}
			if err := ks[i].Update(db); err != nil {
				return err
			}
		}
		return nil
	})
}

func (k *Kokizami) Delete(id int) error {
	return k.execWithDB(func(db models.XODB) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		return ki.Delete(db)
	})
}

func (k *Kokizami) List() ([]*models.Kizami, error) {
	var ks []*models.Kizami
	return ks, k.execWithDB(func(db models.XODB) error {
		var err error
		ks, err = models.AllKizami(db)
		return err
	})
}

func (k *Kokizami) Summary(yyyymm string) ([]*models.Elapsed, error) {
	var s []*models.Elapsed
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}
	return s, k.execWithDB(func(db models.XODB) error {
		var err error
		s, err = models.ElapsedWithQuery(db, yyyymm)
		return err
	})

	return s, nil
}
