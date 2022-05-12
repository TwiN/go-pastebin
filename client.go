package pastebin

import (
	"net/http"
	"time"
)

var httpClient *http.Client

// getHTTPClient returns the shared HTTP client
func getHTTPClient() *http.Client {
	if httpClient == nil {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.MaxIdleConnsPerHost = 50
		httpClient = &http.Client{
			Timeout:   time.Second * 10,
			Transport: tr,
		}
	}
	return httpClient
}
