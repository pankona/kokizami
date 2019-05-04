package models

import "fmt"

// CreateTagTable creates table for tag model
func CreateTagTable(db XODB) error {
	// sql query
	const sqlstr = "CREATE TABLE IF NOT EXISTS tag (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", tag VARCHAR(255) NOT NULL" +
		", UNIQUE(tag) ON CONFLICT IGNORE" +
		")"
	XOLog(sqlstr)
	_, err := db.Exec(sqlstr)
	return err
}

// AllTags returns all tags from tag table
func AllTags(db XODB) ([]*Tag, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, tag ` +
		`FROM tag`

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
	res := []*Tag{}
	for q.Next() {
		t := Tag{
			_exists: true,
		}

		// scan
		err = q.Scan(&t.ID, &t.Tag)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}
