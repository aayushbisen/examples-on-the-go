package main

import (
	"context"
	"fmt"
	"my-go-learning/logstream/internal/pipeline"
	"os"
	"runtime/pprof"
	"sync"
)

func main() {
	file, err := os.Create("cpu.prof")
	if err != nil {
		return
	}
	defer file.Close()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
	fmt.Println("===LOGSTREAM STARTED===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := pipeline.Producer{}

	rawLogs := p.Produce(ctx)

	t := pipeline.Transformer{WorkerCount: 5}

	transformedLogs := t.Transform(ctx, rawLogs)

	f := pipeline.Filter{
		FieldName:   "Level",
		TargetValue: "ERROR",
	}

	filteredLogs := f.FilterLogs(ctx, transformedLogs)

	c := pipeline.Consumer{}

	var wg sync.WaitGroup

	wg.Add(1)
	go c.Consume(filteredLogs, &wg)
	wg.Wait()
}
