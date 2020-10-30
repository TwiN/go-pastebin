package pastebin

import (
	"net/http"
	"time"
)

var client HttpClient

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// getHTTPClient returns the shared HTTP client
func getHTTPClient() HttpClient {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return client
}
