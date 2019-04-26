package kokizami

import (
	"time"
)

type Elapsed struct {
	Desc    string
	Count   int
	Elapsed time.Duration
}
