package pipeline

import (
	"context"
	"fmt"
	"my-go-learning/logstream/internal/types"
	"strings"
	"sync"
	"time"
)

type Transformer struct {
	WorkerCount int
}

/**
*  Think about these questions before starting coding:
1. Input: What should the Transform method take as an argument? (Remember the read-only channel concept).
2. Output: What should the Transform method return?
3. The Work: Since the Transformer needs to wait for logs to arrive from the Producer, should the transformation logic happen in the main function or in a
goroutine?
4. Closure: Who is responsible for closing the Transformer's output channel?

Answering some questions before coding any line
1. It should take a channel as arg, func Transformer(rawLogs chan types.LogEntry) {} wrong
2. Transformer method should return the filtered channel
3. Yes since there is a waiting of logs the transform logic should be in go routing
4. The resposible for closing the transformer output channel should be what I think is with consumer since after transforming all the raw messages it should not close its channel so that when consumer consumes it should be alive wrong


Revised Logic for the Transformer:
1. Read from the Producer's channel using a for range loop.
2. Transform the data and send it to the output channel.
3. Once the for range loop finishes (which means the Producer closed its channel), the Transformer should then close its own output channel.
4. This signals to the Consumer: "I'm done sending everything I've processed; you can stop now."
*/

func (t *Transformer) Transform(ctx context.Context, rawLogs <-chan types.LogEntry) <-chan types.LogEntry {
	// reads from that channel, converts the log message to uppercase, and sends it into a second channel.
	transformedLogs := make(chan types.LogEntry)
	var wg sync.WaitGroup

	for i := range t.WorkerCount {
		wg.Add(1)
		fmt.Printf("Starting worker %d", i)
		time.Sleep(100 * time.Millisecond)
		go func() {
			defer wg.Done()
			sum := 0
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker cancelled. Sum: %d\n", sum)
					return
				case logEntry, ok := <-rawLogs:
					if !ok {
						fmt.Printf("Worker finished. Final sum: %d\n", sum)
						return
					}
					for o := range 10000000 {
						sum = sum + o
					}
					upMsg := strings.ToUpper(logEntry.Message)
					transformedLogs <- types.LogEntry{
						Message: upMsg,
						Stamp:   logEntry.Stamp,
						Level:   logEntry.Level,
					}
				}
			}

			// for logEntry := range rawLogs {
			// 	upMsg := strings.ToUpper(logEntry.Message)

			// 	transformedLogs <- types.LogEntry{
			// 		Message: upMsg,
			// 		Stamp:   logEntry.Stamp,
			// 		Level:   logEntry.Level,
			// 	}
			// }
			// close(transformedLogs)
		}()
	}

	go func() {
		wg.Wait()
		close(transformedLogs)
	}()

	return transformedLogs
}
