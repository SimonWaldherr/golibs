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
 ![Linux Build Status](https://simonwaldherr.de/icon/tux.png) | [![Build Status](https://img.shields.io/travis/SimonWaldherr/golibs.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs) | travis-ci.org | test on various go versions  
 ![Windows Build Status](https://simonwaldherr.de/icon/win.png) | [![Build status](https://img.shields.io/appveyor/ci/SimonWaldherr/golibs.svg?style=flat-square)](https://ci.appveyor.com/project/SimonWaldherr/golibs/branch/master) | appveyor.com | test under windows  
 ![Linux Build Status](https://simonwaldherr.de/icon/tux.png) | [![Build status](https://magnum-ci.com/status/e9ccc5689f4135e4021475bfdf0cf527.png)](https://magnum-ci.com/public/c4dba43a1c41dbff557e/builds) | magnum-ci.com | just another ci service  
 ![Linux Build Status](https://simonwaldherr.de/icon/tux1.png) | [![Build Status](https://semaphoreci.com/api/v1/projects/fe1a7a53-a2c0-4bc0-a539-4af4ef13d49f/487313/shields_badge.svg)](https://semaphoreci.com/simonwaldherr/golibs) | semaphoreci.com | yet another ci service  
   | [![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) | coveralls.io | test coverage  
   | [![codecov.io](http://codecov.io/github/SimonWaldherr/golibs/coverage.svg?branch=master)](http://codecov.io/github/SimonWaldherr/golibs?branch=master) | codecov.io | test coverage  
   | [![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://simonwaldherr.de/gocover/golibs/) |  | go tool cover  
   | [![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT) |  | free + open source license  
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

### cache - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/cache) [![Coverage Status](https://img.shields.io/badge/coverage-95%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

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


### regex - [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/regex) [![Coverage Status](https://img.shields.io/badge/coverage-100%25-brightgreen.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)

```go
import "simonwaldherr.de/go/golibs/regex"
```

**regex** is a layer to speed up your regular expression development  

```go
str := regex.ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1")
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
