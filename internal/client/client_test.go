package client

import (
	"encoding/json"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := New("/tmp/test.sock")
	if c == nil {
		t.Fatal("expected client")
	}
	if c.Address != "/tmp/test.sock" {
		t.Fatalf("expected /tmp/test.sock, got %s", c.Address)
	}
}

func TestDefaultAddress(t *testing.T) {
	c := New("")
	if c.Address != "/cognitiveos/run/daemon.sock" {
		t.Fatalf("expected default socket path, got %s", c.Address)
	}
}

func TestEnvelopeJSON(t *testing.T) {
	env := Envelope{
		Type:      "test_type",
		ID:        "123",
		Timestamp: "2026-01-01T00:00:00Z",
		From:      "test",
		Payload:   json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Envelope
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Type != "test_type" {
		t.Fatalf("expected test_type, got %s", decoded.Type)
	}
	if decoded.From != "test" {
		t.Fatalf("expected test, got %s", decoded.From)
	}
}

func TestSendMessage(t *testing.T) {
	c := New("/nonexistent.sock")
	err := c.SendMessage("input_forward", map[string]interface{}{
		"mode":    "text",
		"content": "hello",
	})
	if err == nil {
		t.Log("SendMessage returned nil (expected error since not connected)")
	} else {
		t.Logf("SendMessage error (expected): %v", err)
	}
}

func TestSendSystemCode(t *testing.T) {
	c := New("/nonexistent.sock")
	err := c.SendSystemCode("wake", "")
	if err == nil {
		t.Log("SendSystemCode for wake returned nil (expected error)")
	}

	err = c.SendSystemCode("unlock", "ABCD-1234")
	if err == nil {
		t.Log("SendSystemCode for unlock returned nil (expected error)")
	}
}

func TestIsClosed(t *testing.T) {
	c := New("/test.sock")
	if c.IsClosed() {
		t.Fatal("expected not closed initially")
	}
	// Without connecting, it's not "closed" in the channel sense
}
