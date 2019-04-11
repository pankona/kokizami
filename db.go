package kokizami

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
	// go-sqlite3 is only imported here
	_ "github.com/mattn/go-sqlite3"
)

// DBInterface represents interface of DB
type DBInterface interface {
	openDB() error
	close()
	createTable() error
	start(desc string) (*kizami, error)
	get(id int) (*kizami, error)
	edit(id int, field, newValue string) (*kizami, error)
	count() (int, error)
	list(start, end int) ([]*kizami, error)
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
	err := db.conn.Close()
	if err != nil {
		log.Printf("failed to close DB: %v", err)
	}
}

func (db *DB) createTable() error {
	q := "CREATE TABLE todo ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", desc VARCHAR(255) NOT NULL"
	q += ", started_at TIMESTAMP DEFAULT (DATETIME('now'))"
	q += ", stopped_at TIMESTAMP DEFAULT (DATETIME('1970-01-01'))"
	q += ")"

	_, err := db.conn.Exec(q)
	if err != nil {
		return err
	}

	q = "CREATE TABLE tag ("
	q += " id INTEGER PRIMARY KEY AUTOINCREMENT"
	q += ", tag VARCHAR(255) NOT NULL"
	q += ")"

	_, err = db.conn.Exec(q)
	if err != nil {
		return err
	}

	q = "CREATE TABLE relation ("
	q += " kizami_id INTEGER NOT NULL"
	q += " tag_id INTEGER NOT NULL"
	q += ")"

	_, err = db.conn.Exec(q)
	return err
}

func (db *DB) start(desc string) (*kizami, error) {
	// FIXME: this should not be done here
	q := "UPDATE todo " +
		"SET stopped_at = (DATETIME('now')) " +
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
	t := &kizami{}
	err = db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) get(id int) (*kizami, error) {
	q := "SELECT id, desc, started_at, stopped_at " +
		"FROM todo WHERE id = ?"
	t := &kizami{}
	err := db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) edit(id int, field, newValue string) (*kizami, error) {
	q := "UPDATE todo " +
		"SET " + field + " = '" + newValue + "' " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	if err != nil {
		return nil, err
	}

	q = "SELECT id, desc, started_at, stopped_at " +
		"FROM todo WHERE id = ?"
	t := &kizami{}
	err = db.conn.QueryRow(q, id).Scan(&t.id, &t.desc, &t.startedAt, &t.stoppedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (db *DB) count() (int, error) {
	q := "SELECT count(*) " +
		"FROM todo"
	var num int
	err := db.conn.QueryRow(q).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func (db *DB) list(start, end int) ([]*kizami, error) {
	if end < start {
		return nil, errors.New("end must be bigger than start")
	}
	s := strconv.Itoa(start)
	c := strconv.Itoa(end - start)
	q := "SELECT id, desc, started_at, stopped_at " +
		"FROM todo " +
		"ORDER BY started_at ASC " +
		"LIMIT " + s + ", " + c

	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("query was: %s", q))
	}

	todos := make([]*kizami, 0)
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
			&kizami{
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
		"SET stopped_at = (DATETIME('now')) " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	return err
}

func (db *DB) stopall() error {
	q := "UPDATE todo " +
		"SET stopped_at = (DATETIME('now')) " +
		"WHERE stopped_at = (DATETIME('1970-01-01'))"
	_, err := db.conn.Exec(q)
	return err
}

func (db *DB) delete(id int) error {
	q := "DELETE from todo " +
		"WHERE id = " + strconv.Itoa(id)
	_, err := db.conn.Exec(q)
	return err
}
