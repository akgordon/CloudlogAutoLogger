package agg_logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type agg_logger struct {
	filePtr *os.File
}

var theLogger *agg_logger = nil

func Get() *agg_logger {
	if theLogger == nil {
		theLogger = &agg_logger{}
	}
	return theLogger
}

func (fl *agg_logger) Open(logName string) {
	// Open or create the log file
	var err error

	fl.filePtr, err = os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Logging to file:%s", logName)
}

func (fl *agg_logger) Close() {
	if fl.filePtr != nil {
		fl.filePtr.Close()
	}
}

func (fl *agg_logger) Log(msg string, msg2 string) {
	msg = strings.ReplaceAll(msg, "\x00", "")
	msg2 = strings.ReplaceAll(msg2, "\x00", "")

	t := time.Now()
	//ts := t.Format("2006-01-02 15:04:05")
	ts := t.Format("15:04:05")

	fmt.Println(ts + " " + msg + " " + msg2)
	fl.filePtr.WriteString(msg + " " + msg2 + "\n")
	//log.Print(msg + " " + msg2)
}
