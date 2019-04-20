package models

import "time"

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
