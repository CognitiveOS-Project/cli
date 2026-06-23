package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Envelope struct {
	Type      string          `json:"type"`
	ID        string          `json:"id,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	From      string          `json:"from,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

type Conn struct {
	conn     *net.UnixConn
	reader   *bufio.Scanner
	encoder  *json.Encoder
	mu       sync.Mutex
	done     chan struct{}
	Messages chan Envelope
	Address  string
}

func New(address string) *Conn {
	if address == "" {
		address = "/cognitiveos/run/daemon.sock"
	}
	return &Conn{
		Address:  address,
		done:     make(chan struct{}),
		Messages: make(chan Envelope, 64),
	}
}

func (c *Conn) Connect() error {
	addr, err := net.ResolveUnixAddr("unix", c.Address)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewScanner(conn)
	c.reader.Buffer(make([]byte, 65536), 1048576)
	c.encoder = json.NewEncoder(conn)

	go c.readLoop()
	return nil
}

func (c *Conn) Close() {
	select {
	case <-c.done:
		return
	default:
		close(c.done)
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Conn) Send(env Envelope) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return fmt.Errorf("connection closed")
	default:
	}
	if c.encoder == nil {
		return fmt.Errorf("not connected")
	}
	return c.encoder.Encode(env)
}

func (c *Conn) SendMessage(msgType string, payload interface{}) error {
	env := Envelope{
		Type:      msgType,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		From:      "cli",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	env.Payload = b
	return c.Send(env)
}

func (c *Conn) SendInput(content string) error {
	return c.SendMessage("input_forward", map[string]interface{}{
		"mode":    "text",
		"content": content,
		"context": map[string]interface{}{
			"session_id": fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		},
	})
}

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

func (c *Conn) RequestStatus() error {
	return c.SendMessage("status_request", map[string]interface{}{})
}

func (c *Conn) RequestAudit() error {
	return c.SendMessage("audit_request", map[string]interface{}{})
}

func (c *Conn) IsClosed() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *Conn) readLoop() {
	defer c.Close()

	for {
		select {
		case <-c.done:
			return
		default:
		}

		if !c.reader.Scan() {
			return
		}

		var env Envelope
		if err := json.Unmarshal(c.reader.Bytes(), &env); err != nil {
			continue
		}

		select {
		case c.Messages <- env:
		default:
		}
	}
}
