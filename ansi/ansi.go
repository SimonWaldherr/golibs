package ansi

import (
	"fmt"
	"log"
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
func Color(str interface{}, col Col) string {
	return fmt.Sprintf("\033[3%vm%v\033[0m", col, str)
}

// Bold surrounds str with the code for bold styled text
func Bold(str interface{}) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

// Underline surrounds str with the code for underlined text
func Underline(str interface{}) string {
	return fmt.Sprintf("\033[4m%v\033[0m", str)
}

func Log(valuea ...interface{}) interface{} {
	if valuea[1] != nil {
		log.Println(Color(valuea[1], Red))
	}

	return valuea[0]
}
