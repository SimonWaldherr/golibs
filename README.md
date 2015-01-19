#golibs

[![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg)](https://coveralls.io/r/SimonWaldherr/golibs)
[![Build Status](https://travis-ci.org/SimonWaldherr/golibs.svg)](https://travis-ci.org/SimonWaldherr/golibs)

golang functions (to be included in other projects)

```sh
go get github.com/simonwaldherr/golibs/...
```

##ansi

[![GoDoc](https://godoc.org/github.com/SimonWaldherr/golibs/ansi?status.svg)](https://godoc.org/github.com/SimonWaldherr/golibs/ansi)  

```
go get github.com/simonwaldherr/golibs/ansi
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

[![GoDoc](https://godoc.org/github.com/SimonWaldherr/golibs/as?status.svg)](https://godoc.org/github.com/SimonWaldherr/golibs/as)  

```
go get github.com/simonwaldherr/golibs/as
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

[![GoDoc](https://godoc.org/github.com/SimonWaldherr/golibs/cache?status.svg)](https://godoc.org/github.com/SimonWaldherr/golibs/cache)  

```
go get github.com/simonwaldherr/golibs/cache
```

##stack

[![GoDoc](https://godoc.org/github.com/SimonWaldherr/golibs/stack?status.svg)](https://godoc.org/github.com/SimonWaldherr/golibs/stack)  

```
go get github.com/simonwaldherr/golibs/stack
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

[![GoDoc](https://godoc.org/github.com/SimonWaldherr/golibs/xmath?status.svg)](https://godoc.org/github.com/SimonWaldherr/golibs/xmath)  

```
go get github.com/simonwaldherr/golibs/xmath
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
