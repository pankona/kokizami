package todo

import (
	"database/sql"
	"strconv"
	"time"

	// go-sqlite3 is only imported here
	_ "github.com/mattn/go-sqlite3"
)

// DBInterface represents interface of DB
type DBInterface interface {
	openDB() error
	close()
	createTable() error
	start(desc string) (*ToDo, error)
	get(id int) (*ToDo, error)
	edit(id int, field, newValue string) (*ToDo, error)
	list() ([]*ToDo, error)
	stop(id int) error
	stopall() error
	delete(id int) error
}

// DB represents DB instance
type DB struct {
	DBInterface
	conn   *sql.DB
	dbpath string
}

func newDB(dbpath string) *DB {
	return &DB{dbpath: dbpath}
}

func (db *DB) openDB() error {
	conn, err := sql.Open("sqlite3", db.dbpath)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *DB) close() {
	db.conn.Close()
}

func (db *DB) createTable() error {
	q := "CREATE TABLE todo ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", desc VARCHAR(255) NOT NULL"
	q += ", started_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))"
	q += ", stopped_at TIMESTAMP DEFAULT (DATETIME('1970-01-01'))"
	q += ")"

	_, err := db.conn.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) start(desc string) (*ToDo, error) {
	// FIXME: this should not be done here
	q := "UPDATE todo " +
		"SET stopped_at = (DATETIME('now','localtime')) " +
		"WHERE stopped_at = (DATETIME('1970-01-01'))"
	_, err := db.conn.Exec(q)
	if err != nil {
		return nil, err
	}

	q = "INSERT INTO todo (desc) " +
		"VALUES ('" + desc + "')"
	result, err := db.conn.Exec(q)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	q = "SELECT id, desc, started_at, stopped_at " +
		"FROM todo WHERE id = ?"
	t := &ToDo{}
	err = db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) get(id int) (*ToDo, error) {
	q := "SELECT id, desc, started_at, stopped_at " +
		"FROM todo WHERE id = ?"
	t := &ToDo{}
	err := db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) edit(id int, field, newValue string) (*ToDo, error) {
	q := "UPDATE todo " +
		"SET " + field + " = '" + newValue + "' " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	if err != nil {
		return nil, err
	}

	q = "SELECT id, desc, started_at, stopped_at " +
		"FROM todo WHERE id = ?"
	t := &ToDo{}
	err = db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) list() ([]*ToDo, error) {
	q := "SELECT id, desc, started_at, stopped_at " +
		"FROM todo"

	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, err
	}

	todos := make([]*ToDo, 0, 0)
	var id int
	var desc string
	var startedAt time.Time
	var stoppedAt time.Time
	for rows.Next() {
		err = rows.Scan(&id, &desc, &startedAt, &stoppedAt)
		if err != nil {
			panic(err.Error())
		}
		todos = append(todos,
			&ToDo{
				id:        id,
				desc:      desc,
				startedAt: startedAt,
				stoppedAt: stoppedAt,
			})
	}
	return todos, nil
}

func (db *DB) stop(id int) error {
	q := "UPDATE todo " +
		"SET stopped_at = (DATETIME('now','localtime')) " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) stopall() error {
	q := "UPDATE todo " +
		"SET stopped_at = (DATETIME('now','localtime')) " +
		"WHERE stopped_at = (DATETIME('1970-01-01'))"
	_, err := db.conn.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) delete(id int) error {
	q := "DELETE from todo " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
