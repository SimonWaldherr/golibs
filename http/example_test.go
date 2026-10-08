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
			_, _ = fmt.Fprintf(rw, "hello gopher\n")
		} else {
			_, _ = fmt.Fprintf(rw, "hello world\n")
		}
	}))
	defer ts.Close()

	client := http.Client(time.Second * 15)
	resp, err := client.Get(ts.URL)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

	// Output: hello world
}

func ExampleGetString() {
	ts := httptest.NewServer(ohttp.HandlerFunc(func(rw ohttp.ResponseWriter, req *ohttp.Request) {
		userAgent := req.UserAgent()
		if userAgent == "Golang_Bot/1.0" {
			_, _ = fmt.Fprintf(rw, "hello gopher\n")
		} else {
			_, _ = fmt.Fprintf(rw, "hello world\n")
		}
	}))
	defer ts.Close()

	resp, err := http.GetString(ts.URL)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)

	// Output: hello world
}
