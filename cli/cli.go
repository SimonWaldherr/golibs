package cli

import (
	"fmt"
)

type Col int

// defines supported colors
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

// Color adds the color code of col to str and returns as string
func Color(str string, col Col) string {
	return fmt.Sprintf("\033[3%vm%v\033[0m", col, str)
}

// Bold surrounds str with the code for bold styled text
func Bold(str string) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

// Underline surrounds str with the code for underlined text
func Underline(str string) string {
	return fmt.Sprintf("\033[4m%v\033[0m", str)
}
