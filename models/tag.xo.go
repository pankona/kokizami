// Package models contains the types for schema ''.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// Tag represents a row from 'tag'.
type Tag struct {
	ID    int    `json:"id"`    // id
	Label string `json:"label"` // label

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Tag exists in the database.
func (t *Tag) Exists() bool {
	return t._exists
}

// Deleted provides information if the Tag has been deleted from the database.
func (t *Tag) Deleted() bool {
	return t._deleted
}

// Insert inserts the Tag to the database.
func (t *Tag) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO tag (` +
		`label` +
		`) VALUES (` +
		`?` +
		`)`

	// run query
	XOLog(sqlstr, t.Label)
	res, err := db.Exec(sqlstr, t.Label)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	t.ID = int(id)
	t._exists = true

	return nil
}

// Update updates the Tag in the database.
func (t *Tag) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if t._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE tag SET ` +
		`label = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, t.Label, t.ID)
	_, err = db.Exec(sqlstr, t.Label, t.ID)
	return err
}

// Save saves the Tag to the database.
func (t *Tag) Save(db XODB) error {
	if t.Exists() {
		return t.Update(db)
	}

	return t.Insert(db)
}

// Delete deletes the Tag from the database.
func (t *Tag) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return nil
	}

	// if deleted, bail
	if t._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM tag WHERE id = ?`

	// run query
	XOLog(sqlstr, t.ID)
	_, err = db.Exec(sqlstr, t.ID)
	if err != nil {
		return err
	}

	// set deleted
	t._deleted = true

	return nil
}

// TagByLabel retrieves a row from 'tag' as a Tag.
//
// Generated from index 'sqlite_autoindex_tag_1'.
func TagByLabel(db XODB, label string) (*Tag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, label ` +
		`FROM tag ` +
		`WHERE label = ?`

	// run query
	XOLog(sqlstr, label)
	t := Tag{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, label).Scan(&t.ID, &t.Label)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// TagByID retrieves a row from 'tag' as a Tag.
//
// Generated from index 'tag_id_pkey'.
func TagByID(db XODB, id int) (*Tag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, label ` +
		`FROM tag ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	t := Tag{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&t.ID, &t.Label)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
