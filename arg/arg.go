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
	argNr        int
}

var flagVar = make(map[string]argument)
var values = make(map[string]interface{})
var i = 1

func String(name, defaultValue, usage string, timeout time.Duration) {
	flagVar[name] = argument{
		name:         name,
		defaultValue: defaultValue,
		flagValue:    flag.String(name, "", usage),
		usage:        usage,
		timeout:      timeout,
		argNr:        -1,
	}
}

func StringArg(name, defaultValue, usage string, timeout time.Duration) {
	flagVar[name] = argument{
		name:         name,
		defaultValue: defaultValue,
		flagValue:    flag.String(name, "", usage),
		usage:        usage,
		timeout:      timeout,
		argNr:        i,
	}
	i++
}

func Parse() {
	var value interface{}

	flag.Parse()

	for v := range flagVar {
		if *flagVar[v].flagValue != "" {
			value = *flagVar[v].flagValue
		} else if !(flagVar[v].argNr >= 1 && flagVar[v].argNr < len(os.Args)) && flagVar[v].timeout >= time.Second*1 {
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
		} else if (flagVar[v].argNr >= 1 && flagVar[v].argNr < len(os.Args)) {
			value = os.Args[flagVar[v].argNr]
		} else {
			value = flagVar[v].defaultValue
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
