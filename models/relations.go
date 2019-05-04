package models

import "bytes"

type Relations []Relation

// BulkInsert inserts multiple relations at once
func (rs Relations) BulkInsert(db XODB) error {
	buf := bytes.NewBuffer([]byte{})

	q1 := []byte("INSERT INTO relation(kizami_id, tag_id)")
	_, err := buf.Write(q1)
	if err != nil {
		return err
	}

	q2 := []byte(" SELECT ? AS kizami_id, ? AS tag_id")
	q3 := []byte(" UNION SELECT ?, ?")

	args := make([]interface{}, len(rs)*2)
	for i, v := range rs {
		if i == 0 {
			_, err = buf.Write(q2)
		} else {
			_, err = buf.Write(q3)
		}
		if err != nil {
			return err
		}

		args[i*2] = v.KizamiID
		args[i*2+1] = v.TagID
	}

	// run query
	sqlstr := buf.String()
	XOLog(sqlstr, args...)
	_, err = db.Exec(sqlstr, args...)

	return err
}
