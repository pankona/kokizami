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

func sqTime(t time.Time) xoutil.SqTime {
	return xoutil.SqTime{Time: t}
}

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

// EnableVerboseQuery toggles debug logging by argument
func (k *Kokizami) EnableVerboseQuery(enable bool) {
	models.XOLog = func(s string, p ...interface{}) {
		if enable {
			fmt.Printf("-------------------------------------\nQUERY: %s\n  VAL: %v\n", s, p)
		}
	}
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

		m.Desc = ki.Desc
		m.StartedAt = sqTime(ki.StartedAt)
		m.StoppedAt = sqTime(ki.StoppedAt)

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

// SummaryByTag returns total elapsed time of Kizamis in specified month grouped by tag
func (k *Kokizami) SummaryByTag(yyyymm string) ([]Elapsed, error) {
	var s []Elapsed

	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	return s, k.execWithDB(func(db database) error {
		ms, err := models.ElapsedOfMonthByTag(db, yyyymm)
		if err != nil {
			return err
		}

		s = make([]Elapsed, len(ms))
		for i := range ms {
			s[i].Tag = ms[i].Tag
			s[i].Desc = ms[i].Desc
			s[i].Count = ms[i].Count
			s[i].Elapsed = ms[i].Elapsed
		}
		return nil
	})
}

// SummaryByDesc returns total elapsed time of Kizamis in specified month grouped by desc
func (k *Kokizami) SummaryByDesc(yyyymm string) ([]Elapsed, error) {
	var s []Elapsed

	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	return s, k.execWithDB(func(db database) error {
		ms, err := models.ElapsedOfMonthByDesc(db, yyyymm)
		if err != nil {
			return err
		}

		s = make([]Elapsed, len(ms))
		for i := range ms {
			s[i].Tag = ms[i].Tag
			s[i].Desc = ms[i].Desc
			s[i].Count = ms[i].Count
			s[i].Elapsed = ms[i].Elapsed
		}
		return nil
	})
}

// AddTag adds a new tag
func (k *Kokizami) AddTag(tag string) (Tag, error) {
	var t Tag
	return t, k.execWithDB(func(db database) error {
		entry := models.Tag{Tag: tag}
		err := entry.Insert(db)
		if err != nil {
			return err
		}

		m, err := models.TagByTag(db, tag)
		if err != nil {
			return err
		}

		t = toTag(m)
		return nil
	})
}

// AddTags adds a new tags
func (k *Kokizami) AddTags(tags []string) error {
	return k.execWithDB(func(db database) error {
		ts := models.Tags(make([]models.Tag, len(tags)))
		for i := range ts {
			ts[i].Tag = tags[i]
		}
		return ts.BulkInsert(db)
	})
}

// DeleteTag deletes a specified tag
func (k *Kokizami) DeleteTag(id int) error {
	return k.execWithDB(func(db database) error {
		m, err := models.TagByID(db, id)
		if err != nil {
			return err
		}
		return m.Delete(db)
	})
}

// Tags returns list of tags
func (k *Kokizami) Tags() ([]Tag, error) {
	var ts []Tag
	return ts, k.execWithDB(func(db database) error {
		ms, err := models.AllTags(db)
		if err != nil {
			return err
		}

		ts = make([]Tag, len(ms))
		for i := range ms {
			ts[i].ID = ms[i].ID
			ts[i].Tag = ms[i].Tag
		}

		return nil
	})
}

// Tagging makes relation between specified kizami and tags
func (k *Kokizami) Tagging(kizamiID int, tagIDs []int) error {
	return k.execWithDB(func(db database) error {
		rs := models.Relations(make([]models.Relation, len(tagIDs)))
		for i := range rs {
			rs[i].KizamiID = kizamiID
			rs[i].TagID = tagIDs[i]

		}
		return rs.BulkInsert(db)
	})
}

// Untagging removes all tags from specified kizami
func (k *Kokizami) Untagging(kizamiID int) error {
	return k.execWithDB(func(db database) error {
		return models.DeleteRelationsByKizamiID(db, kizamiID)
	})
}

// TagsByKizamiID returns tags of specified kizami
func (k *Kokizami) TagsByKizamiID(kizamiID int) ([]Tag, error) {
	var ts []Tag
	return ts, k.execWithDB(func(db database) error {
		ms, err := models.TagsByKizamiID(db, kizamiID)
		if err != nil {
			return err
		}

		ts = make([]Tag, len(ms))
		for i := range ms {
			ts[i].ID = ms[i].ID
			ts[i].Tag = ms[i].Tag
		}

		return nil
	})
}

// TagsByTags returns tags by specified tags
func (k *Kokizami) TagsByTags(tags []string) ([]Tag, error) {
	var ts []Tag
	return ts, k.execWithDB(func(db database) error {
		ms, err := models.TagsByTags(db, tags)
		if err != nil {
			return err
		}

		ts = make([]Tag, len(ms))
		for i := range ms {
			ts[i] = toTag(&ms[i])
		}

		return nil
	})
}
