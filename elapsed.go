package kokizami

import (
	"time"
)

// Elapsed represents elapsed time of each Kizami
type Elapsed struct {
	Desc    string
	Count   int
	Elapsed time.Duration
}
