package kokizami

import (
	"strconv"
	"time"
)

// Kizamier represents interface of Kizami
type Kizamier interface {
	ID() int
	Desc() string
	StartedAt() time.Time
	StoppedAt() time.Time
	Elapsed() time.Duration
	String() string
}

type kizami struct {
	id        int
	desc      string
	startedAt time.Time
	stoppedAt time.Time
}

// String returns string representation of a kizami.
// note that the timestamps is not considered time zone.
func (k *kizami) String() string {
	return strconv.Itoa(k.id) + "\t" +
		k.desc + "\t" +
		k.startedAt.Format("2006-01-02 15:04:05") + "\t" +
		k.stoppedAt.Format("2006-01-02 15:04:05") + "\t" +
		k.Elapsed().String()
}

// ID returns kizami's id
func (k *kizami) ID() int {
	return k.id
}

// Desc returns kizami's description
func (k *kizami) Desc() string {
	return k.desc
}

// StartedAt returns kizami's startedAt
func (k *kizami) StartedAt() time.Time {
	return k.startedAt
}

// StoppedAt returns kizami's stoppedAt
func (k *kizami) StoppedAt() time.Time {
	return k.stoppedAt
}

// Elapsed returns kizami's elapsed time
func (k *kizami) Elapsed() time.Duration {
	var elapsed time.Duration
	if k.stoppedAt.Unix() == 0 {
		// this kizami is on going. show elapsed time until now.
		now := time.Now().UTC()
		elapsed = now.Sub(k.startedAt)
	} else {
		elapsed = k.stoppedAt.Sub(k.startedAt)
		if elapsed < 0 {
			elapsed = 0
		}
	}
	return elapsed
}
