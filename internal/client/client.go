package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/CognitiveOS-Project/cogsdk/client/daemon"
)

// Envelope is a JSON-RPC style message envelope sent over the daemon socket.
type Envelope struct {
	Type      string          `json:"type"`
	ID        string          `json:"id,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	From      string          `json:"from,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

// Conn is a client connection to the cognitiveosd daemon.
// It wraps cogsdk/client/daemon.Client to provide the same API that the
// rest of the cli codebase expects.
type Conn struct {
	inner   *daemon.Client
	addr    string
	done    chan struct{}
	closeOnce sync.Once
	Messages chan Envelope
}

// New creates a new connection to the cognitiveosd daemon socket.
// If address is empty, the default socket path is used.
func New(address string) *Conn {
	if address == "" {
		address = daemon.DefaultSocketPath()
	}
	return &Conn{
		addr:    address,
		done:    make(chan struct{}),
		Messages: make(chan Envelope, 64),
	}
}

// Address returns the socket address this connection uses.
func (c *Conn) Address() string {
	return c.addr
}

// Connect establishes a connection to the daemon socket.
func (c *Conn) Connect() error {
	cli, err := daemon.DialPath(c.addr)
	if err != nil {
		return fmt.Errorf("daemon connection failed: %w", err)
	}
	cli.SetReadTimeout(5 * time.Second)
	c.inner = cli
	go c.readLoop()
	return nil
}

// Close closes the connection.
func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		close(c.Messages)
	})
	if c.inner != nil {
		_ = c.inner.Close()
	}
}

// Send sends an envelope over the connection.
func (c *Conn) Send(env Envelope) error {
	if c.inner == nil {
		return fmt.Errorf("not connected")
	}
	return c.inner.SendEnvelope(daemon.Envelope(env))
}

// SendMessage sends a typed message over the connection.
func (c *Conn) SendMessage(msgType string, payload interface{}) error {
	if c.inner == nil {
		return fmt.Errorf("not connected")
	}
	return c.inner.SendMessage(msgType, payload)
}

// SendInput sends an input_forward message to the daemon.
func (c *Conn) SendInput(content string) error {
	return c.SendMessage("input_forward", map[string]interface{}{
		"mode":    "text",
		"content": content,
		"context": map[string]interface{}{
			"session_id": fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		},
	})
}

// SendVoice sends a voice command to the daemon.
func (c *Conn) SendVoice() error {
	return c.SendMessage("input_forward", map[string]interface{}{
		"mode":    "voice",
		"content": "",
		"context": map[string]interface{}{
			"session_id": fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		},
	})
}

// SendSystemCode sends a system code to the daemon.
func (c *Conn) SendSystemCode(code string, unlockCode string) error {
	payload := map[string]interface{}{
		"code":   code,
		"origin": "keyboard",
	}
	if unlockCode != "" {
		payload["unlock_code"] = unlockCode
	}
	return c.SendMessage("system_code", payload)
}

// RequestStatus sends a status request to the daemon.
func (c *Conn) RequestStatus() error {
	return c.SendMessage("status_request", map[string]interface{}{})
}

// RequestAudit sends an audit request to the daemon.
func (c *Conn) RequestAudit() error {
	return c.SendMessage("audit_request", map[string]interface{}{})
}

// IsClosed returns true if the connection has been closed.
func (c *Conn) IsClosed() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

// readLoop reads messages from the daemon and forwards them to the Messages channel.
func (c *Conn) readLoop() {
	defer func() {
		c.closeOnce.Do(func() {
			close(c.done)
			close(c.Messages)
		})
		if c.inner != nil {
			_ = c.inner.Close()
		}
	}()

	for {
		select {
		case <-c.done:
			return
		default:
		}

		resp, err := c.inner.ReadResponse("")
		if err != nil {
			if c.IsClosed() {
				return
			}
			continue
		}
		select {
		case c.Messages <- Envelope(resp):
		case <-c.done:
			return
		}
	}
}