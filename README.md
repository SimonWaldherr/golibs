#golibs

[![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg?style=flat-square)](https://coveralls.io/r/SimonWaldherr/golibs)
[![Build Status](https://img.shields.io/travis/SimonWaldherr/golibs.svg?style=flat-square)](https://travis-ci.org/SimonWaldherr/golibs)
[![License MIT](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT)
[![Flattr donate button](https://raw.github.com/balupton/flattr-buttons/master/badge-89x18.gif)](https://flattr.com/submit/auto?user_id=SimonWaldherr&url=http%3A%2F%2Fgithub.com%2FSimonWaldherr%2Fgolibs "Donate monthly to this project using Flattr")

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/)  

golang functions (to be included in other projects)

```sh
go get -u -t simonwaldherr.de/go/golibs/...
```

##ansi

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

##as

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

##cache

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/cache)  

```go
import "simonwaldherr.de/go/golibs/cache"
```


##file

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/file)  

```go
import "simonwaldherr.de/go/golibs/file"
```

**file** wraps around the standard functions to simplify reading and writing on disk

```go
str := "Neque porro quisquam est, qui dolorem ipsum, quia dolor sit, amet, consectetur, adipisci velit."
err := file.Write("filename.txt", str, false)
```

##graphics

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


##regex

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SimonWaldherr/golibs/regex)  

```go
import "simonwaldherr.de/go/golibs/regex"
```

**regex** is a layer to speed up your regular expression development  

```go
str := regex.ReplaceAllString("Ipsum Lorem", "([^ ]+) ([^ ]+)", "$2 $1")
```

##stack

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

##xmath

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
