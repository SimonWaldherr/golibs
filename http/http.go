package http

import (
	//"context"
	"io"
	//"net"
	"net/http"
	//"strings"
	"time"
	//"github.com/rs/dnscache"
)

// UserAgent is the default user agent string used by the HTTP client.
const UserAgent = "Golang_Bot/1.0"

var Transporter *http.Transport
var client http.Client
var clientReady bool = false

func init() {
	//dnsResolver := &dnscache.Resolver{}
	Transporter = &http.Transport{
		/*
			DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
				separator := strings.LastIndex(addr, ":")
				ips, err := dnsResolver.LookupHost(ctx, addr[:separator])
				if err != nil {
					return nil, err
				}
				for _, ip := range ips {
					conn, err = net.Dial(network, ip+addr[separator:])
					if err == nil {
						break
					}
				}
				return
			},*/
		MaxIdleConns:    1024,
		MaxConnsPerHost: 1024,
		IdleConnTimeout: 10 * time.Second,
	}
	/*
		go func() {
			cacheTicker := time.NewTicker(30 * time.Minute)
			defer cacheTicker.Stop()
			for range cacheTicker.C {
				dnsResolver.Refresh(true)
			}
		}()*/
}

var NewRequest = http.NewRequest

func Client(timeout time.Duration) http.Client {
	return http.Client{
		Transport: Transporter,
		Timeout:   timeout,
	}
}

func GetString(url string) (string, error) {
	if clientReady == false {
		client = Client(time.Second * 15)
	}

	resp, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetString2 fetches a URL with a custom timeout and returns the response body as a string.
// The numeric suffix distinguishes it from GetString which uses the shared default client.
func GetString2(url string, timeout time.Duration) (string, error) {
	c := Client(timeout)
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
