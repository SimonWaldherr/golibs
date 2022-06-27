package http

import (
	//"context"
	"io/ioutil"
	//"net"
	"net/http"
	//"strings"
	"time"
	//"github.com/rs/dnscache"
)

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
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
