package logger

import (
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
func InitializeLogger(AccessPath string, DebugPath string) error {
	switch AccessPath {
	case "":
		Access = log.New(ioutil.Discard, "Access: ", log.Ldate|log.Ltime)
	case "STDOUT":
		Access = log.New(os.Stdout, "Access: ", log.Ldate|log.Ltime)
	default:
		file, err := os.OpenFile(AccessPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Access = log.New(file, "", log.Ldate|log.Ltime)

	}
	switch DebugPath {
	case "":
		Debug = log.New(ioutil.Discard, "Access: ", log.Ldate|log.Ltime)
	case "STDOUT":
		Debug = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	default:
		file, err := os.OpenFile(DebugPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Debug = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return nil
}
