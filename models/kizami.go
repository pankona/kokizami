package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateKizamiTable creates table and index for Kizami model
func CreateKizamiTable(db XODB) error {
	// sql query
	sqlstr := "CREATE TABLE IF NOT EXISTS kizami (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", desc VARCHAR(255) NOT NULL" +
		", started_at TIMESTAMP DEFAULT (DATETIME('now'))" +
		", stopped_at TIMESTAMP DEFAULT (DATETIME('1970-01-01'))" +
		")"
	XOLog(sqlstr)
	_, err := db.Exec(sqlstr)
	if err != nil {
		return err
	}

	sqlstr = "CREATE INDEX IF NOT EXISTS index_stopped_at ON kizami(stopped_at)"
	_, err = db.Exec(sqlstr)
	return err
}

// AllKizami returns all Kizami from kizami table
func AllKizami(db XODB) ([]*Kizami, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, desc, started_at, stopped_at ` +
		`FROM kizami`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := q.Close()
		if e != nil {
			XOLog(fmt.Sprintf("failed close query: %v", e))
		}
	}()

	// load results
	res := []*Kizami{}
	for q.Next() {
		k := Kizami{
			_exists: true,
		}

		// scan
		err = q.Scan(&k.ID, &k.Desc, &k.StartedAt, &k.StoppedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &k)
	}

	return res, nil
}

// Elapsed represents elapsed time, that are
// calculated from all kizami items with specified term
type Elapsed struct {
	Tag     string
	Desc    string
	Count   int
	Elapsed time.Duration
}

func elapsedOfMonthBy(db XODB, yyyymm string, groupBy string) ([]*Elapsed, error) {
	sqlstr := fmt.Sprintf(`SELECT `+
		`tag, desc, count(desc), SUM(strftime('%%s', kizami.stopped_at) - strftime('%%s', kizami.started_at)) AS elapsed `+
		`FROM kizami `+
		`LEFT JOIN relation ON kizami.id = relation.kizami_id `+
		`LEFT JOIN tag      ON tag.id    = relation.tag_id `+
		`WHERE started_at LIKE '%s-%%' AND stopped_at NOT LIKE '1970-%%' `+
		`GROUP BY %s`, yyyymm, groupBy) // #nosec
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := q.Close()
		if e != nil {
			XOLog(fmt.Sprintf("failed close query: %v", e))
		}
	}()

	res := []*Elapsed{}
	var (
		sec int64
		tag sql.NullString
	)
	for q.Next() {
		e := Elapsed{}

		err = q.Scan(&tag, &e.Desc, &e.Count, &sec)
		if err != nil {
			return nil, err
		}

		e.Tag = ""
		// #nosec
		if v, _ := tag.Value(); v != nil {
			e.Tag = v.(string)
		}
		e.Elapsed = time.Duration(sec) * time.Second

		res = append(res, &e)
	}

	return res, nil
}

// ElapsedOfMonthByDesc returns each all kizami's total elapsed time
// elapsed in specified month group by desc and tag
func ElapsedOfMonthByDesc(db XODB, yyyymm string) ([]*Elapsed, error) {
	return elapsedOfMonthBy(db, yyyymm, "desc, tag")
}

// ElapsedOfMonthByTag returns each all kizami's total
// elapsed time elapsed in specified month group by tag
func ElapsedOfMonthByTag(db XODB, yyyymm string) ([]*Elapsed, error) {
	return elapsedOfMonthBy(db, yyyymm, "tag")
}
