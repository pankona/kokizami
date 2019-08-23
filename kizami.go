package kokizami

import (
	"time"
)

// Kizami represents a task
type Kizami struct {
	ID        int
	Desc      string
	StartedAt time.Time
	StoppedAt time.Time
}

// Elapsed returns kizami's elapsed time
func (k *Kizami) Elapsed() time.Duration {
	var elapsed time.Duration
	if k.StoppedAt.Unix() == 0 {
		// this Kizami is on going. Show elapsed time until now.
		now := time.Now().UTC()
		elapsed = now.Sub(k.StartedAt)
	} else {
		elapsed = k.StoppedAt.Sub(k.StartedAt)
		if elapsed < 0 {
			elapsed = 0
		}
	}
	return elapsed
}

type KizamiRepository interface {
	AllKizami() ([]*Kizami, error)
	Insert(desc string) (*Kizami, error)
	Update(k *Kizami) error
	Delete(k *Kizami) error
	KizamiByID(id int) (*Kizami, error)
	KizamisByStoppedAt(t time.Time) ([]*Kizami, error)
	Tagging(kizamiID int, tagIDs []int) error
}
