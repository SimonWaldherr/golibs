// +build local

package main

import (
	"../file"
	"../cachedfile"
	"log"
	"time"
)

func main() {
	var fn, fs, ca string

	fn = "./test.txt"

	log.Println(file.GetAbsolutePath(fn))
	cachedfile.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	cachedfile.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = cachedfile.CachedRead(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	cachedfile.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)
	cachedfile.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(5 * time.Second)

	fs, _ = file.Read(fn)
	ca, _ = cachedfile.CachedRead(fn)
	log.Printf("\nfile:\t%v\ncached:\t%v\n\n", fs, ca)

	cachedfile.CachedWrite(fn, "\nFoobar\n", true)
	time.Sleep(25 * time.Second)
	cachedfile.StopCache()
}
