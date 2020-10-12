package pastebin

import (
	"net/http"
	"time"
)

var client HttpClient

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func getHttpClient() HttpClient {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return client
}
