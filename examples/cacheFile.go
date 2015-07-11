// +build local

package main

import (
	"../file"
	"log"
	"time"
)

func main() {
	var fn, fs, ca string

	fn = "./test.txt"

	log.Println(file.GetAbsolutePath(fn))
	file.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	file.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = file.CachedRead(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	file.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	file.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = file.CachedRead(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	file.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(25 * time.Second)
	file.StopCache()
}
