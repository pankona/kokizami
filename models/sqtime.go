package models

import (
	"time"

	"github.com/xo/xoutil"
)

// SqTime converts time.Time to xoutil.SqTime
func SqTime(t time.Time) xoutil.SqTime {
	return xoutil.SqTime{Time: t}
}
