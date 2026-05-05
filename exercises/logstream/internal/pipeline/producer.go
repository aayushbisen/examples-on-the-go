package pipeline

import (
	"fmt"
	"math/rand"
	"my-go-learning/logstream/internal/types"
	"time"
)

type Producer struct{}

func (p *Producer) Produce() <-chan types.LogEntry {
	rawLogs := make(chan types.LogEntry)

	go func() {
		for i := range 50 {

			levels := []string{"INFO", "WARN", "ERROR"}
			idx := rand.Intn(len(levels))
			level := levels[idx]

			rawLogs <- types.LogEntry{
				Stamp:   time.Now(),
				Level:   level,
				Message: fmt.Sprintf("Log message %d", i),
			}
		}

		close(rawLogs)
	}()

	return rawLogs
}
