package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Access *log.Logger
	Debug  *log.Logger
)

//LogType

//LOGDst

//could be nil...
//For Debug Log, we want it to log file location and code line.
//For Error Log, we want it to log time. time username id destination
func InitializeLogger(config Config) error {
	AccessPath := config.AccessPath
	DebugPath := config.DebugPath
	var accessWriter io.Writer
	switch AccessPath {
	case "":
		accessWriter = ioutil.Discard
		Access = log.New(ioutil.Discard, "Access: ", log.Ldate|log.Ltime)
	case "STDOUT":
		accessWriter = os.Stdout
	default:
		file, err := os.OpenFile(AccessPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		accessWriter = file
	}
	if config.IsMulti {
		accessWriter = io.MultiWriter(config.LogWriter, accessWriter)
	}
	Access = log.New(accessWriter, "Access: ", log.Ldate|log.Ltime)

	var debugWriter io.Writer
	switch DebugPath {
	case "":
		debugWriter = ioutil.Discard

	case "STDOUT":
		debugWriter = os.Stdout
	default:
		file, err := os.OpenFile(DebugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		debugWriter = file
	}
	if config.IsMulti {
		debugWriter = io.MultiWriter(debugWriter, config.LogWriter)
	}
	Debug = log.New(debugWriter, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
