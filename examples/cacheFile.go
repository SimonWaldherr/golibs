//go:build local
// +build local

package main

import (
	"log"
	"time"

	"simonwaldherr.de/go/golibs/cachedfile"
	"simonwaldherr.de/go/golibs/file"
)

func main() {
	var fn, fs, ca string

	fn = "./test.txt"

	log.Println(file.GetAbsolutePath(fn))
	cachedfile.Write(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	cachedfile.Write(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = cachedfile.Read(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	cachedfile.Write(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	cachedfile.Write(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = cachedfile.Read(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	cachedfile.Write(fn, "\nFoobar\n", true)
	time.Sleep(25 * time.Second)
	cachedfile.Stop()
}
