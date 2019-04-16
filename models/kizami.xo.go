// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	"github.com/xo/xoutil"
)

// Kizami represents a row from 'kizami'.
type Kizami struct {
	ID        int           `json:"id"`         // id
	Desc      string        `json:"desc"`       // desc
	StartedAt xoutil.SqTime `json:"started_at"` // started_at
	StoppedAt xoutil.SqTime `json:"stopped_at"` // stopped_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Kizami exists in the database.
func (k *Kizami) Exists() bool {
	return k._exists
}

// Deleted provides information if the Kizami has been deleted from the database.
func (k *Kizami) Deleted() bool {
	return k._deleted
}

// Insert inserts the Kizami to the database.
func (k *Kizami) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if k._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO kizami (` +
		`desc, started_at, stopped_at` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, k.Desc, k.StartedAt, k.StoppedAt)
	res, err := db.Exec(sqlstr, k.Desc, k.StartedAt, k.StoppedAt)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	k.ID = int(id)
	k._exists = true

	return nil
}

// Update updates the Kizami in the database.
func (k *Kizami) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !k._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if k._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE kizami SET ` +
		`desc = ?, started_at = ?, stopped_at = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, k.Desc, k.StartedAt, k.StoppedAt, k.ID)
	_, err = db.Exec(sqlstr, k.Desc, k.StartedAt, k.StoppedAt, k.ID)
	return err
}

// Save saves the Kizami to the database.
func (k *Kizami) Save(db XODB) error {
	if k.Exists() {
		return k.Update(db)
	}

	return k.Insert(db)
}

// Delete deletes the Kizami from the database.
func (k *Kizami) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !k._exists {
		return nil
	}

	// if deleted, bail
	if k._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM kizami WHERE id = ?`

	// run query
	XOLog(sqlstr, k.ID)
	_, err = db.Exec(sqlstr, k.ID)
	if err != nil {
		return err
	}

	// set deleted
	k._deleted = true

	return nil
}

// KizamiByID retrieves a row from 'kizami' as a Kizami.
//
// Generated from index 'kizami_id_pkey'.
func KizamiByID(db XODB, id int) (*Kizami, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, desc, started_at, stopped_at ` +
		`FROM kizami ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	k := Kizami{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&k.ID, &k.Desc, &k.StartedAt, &k.StoppedAt)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

// KizamisByStoppedAt retrieves a row from 'kizami' as a Kizami.
//
// Generated from index 'stopped_at_index'.
func KizamisByStoppedAt(db XODB, stoppedAt xoutil.SqTime) ([]*Kizami, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, desc, started_at, stopped_at ` +
		`FROM kizami ` +
		`WHERE stopped_at = ?`

	// run query
	XOLog(sqlstr, stoppedAt)
	q, err := db.Query(sqlstr, stoppedAt)
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
