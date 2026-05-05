package pipeline

import (
	"fmt"
	"my-go-learning/logstream/internal/types"
	"sync"
)

type Consumer struct{}

func (c *Consumer) Consume(t <-chan types.LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	for entry := range t {
		fmt.Println(entry)
	}

}
