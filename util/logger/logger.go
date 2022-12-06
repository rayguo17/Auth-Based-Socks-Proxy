package logger

import (
	"log"
	"os"
)

var (
	Access *Logger
	Debug  *Logger
)

type Logger struct {
	logger *log.Logger
	dst    int
}

//LogType

//LOGDst
const (
	LOG_NIL  int = 0
	LOG_STD      = 1
	LOG_FILE     = 2
)

func (l *Logger) IsNull() bool {
	return l.dst == LOG_NIL
}
func (l *Logger) Println(args ...any) {
	if l.IsNull() {
		return
	}
	l.logger.Println(args)

}
func (l *Logger) Fatal(v ...any) {
	if l.IsNull() {
		return
	}
	l.logger.Fatal(v)
}
func (l *Logger) Print(args ...any) {
	if l.IsNull() {
		return
	}
	l.logger.Println(args)
}
func (l *Logger) Printf(format string, args ...any) {
	if l.IsNull() {
		return
	}
	l.logger.Printf(format, args)
}

//could be nil...
//For Debug Log, we want it to log file location and code line.
//For Error Log, we want it to log time. time username id destination
func InitializeLogger(AccessPath string, DebugPath string) error {
	switch AccessPath {
	case "":
		Access = &Logger{dst: LOG_NIL}
	case "STDOUT":
		Access = &Logger{
			logger: log.New(os.Stdout, "Access: ", log.Ldate|log.Ltime),
			dst:    LOG_STD,
		}
	default:
		file, err := os.OpenFile(AccessPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Access = &Logger{
			logger: log.New(file, "", log.Ldate|log.Ltime),
			dst:    LOG_FILE,
		}

	}
	switch DebugPath {
	case "":
		Debug = &Logger{dst: LOG_NIL}
	case "STDOUT":
		Debug = &Logger{
			logger: log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile),
			dst:    LOG_STD,
		}
	default:
		file, err := os.OpenFile(DebugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Debug = &Logger{
			logger: log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile),
			dst:    LOG_FILE,
		}
	}
	return nil
}
