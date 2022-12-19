package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogWriter struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) *LogWriter {
	return &LogWriter{ctx: ctx}
}
func (l *LogWriter) Write(p []byte) (int, error) {
	str := string(p)
	runtime.EventsEmit(l.ctx, "log", str)
	return len(str), nil
}
