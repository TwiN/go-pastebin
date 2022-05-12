package pastebin

import (
	"net/http"
	"testing"
	"time"
)

func TestValidateClient(t *testing.T) {
	httpClient = nil
	client := getHTTPClient()
	if client == nil {
		t.Error("client should not be nil")
	}
	if client.Timeout != 10*time.Second {
		t.Error("default timeout should be 10 seconds")
	}
	tr := client.Transport.(*http.Transport)
	if tr.MaxIdleConnsPerHost != 50 {
		t.Error("default max conns idle conns per host should be 150")
	}
	if tr.MaxIdleConns != 100 {
		t.Error("default max conns idle conns should be 100")
	}
}
