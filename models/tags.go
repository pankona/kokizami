package models

import (
	"bytes"
	"fmt"
)

// Tags represents a list of tags
type Tags []Tag

// BulkInsert inserts multiple tags at once
func (ts Tags) BulkInsert(db XODB) error {
	if len(ts) == 0 {
		return nil
	}

	buf := bytes.NewBuffer([]byte{})

	q1 := []byte("INSERT INTO tag(label)")
	_, err := buf.Write(q1)
	if err != nil {
		return err
	}

	q2 := []byte(" SELECT ? AS label")
	q3 := []byte(" UNION SELECT ?")

	args := make([]interface{}, len(ts))
	for i, v := range ts {
		if i == 0 {
			_, err = buf.Write(q2)
		} else {
			_, err = buf.Write(q3)
		}
		if err != nil {
			return err
		}

		args[i] = v.Label
	}

	// run query
	sqlstr := buf.String()
	XOLog(sqlstr, args...)
	_, err = db.Exec(sqlstr, args...)

	return err
}

// TagsByTags retrieves tags by tags at once
func TagsByLabels(db XODB, tags []string) ([]*Tag, error) {
	if len(tags) == 0 {
		return []*Tag{}, nil
	}

	buf := bytes.NewBuffer([]byte{})

	q1 := []byte("SELECT * FROM tag")
	_, err := buf.Write(q1)
	if err != nil {
		return nil, err
	}

	q2 := []byte(" WHERE label = ?")
	q3 := []byte(" OR label = ?")

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		if i == 0 {
			_, err = buf.Write(q2)
		} else {
			_, err = buf.Write(q3)
		}
		if err != nil {
			return nil, err
		}

		args[i] = v
	}

	// run query
	sqlstr := buf.String()
	XOLog(sqlstr, args...)
	q, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := q.Close()
		if e != nil {
			XOLog(fmt.Sprintf("failed close query: %v", e))
		}
	}()

	res := []*Tag{}
	for q.Next() {
		t := &Tag{
			_exists: true,
		}

		err = q.Scan(&t.ID, &t.Label)
		if err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	return res, err
}
