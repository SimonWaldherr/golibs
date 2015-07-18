// the arg package simplifies cli flags (arguments)
package arg

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

type argument struct {
	name         string
	defaultValue string
	flagValue    *string
	usage        string
	timeout      time.Duration
}

var flagVar = make(map[string]argument)
var values = make(map[string]interface{})

func String(name, defaultValue, usage string, timeout time.Duration) {
	flagVar[name] = argument{
		name:         name,
		defaultValue: defaultValue,
		flagValue:    flag.String(name, "", usage),
		usage:        usage,
		timeout:      timeout,
	}
}

func Parse() {
	var value interface{}

	flag.Parse()

	for v := range flagVar {
		if *flagVar[v].flagValue == "" {
			fmt.Printf("# %v: ", flagVar[v].usage)
			scanner := bufio.NewScanner(os.Stdin)
			ch := make(chan bool, 1)
			go func() {
				defer func() {
					close(ch)
				}()
				ch <- scanner.Scan()
			}()
			select {
			case <-ch:
				value = scanner.Text()
			case <-time.After(flagVar[v].timeout):
				value = flagVar[v].defaultValue
			}
		} else {
			value = *flagVar[v].flagValue
		}
		values[v] = value
	}
}

func Get(name string) interface{} {
	return values[name]
}

func Dump() string {
	return fmt.Sprintf("%v\n", flagVar)
}
