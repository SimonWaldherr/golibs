package log

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stderr, "WARNING: ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

func Change(infoHandle, warningHandle, errorHandle, fatalHandle io.Writer,
	infoFlag, warningFlag, errorFlag, fatalFlat int) {

	Info = log.New(infoHandle,
		"INFO: ", infoFlag)

	Warning = log.New(warningHandle,
		"WARNING: ", warningFlag)

	Error = log.New(errorHandle,
		"ERROR: ", errorFlag)

	Fatal = log.New(fatalHandle,
		"FATAL: ", fatalFlat)
}
