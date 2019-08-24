package models

import "fmt"

// CreateRelationTable creates table for relation model
func CreateRelationTable(db XODB) error {
	// sql query
	const sqlstr = "CREATE TABLE IF NOT EXISTS relation (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", kizami_id INTEGER NOT NULL" +
		", tag_id INTEGER NOT NULL" +
		", UNIQUE(kizami_id, tag_id) ON CONFLICT IGNORE" +
		")"
	XOLog(sqlstr)
	_, err := db.Exec(sqlstr)
	return err
}

// TagsByKizamiID returns tags related to specified kizami
func TagsByKizamiID(db XODB, kizamiID int) ([]*Tag, error) {
	// sql query
	const sqlstr = `SELECT tag.id, tag.tag` +
		` FROM relation` +
		` INNER JOIN tag` +
		` ON relation.tag_id = tag.id` +
		` WHERE kizami_id = ?`
	// run query
	XOLog(sqlstr, kizamiID)
	q, err := db.Query(sqlstr, kizamiID)
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
		err = q.Scan(&t.ID, &t.Label)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// DeleteRelationsByKizamiID removes all tags from specified kizami
func DeleteRelationsByKizamiID(db XODB, kizamiID int) error {
	const sqlstr = `DELETE FROM` +
		` relation` +
		` WHERE kizami_id = ?`
	XOLog(sqlstr, kizamiID)
	_, err := db.Exec(sqlstr, kizamiID)
	return err
}
