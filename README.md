#golibs

golang functions (to be included in other projects)

```sh
go get -u -t simonwaldherr.de/go/golibs/...
```

##coverage & tests

```sh
go test ./...
```

 . | service | info
---|---------|------
[![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs) | coveralls.io | test coverage  
[![Build Status](https://img.shields.io/travis/SimonWaldherr/golibs.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs) | travis-ci.org | test on various go versions  
[![Build status](https://img.shields.io/appveyor/ci/SimonWaldherr/golibs.svg?style=flat-square)](https://ci.appveyor.com/project/SimonWaldherr/golibs/branch/master) | appveyor.com | test under windows  
[![Build status](https://magnum-ci.com/status/e9ccc5689f4135e4021475bfdf0cf527.png)](https://magnum-ci.com/public/c4dba43a1c41dbff557e/builds) | magnum-ci.com | yet another ci service  
[![License MIT](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT) |  | free + open source license  
[![Flattr donate button](https://img.shields.io/badge/flattr-donate-orange.svg?style=flat-square)](https://flattr.com/submit/auto?user_id=SimonWaldherr&url=http%3A%2F%2Fgithub.com%2FSimonWaldherr%2Fgolibs "Donate monthly to this project using Flattr") | flattr.com | micro donation  
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/) | godoc.org | documentation  


each new build gets tested in multiple steps:

* on development i regularly type ```go test``` from time to time to check the test suite
* also there are a few go apps in the [examples](https://github.com/SimonWaldherr/golibs/tree/master/examples)-folder which i test to build and run
* on commit, git automatically runs the [pre-commit](https://github.com/SimonWaldherr/golibs/blob/master/pre-commit)-hook shell script
* after a commit gets pushed to **GitHub**, the following tests are started via Webhooks and Services
	* **Travis CI** build the lib and all tests on docker containers with the go versions noted in [.travis.yml](https://github.com/SimonWaldherr/golibs/blob/master/.travis.yml)
	* **appveyor** builds the lib on Windows Server to test against the Microsoft Infrastructure ([.appveyor.yml](https://github.com/SimonWaldherr/golibs/blob/master/.appveyor.yml))
	* **magnum-ci**, another ci service build the lib and tests on a linux machine

##sublibs

###ansi

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/ansi)  

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

###as

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/as)  

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

###cache

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/cache)  

```go
import "simonwaldherr.de/go/golibs/cache"
```


###file

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/file)  

```go
import "simonwaldherr.de/go/golibs/file"
```

**file** wraps around the standard functions to simplify reading and writing on disk

```go
str := "Neque porro quisquam est, qui dolorem ipsum, quia dolor sit, amet, consectetur, adipisci velit."
err := file.Write("filename.txt", str, false)
```

###graphics

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/graphics)  

```go
import "simonwaldherr.de/go/golibs/graphics"
```

with **graphics** you can manipulate images  

```go
img := graphics.EachPixel(file, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
	return g, b, r, a
})
```


###re

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/re)  

```go
import "simonwaldherr.de/go/golibs/re"
```

**re** helps you whenever you have to do something multiple times  

```go
data, stop := re.Do(time.Second * 5, func(data chan<- interface{}) {
	data <- fmt.Sprintf("%v\n", time.Now().Format("02.01.2006 15:04:05"))
})
```


###regex

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/regex)  

```go
import "simonwaldherr.de/go/golibs/regex"
```

**regex** is a layer to speed up your regular expression development  

```go
str := regex.ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1")
```


###ssl

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/ssl)  

```go
import "simonwaldherr.de/go/golibs/ssl"
```

**ssl** generates ssl certificates for http**s**  

```go
err := ssl.Generate(options)
```


###stack

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/stack)  

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

###xmath

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/xmath)  

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
