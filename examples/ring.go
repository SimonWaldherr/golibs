// +build local

package main

import (
	"fmt"
	"log"
	"math/rand"
	"simonwaldherr.de/go/golibs/ansi"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/stack"
	"time"
)

func main() {
	r := stack.Ring()
	r.Init(10)
	r.SetSize(15)
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
			rand.Seed(time.Now().UnixNano())
			ran := as.String(rand.Intn(127))
			pos := r.Push(ran)
			log.Printf("Value %v added at %v\n", ansi.Bold(ansi.Color(ran, ansi.Red)), ansi.Bold(ansi.Color(as.String(pos), ansi.Blue)))
		}
	}()
	var last int = 0
	var nlast int = 0
	var vals []string
	var i int

	for {
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			fmt.Println("only numbers")
		} else {
			ix := as.String(i)
			nlast = r.Push(ix)
			log.Printf("Value %v added at %v\n", ansi.Bold(ansi.Color(ix, ansi.Green)), ansi.Bold(ansi.Color(as.String(nlast), ansi.Blue)))
			log.Printf("Get from %v\n", last)
			vals = r.Get(last)
			for _, v := range vals {
				if v != "" {
					log.Printf("\t%v\n", ansi.Color(as.String(v), ansi.Yellow))
				}
			}
			last = nlast
		}
	}
}
