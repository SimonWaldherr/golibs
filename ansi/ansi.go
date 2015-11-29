// Package ansi can print colored and styled text to your terminal.
package ansi

import (
	"fmt"
	"log"
)

// Col defines supported colors
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

// defines text styling options
const (
	TReset = Col(iota)
	TBold
	TFaint
	TItalic
	TUnderline
	TBlinkSlow
	TBlinkFast
	TNegative
	TConceal
	TCrossedOut
)

// Foreground text colors
const (
	FgBlack = Col(iota + 30)
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Background text colors
const (
	BgBlack = Col(iota + 40)
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Set sets styling options on strings and stringable interfaces
func Set(str interface{}, Attribute ...Col) string {
	var rstr = fmt.Sprint(str)
	for _, attr := range Attribute {
		rstr = fmt.Sprintf("\033[%vm%v\033[0m", attr, rstr)
	}
	return rstr
}

// Color adds the color code of col as text color to str and returns as string
func Color(str interface{}, col Col) string {
	return fmt.Sprintf("\033[3%vm%v\033[0m", col, str)
}

// BgColor adds the color code of col as background color to str and returns as string
func BgColor(str interface{}, col Col) string {
	return fmt.Sprintf("\033[4%vm%v\033[0m", col, str)
}

// Bold surrounds str with the code for bold styled text
func Bold(str interface{}) string {
	return fmt.Sprintf("\033[1m%v\033[0m", str)
}

// Underline surrounds str with the code for underlined text
func Underline(str interface{}) string {
	return fmt.Sprintf("\033[4m%v\033[0m", str)
}

// Log prints red text via log package
func Log(valuea ...interface{}) interface{} {
	if valuea[1] != nil {
		log.Println(Color(valuea[1], Red))
	}

	return valuea[0]
}
