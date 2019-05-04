package models

import "bytes"

type Tags []Tag

// BulkInsert inserts multiple tags at once
func (ts Tags) BulkInsert(db XODB) error {
	buf := bytes.NewBuffer([]byte{})

	q1 := []byte("INSERT INTO tag(tag)")
	_, err := buf.Write(q1)
	if err != nil {
		return err
	}

	q2 := []byte(" SELECT ? AS tag")
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

		args[i] = v.Tag
	}

	// run query
	sqlstr := buf.String()
	XOLog(sqlstr, args...)
	_, err = db.Exec(sqlstr, args...)

	return err
}
