package pipeline

import (
	"context"
	"my-go-learning/logstream/internal/types"
	"reflect"
)

type Filter struct {
	FieldName   string
	TargetValue string
}

func (f *Filter) FilterLogs(ctx context.Context, logs <-chan types.LogEntry) <-chan types.LogEntry {
	fLogs := make(chan types.LogEntry)
	go func() {
		for {

			select {
			case <-ctx.Done():
				close(fLogs)
				return

			case l, ok := <-logs:
				if !ok {
					return
				}
				valueL := reflect.ValueOf(l)
				levelVal := valueL.FieldByName(f.FieldName)

				if levelVal.IsValid() {
					if levelVal.String() == f.TargetValue {
						fLogs <- l
					}
				}
			}

		}

		close(fLogs)
	}()
	return fLogs
}
