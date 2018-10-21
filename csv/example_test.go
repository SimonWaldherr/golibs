package csv_test

import (
    csv "."
    "fmt"
)

var userdata string = `id;name;email
0;John Doe;jDoe@example.org
1;Jane Doe;jane.doe@example.com
2;Max Mustermann;m.mustermann@alpha.tld`

func Example() {
    csvmap, k := csv.LoadCSVfromString(userdata)
    for _, user := range csvmap {
        fmt.Println(user[k["name"]])
    }
    
    // Output: John Doe
    // Jane Doe
    // Max Mustermann
}

func Example_second() {
    csvmap, k := csv.LoadCSVfromFile("./test.csv")
    for _, user := range csvmap {
        fmt.Println(user[k["name"]])
    }

    // Output: John Doe
    // Jane Doe
    // Max Mustermann
}
