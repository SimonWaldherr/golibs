package specialdays_test

import (
	"fmt"
	"sort"
	"simonwaldherr.de/go/golibs/specialdays"
)

func ExampleGetSpecialDays() {
	days := specialdays.GetSpecialDays(2019)
	
	keys := make([]string, 0, len(days))
	
		for k := range days{
			keys = append(keys, k)
		}
	sort.Strings(keys)
	
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, days[k].Format("2006-01-02"))
	}

	// Output: 
	// 1. Weihnachtstag: 2019-12-25
	// 2. Weihnachtstag: 2019-12-26
	// Allerheiligen: 2019-11-01
	// Buß- und Bettag: 2019-11-20
	// Christi Himmelfahrt: 2019-05-30
	// Fronleichnam: 2019-06-20
	// Halloween: 2019-10-31
	// Heiligabend: 2019-12-24
	// Heilige Drei Könige: 2019-01-06
	// Karfreitag: 2019-04-19
	// Karnevalsbeginn: 2019-11-11
	// Mariä Himmelfahrt: 2019-08-15
	// Muttertag: 2019-05-12
	// Neujahrstag: 2019-01-01
	// Nikolaustag: 2019-12-06
	// Ostermontag: 2019-04-22
	// Ostersonntag: 2019-04-21
	// Pfingstmontag: 2019-06-10
	// Pfingstsonntag: 2019-06-09
	// Silvester: 2019-12-31
	// Tag der Arbeit: 2019-05-01
	// Tag der Deutschen Einheit: 2019-10-03
	// Valentinstag: 2019-02-14
	// erster Advent: 2019-12-01
}
