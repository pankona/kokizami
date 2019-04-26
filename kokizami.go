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

type Kizami struct {
	ID        int
	Desc      string
	StartedAt time.Time
	StoppedAt time.Time
}

func (k *Kizami) toModel() *models.Kizami {
	return &models.Kizami{
		ID:        k.ID,
		Desc:      k.Desc,
		StartedAt: SqTime(k.StartedAt),
		StoppedAt: SqTime(k.StoppedAt),
	}
}

func toKizami(m *models.Kizami) *Kizami {
	return &Kizami{
		ID:        m.ID,
		Desc:      m.Desc,
		StartedAt: m.StartedAt.Time,
		StoppedAt: m.StoppedAt.Time,
	}
}

// Elapsed returns kizami's elapsed time
func (k *Kizami) Elapsed() time.Duration {
	var elapsed time.Duration
	if k.StoppedAt.Unix() == 0 {
		// this Kizami is on going. Show elapsed time until now.
		now := time.Now().UTC()
		elapsed = now.Sub(k.StartedAt)
	} else {
		elapsed = k.StoppedAt.Sub(k.StartedAt)
		if elapsed < 0 {
			elapsed = 0
		}
	}
	return elapsed
}

type Elapsed struct {
	Desc    string
	Count   int
	Elapsed time.Duration
}

func (e *Elapsed) toModel() *models.Elapsed {
	return (*models.Elapsed)(e)
}

func toElapsed(m *models.Elapsed) *Elapsed {
	return (*Elapsed)(m)
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

func (k *Kokizami) Start(desc string) (*Kizami, error) {
	var ki *Kizami
	return ki, k.execWithDB(func(db models.XODB) error {
		entry := &models.Kizami{
			Desc:      desc,
			StartedAt: SqTime(time.Now()),
			StoppedAt: SqTime(initialTime),
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

func (k *Kokizami) Get(id int) (*Kizami, error) {
	var ki *Kizami
	return ki, k.execWithDB(func(db models.XODB) error {
		m, err := models.KizamiByID(db, id)
		ki = toKizami(m)
		return err
	})
}

func (k *Kokizami) Edit(ki *Kizami) (*Kizami, error) {
	return ki, k.execWithDB(func(db models.XODB) error {
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

func (k *Kokizami) Stop(id int) error {
	return k.execWithDB(func(db models.XODB) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		ki.StoppedAt = SqTime(time.Now())
		return ki.Update(db)
	})
}

func (k *Kokizami) StopAll() error {
	return k.execWithDB(func(db models.XODB) error {
		ks, err := models.KizamisByStoppedAt(db, SqTime(initialTime))
		if err != nil {
			return err
		}
		now := time.Now()
		for i := range ks {
			ks[i].StoppedAt = SqTime(now)
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

func (k *Kokizami) List() ([]Kizami, error) {
	var ks []Kizami
	return ks, k.execWithDB(func(db models.XODB) error {
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
}

func SqTime(t time.Time) xoutil.SqTime {
	return xoutil.SqTime{Time: t}
}
