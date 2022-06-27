package http_test

import (
	"fmt"
	"io/ioutil"
	webserver "net/http"
	"time"

	"simonwaldherr.de/go/golibs/http"
)

func init() {
	webserver.HandleFunc("/hello", func(rw webserver.ResponseWriter, req *webserver.Request) {
		userAgent := req.UserAgent()
		if userAgent == "Golang_Bot/1.0" {
			fmt.Fprintf(rw, "hello gopher\n")
		} else {
			fmt.Fprintf(rw, "hello world\n")
		}
	})
	go webserver.ListenAndServe(":8081", nil)
}

func ExampleClient() {
	client := http.Client(time.Second * 15)
	resp, err := client.Get("http://localhost:8081/hello")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))

	// Output: hello world
}

func ExampleGetString() {
	resp, err := http.GetString("http://localhost:8081/hello")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)

	// Output: hello world
}

func ExampleUserAgent() {
	client := http.Client(time.Second * 15)

	req, _ := http.NewRequest("GET", "http://localhost:8081/hello", nil)
	req.Header.Set("User-Agent", "Golang_Bot/1.0")

	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	// Output: hello gopher
}

func ExampleGetString2() {
	resp, err := http.GetString("https://google.de/")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
