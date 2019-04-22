package models

import (
	"fmt"
	"time"
)

func CreateKizamiTable(db XODB) error {
	var err error

	// sql query
	sqlstr := "CREATE TABLE IF NOT EXISTS kizami (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", desc VARCHAR(255) NOT NULL" +
		", started_at TIMESTAMP DEFAULT (DATETIME('now'))" +
		", stopped_at TIMESTAMP DEFAULT (DATETIME('1970-01-01'))" +
		")"
	XOLog(sqlstr)
	_, err = db.Exec(sqlstr)
	if err != nil {
		return err
	}

	sqlstr = "CREATE INDEX IF NOT EXISTS index_stopped_at ON kizami(stopped_at)"
	_, err = db.Exec(sqlstr)
	return err
}

func AllKizami(db XODB) ([]*Kizami, error) {
	var err error

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
	defer q.Close()

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

// Elapsed returns kizami's elapsed time
func (k *Kizami) Elapsed() time.Duration {
	var elapsed time.Duration
	if k.StoppedAt.Unix() == 0 {
		// this Kizami is on going. Show elapsed time until now.
		now := time.Now().UTC()
		elapsed = now.Sub(k.StartedAt.Time)
	} else {
		elapsed = k.StoppedAt.Sub(k.StartedAt.Time)
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

func ElapsedWithQuery(db XODB, yyyymm string) ([]*Elapsed, error) {
	sqlstr := fmt.Sprintf(`SELECT `+
		`desc, COUNT(desc), SUM(strftime('%%s', stopped_at) - strftime('%%s', started_at)) AS elapsed `+
		`FROM kizami `+
		`WHERE started_at LIKE '%s-%%' `+
		`GROUP BY desc`, yyyymm)
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	res := []*Elapsed{}
	var sec int64
	for q.Next() {
		e := Elapsed{}

		// scan
		err = q.Scan(&e.Desc, &e.Count, &sec)
		if err != nil {
			return nil, err
		}
		e.Elapsed = time.Duration(sec) * time.Second

		res = append(res, &e)
	}

	return res, nil
}
