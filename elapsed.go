package kokizami

import (
	"time"

	"github.com/pankona/kokizami/models"
)

type Elapsed struct {
	Desc    string
	Count   int
	Elapsed time.Duration
}

func (e *Elapsed) toModel() *models.Elapsed {
	return (*models.Elapsed)(e)
}

func toElapsed(m *models.Elapsed) *Elapsed {
	return (*Elapsed)(m)
}
