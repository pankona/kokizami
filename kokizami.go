package kokizami

import (
	"database/sql"
	"fmt"
	"time"

	// go-sqlite3 is imported only here
	_ "github.com/mattn/go-sqlite3"
	"github.com/pankona/kokizami/models"
)

// Kokizami represents a instance of kokizami
// Kokizami provides most APIs of kokizami library
type Kokizami struct {
	DBPath string
	db     *sql.DB
	now    func() time.Time

	TagRepo TagRepository
}

// initialTime is used to insert a time value that indicates initial value of time.
func initialTime() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	if err != nil {
		panic(fmt.Sprintf("failed to parse time for initial value for time: %v", err))
	}
	return t.UTC()
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
	if k.now == nil {
		k.now = time.Now
	}

	if k.db == nil {
		db, err := sql.Open("sqlite3", k.DBPath)
		if err != nil {
			return err
		}
		k.db = db
	}

	if err := models.CreateKizamiTable(k.db); err != nil {
		return fmt.Errorf("failed to create kizami table: %v", err)
	}
	if err := models.CreateTagTable(k.db); err != nil {
		return fmt.Errorf("failed to create tag table: %v", err)
	}
	if err := models.CreateRelationTable(k.db); err != nil {
		return fmt.Errorf("failed to create relation table: %v", err)
	}

	return nil
}

// DB returns db. this function is temporary use for refactoring.
func (k *Kokizami) DB() *sql.DB {
	return k.db
}

// Finalize finalizes kokizami
func (k *Kokizami) Finalize() error {
	if k.db == nil {
		return nil
	}

	return k.db.Close()
}

// Start starts a new kizami with specified desc
func (k *Kokizami) Start(desc string) (*Kizami, error) {
	if len(desc) == 0 {
		return nil, fmt.Errorf("desc must not be empty")
	}

	entry := &models.Kizami{
		Desc:      desc,
		StartedAt: models.SqTime(k.now().UTC()),
		StoppedAt: models.SqTime(initialTime()),
	}
	err := entry.Insert(k.db)
	if err != nil {
		return nil, err
	}

	m, err := models.KizamiByID(k.db, entry.ID)
	return toKizami(m), err
}

// Get returns a Kizami by specified ID
func (k *Kokizami) Get(id int) (*Kizami, error) {
	m, err := models.KizamiByID(k.db, id)
	if err != nil {
		return nil, err
	}
	return toKizami(m), nil
}

// Edit edits a specified kizami and update its model
func (k *Kokizami) Edit(ki *Kizami) (*Kizami, error) {
	m, err := models.KizamiByID(k.db, ki.ID)
	if err != nil {
		return nil, err
	}

	m.Desc = ki.Desc
	m.StartedAt = models.SqTime(ki.StartedAt.UTC())
	m.StoppedAt = models.SqTime(ki.StoppedAt.UTC())

	err = m.Update(k.db)
	if err != nil {
		return nil, err
	}
	m, err = models.KizamiByID(k.db, ki.ID)
	return toKizami(m), err
}

// Stop stops a on-going kizami by specified ID
func (k *Kokizami) Stop(id int) error {
	ki, err := models.KizamiByID(k.db, id)
	if err != nil {
		return err
	}
	ki.StoppedAt = models.SqTime(k.now().UTC())
	return ki.Update(k.db)
}

// StopAll stops all on-going kizamis
func (k *Kokizami) StopAll() error {
	ks, err := models.KizamisByStoppedAt(k.db, models.SqTime(initialTime()))
	if err != nil {
		return err
	}
	now := k.now().UTC()
	for i := range ks {
		ks[i].StoppedAt = models.SqTime(now)
		if err := ks[i].Update(k.db); err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes a kizami by specified ID
func (k *Kokizami) Delete(id int) error {
	ki, err := models.KizamiByID(k.db, id)
	if err != nil {
		return err
	}
	return ki.Delete(k.db)
}

// List returns all Kizamis
func (k *Kokizami) List() ([]Kizami, error) {
	ms, err := models.AllKizami(k.db)
	if err != nil {
		return nil, err
	}

	ks := make([]Kizami, len(ms))
	for i := range ms {
		ks[i].ID = ms[i].ID
		ks[i].Desc = ms[i].Desc
		ks[i].StartedAt = ms[i].StartedAt.Time
		ks[i].StoppedAt = ms[i].StoppedAt.Time
	}

	return ks, nil
}

// SummaryByTag returns total elapsed time of Kizamis in specified month grouped by tag
func (k *Kokizami) SummaryByTag(yyyymm string) ([]Elapsed, error) {
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	ms, err := models.ElapsedOfMonthByTag(k.db, yyyymm)
	if err != nil {
		return nil, err
	}

	s := make([]Elapsed, len(ms))
	for i := range ms {
		s[i].Tag = ms[i].Tag
		s[i].Desc = ms[i].Desc
		s[i].Count = ms[i].Count
		s[i].Elapsed = ms[i].Elapsed
	}

	return s, nil
}

// SummaryByDesc returns total elapsed time of Kizamis in specified month grouped by desc
func (k *Kokizami) SummaryByDesc(yyyymm string) ([]Elapsed, error) {
	// validate input
	_, err := time.Parse("2006-01", yyyymm)
	if err != nil {
		return nil, fmt.Errorf("invalid argument format. should be yyyy-mm: %v", err)
	}

	ms, err := models.ElapsedOfMonthByDesc(k.db, yyyymm)
	if err != nil {
		return nil, err
	}

	s := make([]Elapsed, len(ms))
	for i := range ms {
		s[i].Tag = ms[i].Tag
		s[i].Desc = ms[i].Desc
		s[i].Count = ms[i].Count
		s[i].Elapsed = ms[i].Elapsed
	}

	return s, nil
}

// AddTags adds a new tags
func (k *Kokizami) AddTags(labels []string) error {
	return k.TagRepo.InsertTags(labels)
}

// DeleteTag deletes a specified tag
func (k *Kokizami) DeleteTag(id int) error {
	return k.TagRepo.Delete(id)
}

// Tags returns list of tags
func (k *Kokizami) Tags() ([]*Tag, error) {
	ms, err := k.TagRepo.FindAllTags()
	if err != nil {
		return nil, err
	}

	ts := make([]*Tag, len(ms))
	for i := range ms {
		ts[i].ID = ms[i].ID
		ts[i].Tag = ms[i].Tag
	}

	return ts, nil
}

// Tagging makes relation between specified kizami and tags
func (k *Kokizami) Tagging(kizamiID int, tagIDs []int) error {
	rs := models.Relations(make([]models.Relation, len(tagIDs)))
	for i := range rs {
		rs[i].KizamiID = kizamiID
		rs[i].TagID = tagIDs[i]

	}
	return rs.BulkInsert(k.db)
}

// Untagging removes all tags from specified kizami
func (k *Kokizami) Untagging(kizamiID int) error {
	return models.DeleteRelationsByKizamiID(k.db, kizamiID)
}

// TagsByKizamiID returns tags of specified kizami
func (k *Kokizami) TagsByKizamiID(kizamiID int) ([]*Tag, error) {
	ms, err := k.TagRepo.FindTagsByKizamiID(kizamiID)
	if err != nil {
		return nil, err
	}

	ts := make([]*Tag, len(ms))
	for i := range ms {
		ts[i].ID = ms[i].ID
		ts[i].Tag = ms[i].Tag
	}

	return ts, nil
}

// TagsByLabels returns tags by specified tags
func (k *Kokizami) TagsByLabels(labels []string) ([]*Tag, error) {
	return k.TagRepo.FindTagsByLabels(labels)
}
