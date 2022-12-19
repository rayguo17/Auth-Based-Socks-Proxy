package logger

import "io"

type Config struct {
	DebugPath  string
	AccessPath string
	LogWriter  io.Writer
	IsMulti    bool
}
