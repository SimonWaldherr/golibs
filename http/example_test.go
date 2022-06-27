package http_test

import (
	"fmt"
	"io/ioutil"
	ohttp "net/http"
	"net/http/httptest"
	"time"

	"simonwaldherr.de/go/golibs/http"
)

func ExampleClient() {
	ts := httptest.NewServer(ohttp.HandlerFunc(func(rw ohttp.ResponseWriter, req *ohttp.Request) {
		userAgent := req.UserAgent()
		if userAgent == "Golang_Bot/1.0" {
			fmt.Fprintf(rw, "hello gopher\n")
		} else {
			fmt.Fprintf(rw, "hello world\n")
		}
	}))
	defer ts.Close()

	client := http.Client(time.Second * 15)
	resp, err := client.Get(ts.URL)

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
	ts := httptest.NewServer(ohttp.HandlerFunc(func(rw ohttp.ResponseWriter, req *ohttp.Request) {
		userAgent := req.UserAgent()
		if userAgent == "Golang_Bot/1.0" {
			fmt.Fprintf(rw, "hello gopher\n")
		} else {
			fmt.Fprintf(rw, "hello world\n")
		}
	}))
	defer ts.Close()

	resp, err := http.GetString(ts.URL)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)

	// Output: hello world
}

func ExampleUserAgent() {
	ts := httptest.NewServer(ohttp.HandlerFunc(func(rw ohttp.ResponseWriter, req *ohttp.Request) {
		userAgent := req.UserAgent()
		if userAgent == "Golang_Bot/1.0" {
			fmt.Fprintf(rw, "hello gopher\n")
		} else {
			fmt.Fprintf(rw, "hello world\n")
		}
	}))
	defer ts.Close()

	client := http.Client(time.Second * 15)

	req, _ := http.NewRequest("GET", ts.URL, nil)
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
