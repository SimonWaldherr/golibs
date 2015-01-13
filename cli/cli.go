package cli

import (
	"fmt"
)

type Col int

const (
	Black = Col(iota)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func Color(str string, col Col) string {
	return fmt.Sprintf("\033[3%vm%v\033[0m", col, str)
}

func Bold(str string) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

func Underline(str string) string {
	return fmt.Sprintf("\033[4m%v\033[0m", str)
}
