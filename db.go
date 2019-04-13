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
	dbpath string
}

func newDB(dbpath string) *DB {
	return &DB{dbpath: dbpath}
}

func (db *DB) execWithFunc(f func(conn *sql.DB) error) error {
	conn, err := sql.Open("sqlite3", db.dbpath)
	if err != nil {
		return err
	}
	defer func() {
		e := conn.Close()
		if e != nil {
			log.Printf("failed to close DB connection: %v", e)
		}
	}()

	return f(conn)
}

func (db *DB) createTable() error {
	return db.execWithFunc(func(conn *sql.DB) error {
		q := "CREATE TABLE IF NOT EXISTS todo (" +
			" id INTEGER PRIMARY KEY AUTOINCREMENT" +
			", desc VARCHAR(255) NOT NULL" +
			", started_at TIMESTAMP DEFAULT (DATETIME('now'))" +
			", stopped_at TIMESTAMP DEFAULT (DATETIME('1970-01-01'))" +
			")"
		_, err := conn.Exec(q)
		if err != nil {
			return err
		}

		q = "CREATE TABLE IF NOT EXISTS tag (" +
			" id INTEGER PRIMARY KEY AUTOINCREMENT" +
			", tag VARCHAR(255) NOT NULL" +
			")"
		_, err = conn.Exec(q)
		if err != nil {
			return err
		}

		q = "CREATE TABLE IF NOT EXISTS relation (" +
			" kizami_id INTEGER NOT NULL" +
			")"
		_, err = conn.Exec(q)
		return err
	})
}

func (db *DB) start(desc string) (*kizami, error) {
	k := &kizami{}
	return k, db.execWithFunc(func(conn *sql.DB) error {
		q := "UPDATE todo " +
			"SET stopped_at = (DATETIME('now')) " +
			"WHERE stopped_at = (DATETIME('1970-01-01'))"
		_, err := conn.Exec(q)
		if err != nil {
			return err
		}

		q = "INSERT INTO todo (desc) " +
			"VALUES ('" + desc + "')"
		result, err := conn.Exec(q)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		q = "SELECT id, desc, started_at, stopped_at " +
			"FROM todo WHERE id = ?"
		return conn.QueryRow(q, id).Scan(&k.id, &k.desc, &k.startedAt, &k.stoppedAt)
	})
}

func (db *DB) get(id int) (*kizami, error) {
	k := &kizami{}
	return k, db.execWithFunc(func(conn *sql.DB) error {
		q := "SELECT id, desc, started_at, stopped_at " +
			"FROM todo WHERE id = ?"
		return conn.QueryRow(q, id).Scan(&k.id, &k.desc, &k.startedAt, &k.stoppedAt)
	})
}

func (db *DB) edit(id int, field, newValue string) (*kizami, error) {
	k := &kizami{}
	return k, db.execWithFunc(func(conn *sql.DB) error {
		q := "UPDATE todo " +
			"SET " + field + " = '" + newValue + "' " +
			"WHERE id = " + strconv.Itoa(id)
		_, err := conn.Exec(q)
		if err != nil {
			return err
		}

		q = "SELECT id, desc, started_at, stopped_at " +
			"FROM todo WHERE id = ?"
		return conn.QueryRow(q, id).Scan(&k.id, &k.desc, &k.startedAt, &k.stoppedAt)
	})
}

func (db *DB) count() (int, error) {
	var n int
	return n, db.execWithFunc(func(conn *sql.DB) error {
		q := "SELECT count(*) " +
			"FROM todo"
		return conn.QueryRow(q).Scan(&n)
	})
}

func (db *DB) list(start, end int) ([]*kizami, error) {
	var ks []*kizami
	return ks, db.execWithFunc(func(conn *sql.DB) error {
		if end < start {
			return errors.New("end must be bigger than start")
		}
		s := strconv.Itoa(start)
		c := strconv.Itoa(end - start)
		q := "SELECT id, desc, started_at, stopped_at " +
			"FROM todo " +
			"ORDER BY started_at ASC " +
			"LIMIT " + s + ", " + c

		rows, err := conn.Query(q)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("query was: %s", q))
		}

		todos := make([]*kizami, 0)
		var (
			id        int
			desc      string
			startedAt time.Time
			stoppedAt time.Time
		)
		for rows.Next() {
			err = rows.Scan(&id, &desc, &startedAt, &stoppedAt)
			if err != nil {
				return err
			}
			todos = append(todos,
				&kizami{
					id:        id,
					desc:      desc,
					startedAt: startedAt,
					stoppedAt: stoppedAt,
				})
		}
		ks = todos
		return nil
	})
}

func (db *DB) stop(id int) error {
	return db.execWithFunc(func(conn *sql.DB) error {
		q := "UPDATE todo " +
			"SET stopped_at = (DATETIME('now')) " +
			"WHERE id = " + strconv.Itoa(id)
		_, err := conn.Exec(q)
		return err
	})
}

func (db *DB) stopall() error {
	return db.execWithFunc(func(conn *sql.DB) error {
		q := "UPDATE todo " +
			"SET stopped_at = (DATETIME('now')) " +
			"WHERE stopped_at = (DATETIME('1970-01-01'))"
		_, err := conn.Exec(q)
		return err
	})
}

func (db *DB) delete(id int) error {
	return db.execWithFunc(func(conn *sql.DB) error {
		q := "DELETE from todo " +
			"WHERE id = " + strconv.Itoa(id)
		_, err := conn.Exec(q)
		return err
	})
}
