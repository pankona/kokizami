package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/models"
	"github.com/xo/xoutil"
)

type kizamiRepo struct {
	db  *sql.DB
	now func() time.Time
}

// SqTime converts time.Time to xoutil.SqTime
func SqTime(t time.Time) xoutil.SqTime {
	return xoutil.SqTime{Time: t}
}

func toKizami(m *models.Kizami) *kokizami.Kizami {
	return &kokizami.Kizami{
		ID:        m.ID,
		Desc:      m.Desc,
		StartedAt: m.StartedAt.Time,
		StoppedAt: m.StoppedAt.Time,
	}
}

func initialTime() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	if err != nil {
		panic(fmt.Sprintf("failed to parse time for initial value for time: %v", err))
	}
	return t.UTC()
}

func (r *kizamiRepo) Insert(desc string) (*kokizami.Kizami, error) {
	m := &models.Kizami{
		Desc:      desc,
		StartedAt: SqTime(r.now().UTC()),
		StoppedAt: SqTime(initialTime()),
	}
	err := m.Insert(r.db)
	if err != nil {
		return nil, err
	}

	return toKizami(m), nil

}

func (r *kizamiRepo) AllKizami() ([]*kokizami.Kizami, error) {
	ms, err := models.AllKizami(r.db)
	if err != nil {
		return nil, err
	}

	ks := make([]kokizami.Kizami, len(ms))
	for i, v := range ms {
		ks[i].ID = v.ID
		ks[i].Desc = v.Desc
		ks[i].StartedAt = v.StartedAt.Time
		ks[i].StoppedAt = v.StoppedAt.Time
	}

	ret := make([]*kokizami.Kizami, len(ms))
	for i := range ms {
		ret[i] = &ks[i]
	}

	return ret, nil
}

func (r *kizamiRepo) Update(k *kokizami.Kizami) error {
	m, err := models.KizamiByID(r.db, k.ID)
	if err != nil {
		return err
	}
	m.Desc = k.Desc
	m.StartedAt = SqTime(k.StartedAt)
	m.StoppedAt = SqTime(k.StoppedAt)

	return m.Update(r.db)
}

func (r *kizamiRepo) Delete(k *kokizami.Kizami) error {
	m, err := models.KizamiByID(r.db, k.ID)
	if err != nil {
		return err
	}
	return m.Delete(r.db)
}

func (r *kizamiRepo) KizamiByID(id int) (*kokizami.Kizami, error) {
	m, err := models.KizamiByID(r.db, id)
	if err != nil {
		return nil, err
	}
	return toKizami(m), nil
}

func (r *kizamiRepo) KizamisByStoppedAt(t time.Time) ([]*kokizami.Kizami, error) {
	ms, err := models.KizamisByStoppedAt(r.db, SqTime(t))
	if err != nil {
		return nil, err
	}

	ks := make([]kokizami.Kizami, len(ms))
	for i, v := range ms {
		ks[i].ID = v.ID
		ks[i].Desc = v.Desc
		ks[i].StartedAt = v.StartedAt.Time
		ks[i].StoppedAt = v.StoppedAt.Time
	}

	ret := make([]*kokizami.Kizami, len(ms))
	for i := range ms {
		ret[i] = &ks[i]
	}

	return ret, nil
}

func (r *kizamiRepo) Tagging(kizamiID int, tagIDs []int) error {
	rs := models.Relations(make([]models.Relation, len(tagIDs)))
	for i := range rs {
		rs[i].KizamiID = kizamiID
		rs[i].TagID = tagIDs[i]
	}
	return rs.BulkInsert(r.db)
}

func (r *kizamiRepo) Untagging(kizamiID int) error {
	return models.DeleteRelationsByKizamiID(r.db, kizamiID)
}
