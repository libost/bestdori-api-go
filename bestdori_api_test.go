package bestdoriapigo

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(nil), // You can set a proxy URL here if needed
		},
	})
	if client == nil {
		t.Error("Expected non-nil client")
	}
}
