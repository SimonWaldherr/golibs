// +build local

package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"simonwaldherr.de/go/golibs/ansi"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/cache"
	"time"
)

func runt() {
	fmt.Printf("NumCPU: %v\n", runtime.NumCPU())
	fmt.Printf("NumCgoCall: %v\n", runtime.NumCgoCall())
	fmt.Printf("NumGoroutine: %v\n", runtime.NumGoroutine())
	fmt.Printf("Version: %v\n", runtime.Version())
}

func size(c *cache.Cache) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	fmt.Printf("Cache size: %v\n", ansi.Color(fmt.Sprintf("%v", c.Size()), ansi.Red))
	fmt.Printf("Memory Acquired: %v\n", ansi.Color(fmt.Sprintf("%v", m.Sys), ansi.Blue))
	fmt.Printf("Memory Used: %v\n", ansi.Color(fmt.Sprintf("%v", m.Alloc), ansi.Yellow))
}

func main() {
	c := cache.New(10*time.Second, 1000*time.Millisecond)
	size(c)

	log.Println("writing first 500 values")
	for i := 0; i < 500; i++ {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

		rand.Seed(time.Now().UnixNano())
		ran := as.String(rand.Intn(999))

		c.Add(fmt.Sprintf("Value %v", i), ran)
	}
	size(c)

	log.Println("delete random values")
	for i := 0; i < 100; i++ {
		c.Delete(fmt.Sprintf("Value %v", c.Get(fmt.Sprintf("Value %v", i))))
	}
	size(c)

	log.Println("writing next 500 values")
	for i := 500; i < 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

		rand.Seed(time.Now().UnixNano())
		ran := as.String(rand.Intn(999))

		c.Add(fmt.Sprintf("Value %v", i), ran)
	}
	size(c)

	log.Println("delete random values")
	for i := 200; i < 300; i++ {
		c.Delete(fmt.Sprintf("Value %v", c.Get(fmt.Sprintf("Value %v", i))))
	}
	size(c)

	log.Println("output all values")
	for i := 0; i < 1000; i++ {
		v := c.Get(fmt.Sprintf("Value %v", i))
		if v == nil {
			fmt.Printf("%v: \t%v\t", i, v)
		} else {
			fmt.Printf("%v: \t%v\t", i, ansi.Color(fmt.Sprintf("%v", v), ansi.Green))
		}

		if i%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
	size(c)
	log.Println("wait 11 seconds")
	time.Sleep(time.Duration(11) * time.Second)
	size(c)
	log.Println("run gc")
	runtime.GC()
	size(c)
	log.Println("wait 11 seconds")
	time.Sleep(time.Duration(11) * time.Second)
	size(c)

	runt()
}
