# golibs


golang functions (to be included in other projects)

```sh
go get -u -t simonwaldherr.de/go/golibs/...
```

## coverage & tests

```sh
go test ./...
```

 . | . | service | info
---|---|---------|------
 ![OSX Build Status](https://simonwaldherr.de/icon/osx.png) | ![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square) |  | go test ./...  
 ![Linux Build Status](https://simonwaldherr.de/icon/tux.png) | [![Travis CI](https://img.shields.io/travis/SimonWaldherr/golibs.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs) | travis-ci.org | test on various go versions  
 ![Windows Build Status](https://simonwaldherr.de/icon/win.png) | [![Appveyor CI](https://img.shields.io/appveyor/ci/SimonWaldherr/golibs.svg?style=flat-square)](https://ci.appveyor.com/project/SimonWaldherr/golibs/branch/master) | appveyor.com | test under windows  
 ![Linux Build Status](https://simonwaldherr.de/icon/tux.png) | [![Semaphore CI](https://semaphoreci.com/api/v1/projects/fe1a7a53-a2c0-4bc0-a539-4af4ef13d49f/487313/shields_badge.svg)](https://semaphoreci.com/simonwaldherr/golibs) | semaphoreci.com | yet another ci service  
 ![Linux Build Status](https://simonwaldherr.de/icon/tux1.png) | [![Wercker CI](https://app.wercker.com/status/c08af8170c3ee4955ca683dd8bf96f07/s/master "wercker status")](https://app.wercker.com/project/bykey/c08af8170c3ee4955ca683dd8bf96f07) | wercker.com | the wercker ci service  
   | [![Go Report Card](http://goreportcard.com/badge/simonwaldherr/golibs)](http://goreportcard.com/report/simonwaldherr/golibs) | goreportcard.com | report card  
   | [![codebeat badge](https://codebeat.co/badges/f697775a-ac7f-4d6f-a7a2-214fadf765d9)](https://codebeat.co/projects/github-com-simonwaldherr-golibs) | codebeat.co | code quality  
   | [![Codacy Badge](https://api.codacy.com/project/badge/Grade/de8a4f17bb7546f7b191f47f1b223b8b)](https://www.codacy.com/app/SimonWaldherr/golibs?utm_source=github.com&utm_medium=referral&utm_content=SimonWaldherr/golibs&utm_campaign=badger) | codacy.com | code quality  
   | [![BCH compliance](https://bettercodehub.com/edge/badge/SimonWaldherr/golibs?branch=master)](https://bettercodehub.com/) | bettercodehub.com | code quality  
   | [![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) | coveralls.io | test coverage  
   | [![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://simonwaldherr.de/gocover/golibs/) |  | go tool cover  
   | [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FSimonWaldherr%2Fgolibs.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FSimonWaldherr%2Fgolibs?ref=badge_shield) | fossa.io | license report  
   | [![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/SimonWaldherr/golibs/master/LICENSE) |  | free + open source license  
   | [![Flattr donate button](https://img.shields.io/badge/flattr-donate-orange.svg?style=flat-square)](https://flattr.com/submit/auto?user_id=SimonWaldherr&url=http%3A%2F%2Fgithub.com%2FSimonWaldherr%2Fgolibs "Donate monthly to this project using Flattr") | flattr.com | micro donation  
   | [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/) | godoc.org | documentation  

each new build gets tested in multiple steps:

* on development i regularly type ```go test``` from time to time to check the test suite
* also there are a few go apps in the [examples](https://github.com/SimonWaldherr/golibs/tree/master/examples)-folder which i test to build and run
* on commit, git automatically runs the [pre-commit](https://github.com/SimonWaldherr/golibs/blob/master/pre-commit)-hook shell script
* after a commit gets pushed to **GitHub**, the following tests are started via Webhooks and Services
	* **Travis CI** build the lib and all tests on docker containers with the go versions noted in [.travis.yml](https://github.com/SimonWaldherr/golibs/blob/master/.travis.yml)
	* **appveyor** builds the lib on Windows Server to test against the Microsoft Infrastructure ([.appveyor.yml](https://github.com/SimonWaldherr/golibs/blob/master/.appveyor.yml))
	* **magnum-ci**, another ci service build the lib and tests on a linux (ubuntu) machine
	* **semaphoreci**, yet another (linux (ubuntu) based) ci service

![codecov.io](https://codecov.io/github/SimonWaldherr/golibs/branch.svg?branch=master)

## ToC

* [ansi](https://github.com/SimonWaldherr/golibs#ansi-----) can print colored and styled text to your terminal
* [arg](https://github.com/SimonWaldherr/golibs#arg-----) simplifies cli flags (arguments)
* [as](https://github.com/SimonWaldherr/golibs#as-----) can convert most standard data types to most other data types
* [cache](https://github.com/SimonWaldherr/golibs#cache-----) is an easy and small caching package
* [cachedfile](https://github.com/SimonWaldherr/golibs#cachedfile-----) simplifies reading and writing from and to disk and adds caching
* [channel](https://github.com/SimonWaldherr/golibs#channel-----) simplifies channel operations, e.g. sending the same data to multiple receivers
* [csv](https://github.com/SimonWaldherr/golibs#csv-----) load and parse CSV file
* [file](https://github.com/SimonWaldherr/golibs#file-----) wraps around the standard functions to simplify reading and writing on disk
* [foreach](https://github.com/SimonWaldherr/golibs#foreach-----) calls a given function for each element of a [ JSON-string ]
* [gcurses](https://github.com/SimonWaldherr/golibs#gcurses-----) enabling the development of text user interface applications
* [gopath](https://github.com/SimonWaldherr/golibs#gopath-----) provides an easy way to get system information
* [graphics](https://github.com/SimonWaldherr/golibs#graphics-----) can manipulate images
* [log](https://github.com/SimonWaldherr/golibs#log-----) helps on logging your applications status
* [re](https://github.com/SimonWaldherr/golibs#re-----) helps you whenever you have to do something multiple times
* [regex](https://github.com/SimonWaldherr/golibs#regex-----) is a layer to speed up your regular expression development
* [rss](https://github.com/SimonWaldherr/golibs#rss-----) is a rss feed parser based on Golangs std xml package
* [ssl](https://github.com/SimonWaldherr/golibs#ssl-----) generates ssl certificates for https
* [stack](https://github.com/SimonWaldherr/golibs#stack-----) can store your values in stacks and rings
* [structs](https://github.com/SimonWaldherr/golibs#structs-----) use structs like maps
* [xmath](https://github.com/SimonWaldherr/golibs#xmath-----) provides a few mathematical functions like Sum, Median, Harmonic-mean, …
* [xtime](https://github.com/SimonWaldherr/golibs#xtime-----) xtime implements a subset of strftime

## WARNING

some functions in this repo enables you to write really bad code - I wrote it anyway because:  

* sometimes you need a quick and dirty solution
* it shows you possible ways, feel free to fork and customize
* mostly bad go code is better than good php code - every aspect in golang is designed to prevent from bad code

BUT: please remember that you can do things more performantly e.g. converting a string to a float is much faster done with strconv.ParseFloat than with as.Float.  

## sublibs

### ansi - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/ansi) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/ansi"
```

**ansi** can print colored and styled text to your terminal:

* green, yellow and red strings:  

```go
log.Println(ansi.Color("INFO: everything is fine", ansi.Green))
log.Println(ansi.Color("WARNING: not everything is fine", ansi.Yellow))
log.Println(ansi.Color("ERROR: OMG!!!", ansi.Red))
```

* bold and underlined text:  

```go
fmt.Printf("this is %v and %v text", ansi.Bold("bold"), ansi.Underline("underlined"))
```

### arg - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/arg) [![Coverage Status](https://img.shields.io/badge/coverage-89%25-green.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/arg"
```

the **arg** package simplifies cli flags (arguments)

### as - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/as) [![Coverage Status](https://img.shields.io/badge/coverage-99%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/as"
```

with **as** you can convert most standard data types to most other data types e.g.

* int to string:  

```go
var x string = as.String(int(32))
```

* string to int:  

```go
var x int = as.Int("32")
```

* string to time:  

```go
var x time.Time = as.Time("31.12.2014")
```

### cache - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/cache) [![Coverage Status](https://img.shields.io/badge/coverage-96%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/cache"
```

### cachedfile - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/cachedfile) [![Coverage Status](https://img.shields.io/badge/coverage-96%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/cachedfile"
```

**cachedfile** simplifies reading and writing from and to disk and adds caching

do

```go
str := "Neque porro quisquam est, qui dolorem ipsum, quia dolor sit, amet, consectetur, adipisci velit."
err := cachedfile.Write("filename.txt", str, false)
```

and in less then 15 minutes this

```go
str, _ := cachedfile.Read("filename.txt")
```

and there will be no file access to disk.
If you kill the App, wait 15 min or call ```cachedfile.Stop()``` the cached content will be exported to disk.

### channel - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/channel) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/channel"
```

**channel** simplifies channel operations, e.g. sending the same data to multiple receivers

```go
con := channel.Init()

receiver01 := con.AddReceiver()
receiver02 := con.AddReceiver()
receiver03 := con.AddReceiver()

go func() {
  fmt.Println(<-receiver01)
}()

go func() {
  fmt.Println(<-receiver02)
}()

go func() {
  fmt.Println(<-receiver03)
}()

transmitter01 := con.AddTransmitter()

transmitter01 <- "Hello World"
```

### csv - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/csv) [![Coverage Status](https://img.shields.io/badge/coverage-93%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/csv"
```

### file - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/file) [![Coverage Status](https://img.shields.io/badge/coverage-93%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/file"
```

**file** wraps around the standard functions to simplify reading and writing on disk

```go
str := "Neque porro quisquam est, qui dolorem ipsum, quia dolor sit, amet, consectetur, adipisci velit."
err := file.Write("filename.txt", str, false)
```

Besides simple reading and writing, the package also contains functions to test file statuses, read large files by small blocks, clear, rename and delete files.
There is even a function for do things with each file in a directory (even with subfolders):

```go
err := file.Each("..", true, func(filename, extension, filepath string, dir bool, fileinfo os.FileInfo) {
	if extension == "go" && !dir {
		t.Logf("%v, %v, %v, %v\n", filename, filepath, dir, fileinfo)
	}
}
```

If you need the absolute path to a file, but only have a relative path, you can use ```file.GetAbsolutePath("~/path/to/file.txt")```.  

### foreach - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/foreach) [![Coverage Status](https://img.shields.io/badge/coverage-90%25-green.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/foreach"
```

**foreach** calls a given function for each element of a [ JSON-string ]

### gcurses - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/gcurses) [![Coverage Status](https://img.shields.io/badge/coverage-90%25-green.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/gcurses"
```

**gcurses** is a terminal control library for Unixoide systems, enabling the development of text user interface applications.  

Named after the "GUI-like" terminal application toolkit [ncurses](https://en.wikipedia.org/wiki/Ncurses) (new curses) which is named after the original [curses](https://en.wikipedia.org/wiki/Curses_(programming_library)). 
This library is in a very early stage.  

```go
package main

import (
  "fmt"
  "simonwaldherr.de/go/golibs/gcurses"
  "time"
)

func main() {
  writer := gcurses.New()

  writer.Start()

  for i := 0; i < 100; i++ {
    fmt.Fprintf(writer, "Count till one hundred: %d\nStill counting ...\n", i)
    time.Sleep(time.Millisecond * 10)
  }

  time.Sleep(time.Millisecond * 500)
  fmt.Fprintln(writer, "Finished counting")

  writer.Stop()
}
```

### gopath - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/gopath) [![Coverage Status](https://img.shields.io/badge/coverage-93%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/gopath"
```

**gopath** provides an easy way to get system information

to read a config file in which is in the same file as the executable, you can do something like this:

```go
package main

import (
  "fmt"
  "path/filepath"
  "simonwaldherr.de/go/golibs/file"
  "simonwaldherr.de/go/golibs/gopath"
)

func main() {
  dir := gopath.Dir()
  config := file.Read(filepath.Join(dir, "config.yaml"))
  
  ...
```

### graphics - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/graphics) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/graphics"
```

with **graphics** you can manipulate images  

```go
img := graphics.EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
	return g, b, r, a
})
```

### log - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/log) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/log"
```


### re - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/re) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/re"
```

**re** helps you whenever you have to do something multiple times  

```go
data, stop := re.Do(time.Second * 5, func(data chan<- interface{}) {
	data <- fmt.Sprintf("%v\n", time.Now().Format("02.01.2006 15:04:05"))
})
```


### regex - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/regex) [![Coverage Status](https://img.shields.io/badge/coverage-88%25-green.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/regex"
```

**regex** is a layer to speed up your regular expression development  

```go
str, err := regex.ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1")
```


### rss - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/rss) [![Coverage Status](https://img.shields.io/badge/coverage-92%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/rss"
```

**rss** is a rss feed parser based on Golangs std xml package  

```go
podcast, err := rss.Read(url)
if err == nil {
  for _, episode := range podcast.Items {
    fmt.Println(episode.Title)
  }
}
```


### ssl - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/ssl) [![Coverage Status](https://img.shields.io/badge/coverage-87%25-green.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/ssl"
```

**ssl** generates ssl certificates for http**s**  

```go
err := ssl.Generate(options)
```


### stack - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/stack) [![Coverage Status](https://img.shields.io/badge/coverage-99%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/stack"
```

with **stack** you can store your values in stacks and rings  

```go
array := stack.Lifo()
array.Push(as.Bytes(12.34))
array.Push(as.Float(13.37))
array.Push(as.String(23.0))
for array.Len() > 0 {
	log.Println(array.Pop())
}
```

### structs - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/structs) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/structs"
```

### xmath - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/xmath) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/xmath"
```

**xmath** provides a few mathematical functions like *Max*, *Min*, *Sum*, *Median*, *Harmonic*-mean, ...

```go
var f = []float64{.5, 1.33, 2.66, 3.99, 13.37, 23.42, 42.000003}

fmt.Printf("Max: %v\n", xmath.Max(f))
fmt.Printf("Min: %v\n", xmath.Min(f))
fmt.Printf("Sum: %v\n", xmath.Sum(f))

fmt.Printf("Median:     %v\n", xmath.Median(f))
fmt.Printf("Arithmetic: %v\n", xmath.Arithmetic(f))
fmt.Printf("Harmonic:   %v\n", xmath.Harmonic(f))
fmt.Printf("Geometric:  %v\n", xmath.Geometric(f))
```

### xtime - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/xtime) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/xtime"
```

**xtime** implements a subset of [strftime](http://man7.org/linux/man-pages/man3/strftime.3.html)

```go
t, _ := time.Now()
fmt.Println(xtime.Fmt("%Y-%m-%d %H:%M:%S", t))
```
