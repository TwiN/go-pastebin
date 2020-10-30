package pastebin

import (
	"net/http"
	"testing"
	"time"
)

func TestValidateClient(t *testing.T) {
	oldRef := getHTTPClient()
	defer func() {
		client = oldRef
	}()
	client = nil
	client := getHTTPClient()
	if client == nil {
		t.Error("client should not be nil")
	}
	cl := client.(*http.Client)
	if cl.Timeout != 10*time.Second {
		t.Error("default timeout should be 10 seconds")
	}
	tr := cl.Transport.(*http.Transport)
	if tr.MaxIdleConnsPerHost != 150 {
		t.Error("default max conns idle conns per host should be 150")
	}
	if tr.MaxIdleConns != 150 {
		t.Error("default max conns idle conns should be 150")
	}
}
