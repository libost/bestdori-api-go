package bestdoriapigo

import (
	"context"
	"net/http"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(nil),
		},
	})
	if client == nil {
		t.Error("Expected non-nil client")
	}
}

// Before running this test, make sure to set the environment variables BESTDORI_API_USER and BESTDORI_API_PSK with valid credentials.
func TestNewClientWithLogin(t *testing.T) {
	if os.Getenv("BESTDORI_API_USER") == "" || os.Getenv("BESTDORI_API_PSK") == "" {
		t.Skip("Environment variables BESTDORI_API_USER and BESTDORI_API_PSK are not set")
	}
	client, err := NewClientWithLogin(context.TODO(), nil, os.Getenv("BESTDORI_API_USER"), os.Getenv("BESTDORI_API_PSK"))
	if err != nil {
		t.Fatalf("NewClientWithLogin returned an error: %v", err)
	}
	if client == nil {
		t.Error("Expected non-nil client")
	}
}
