package pastebin

import (
	"net/http"
	"time"
)

var client *http.Client

func getHttpClient() *http.Client {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return client
}
