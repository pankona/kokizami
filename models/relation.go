package models

func CreateRelationTable(db XODB) error {
	var err error

	// sql query
	const sqlstr = "CREATE TABLE IF NOT EXISTS relation (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", kizami_id INTEGER NOT NULL" +
		", tag_id INTEGER NOT NULL" +
		")"
	XOLog(sqlstr)
	_, err = db.Exec(sqlstr)
	return err
}
