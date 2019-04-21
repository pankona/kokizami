package models

func CreateTagTable(db XODB) error {
	var err error

	// sql query
	const sqlstr = "CREATE TABLE IF NOT EXISTS tag (" +
		" id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL" +
		", tag VARCHAR(255) NOT NULL" +
		")"
	XOLog(sqlstr)
	_, err = db.Exec(sqlstr)
	return err
}
