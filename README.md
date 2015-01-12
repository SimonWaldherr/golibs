#golibs

[![Coverage Status](https://img.shields.io/coveralls/SimonWaldherr/golibs.svg)](https://coveralls.io/r/SimonWaldherr/golibs)
[![Build Status](https://travis-ci.org/SimonWaldherr/golibs.svg)](https://travis-ci.org/SimonWaldherr/golibs)

golang functions (to be included in other projects)

##as

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

##cli

```
go get github.com/simonwaldherr/golibs/cli
```

**cli** can print colored and styled text to your terminal:

* green, yellow and red strings:  

```go
log.Println(cli.Color("INFO: everything is fine", cli.Green))
log.Println(cli.Color("WARNING: not everything is fine", cli.Yellow))
log.Println(cli.Color("ERROR: OMG!!!", cli.Red))
```

* bold and underlined text:  

```go
fmt.Printf("this is %v and %v text", cli.Bold("bold"), cli.Underline("underlined"))
```

##stack

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

```
go get github.com/simonwaldherr/golibs/xmath
```

**xmath** provides a few mathematical functions like *Max*, *Min*, *Sum*, *Median*, *Harmonic*-mean, ...
