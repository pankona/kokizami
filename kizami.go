package kokizami

import (
	"time"

	"github.com/pankona/kokizami/models"
)

// Kizami represents a task
type Kizami struct {
	ID        int
	Desc      string
	StartedAt time.Time
	StoppedAt time.Time
}

func toKizami(m *models.Kizami) *Kizami {
	return &Kizami{
		ID:        m.ID,
		Desc:      m.Desc,
		StartedAt: m.StartedAt.Time,
		StoppedAt: m.StoppedAt.Time,
	}
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
