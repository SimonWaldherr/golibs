package log

import (
	olog "log"
	"os"
	"testing"
)

func Test_Log(t *testing.T) {
	Info.Println("1")
	Warning.Println("2")
	Error.Println("3")
	Fatal.Println("4")
}

func Test_Change(t *testing.T) {
	Change(os.Stdout, os.Stdout, os.Stdout, os.Stdout,
		olog.Ltime|olog.Lshortfile, olog.Ltime|olog.Lshortfile, olog.Ltime|olog.Lshortfile, olog.Ltime|olog.Lshortfile)
	Info.Println("1")
	Warning.Println("2")
	Error.Println("3")
	Fatal.Println("4")
}

func Test_Change2(t *testing.T) {
	Change(os.Stdout, os.Stdout, os.Stdout, os.Stdout, 0, 0, 0, 0)
	Info.Println("1")
	Warning.Println("2")
	Error.Println("3")
	Fatal.Println("4")
}
