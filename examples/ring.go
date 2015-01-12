// +build local

package main

import (
	"fmt"
	"github.com/simonwaldherr/golibs/as"
	"github.com/simonwaldherr/golibs/cli"
	"github.com/simonwaldherr/golibs/stack"
	"log"
	"math/rand"
	"time"
)

func main() {
	r := stack.Ring()
	r.Init()
	r.SetSize(15)
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
			rand.Seed(time.Now().UnixNano())
			ran := as.String(rand.Intn(127))
			pos := r.Push(ran)
			log.Printf("Value %v added at %v\n", cli.Bold(cli.Color(ran, cli.Red)), cli.Bold(cli.Color(as.String(pos), cli.Blue)))
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
			log.Printf("Value %v added at %v\n", cli.Bold(cli.Color(ix, cli.Green)), cli.Bold(cli.Color(as.String(nlast), cli.Blue)))
			log.Printf("Get from %v\n", last)
			vals = r.Get(last)
			for _, v := range vals {
				if v != "" {
					log.Printf("\t%v\n", cli.Color(as.String(v), cli.Yellow))
				}
			}
			last = nlast
		}
	}
}
