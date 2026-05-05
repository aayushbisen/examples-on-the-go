package main

import (
	"fmt"
	"my-go-learning/logstream/internal/pipeline"
	"sync"
)

func main() {
	fmt.Println("===LOGSTREAM STARTED===")

	p := pipeline.Producer{}

	rawLogs := p.Produce()

	t := pipeline.Transformer{}

	filteredLogs := t.Transform(rawLogs)

	c := pipeline.Consumer{}

	var wg sync.WaitGroup

	wg.Add(1)
	go c.Consume(filteredLogs, &wg)
	wg.Wait()
}
