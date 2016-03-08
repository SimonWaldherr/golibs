// +build local

package main

import (
	"fmt"
	"simonwaldherr.de/go/golibs/gcurses"
	"strings"
	"time"
)

func ConcatMultiline(str ...string) string {
	outmap := make(map[int]string)
	var out string

	for _, in := range str {
		s := strings.Split(in, "\n")
		var i = 0
		for _, l := range s {
			outmap[i] += l
			i++
		}
	}
	for i := 0; i < len(outmap); i++ {
		out += outmap[i] + "\n"
	}
	return out
}

var BigChars = map[string]string{
	"0": ` ███
 █ █
 █ █
 █ █
 ███`,
	"1": `  █ 
  █ 
  █ 
  █ 
  █ `,
	"2": ` ███
   █
 ███
 █  
 ███`,
	"3": ` ███
   █
  ██
   █
 ███`,
	"4": ` █ █
 █ █
 ███
   █
   █`,
	"5": ` ███
 █  
 ███
   █
 ███`,
	"6": ` █
 █  
 ███
 █ █
 ███`,
	"7": ` ███
   █
   █
   █
   █`,
	"8": ` ███
 █ █
 ███
 █ █
 ███`,
	"9": ` ███
 █ █
 ███
   █
 ███`,
}

func main() {
	writer := gcurses.New()

	writer.Start()

	for i := 0; i < 100; i++ {
		fmt.Fprintf(writer, "Count till one hundred: %d\nStill counting ...\n", i)
		time.Sleep(time.Millisecond * 10)
	}

	time.Sleep(time.Millisecond * 500)
	fmt.Fprintln(writer, "Finished counting")
	time.Sleep(time.Millisecond * 500)

	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("%02d", i)
		fmt.Fprintf(writer, "Counting with big numbers\n%s\n", ConcatMultiline(BigChars[k[0:1]], BigChars[k[1:2]]))

		time.Sleep(time.Millisecond * 40)
	}

	writer.Stop()
}
