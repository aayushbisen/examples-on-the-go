package types

import "time"

type LogEntry struct {
	Stamp   time.Time
	Level   string
	Message string
}
