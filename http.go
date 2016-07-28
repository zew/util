package util

import (
	"net"
	"net/http"
	"time"
)

// Get a http client
func HttpClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 1110,
		Transport: netTransport,
	}
	return netClient
}
