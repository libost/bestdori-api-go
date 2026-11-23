package bestdoriapigo

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("", 10*time.Second)
	if client == nil {
		t.Error("Expected non-nil client")
	}
}
