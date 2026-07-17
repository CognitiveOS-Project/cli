package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CognitiveOS-Project/cli/internal/client"
)

const defaultTimeout = 30 * time.Second

type outputPayload struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type,omitempty"`
}

// Run connects to the daemon, sends a command, waits for the response, and exits.
// If jsonOutput is true, the full JSON envelope is printed. Otherwise, only the text content.
func Run(conn *client.Conn, command string, jsonOutput bool) error {
	if err := conn.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	if err := conn.SendInput(command); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	deadline := time.After(defaultTimeout)

	for {
		select {
		case env, ok := <-conn.Messages:
			if !ok {
				return fmt.Errorf("connection closed before response")
			}
			if env.Type == "output_deliver" {
				if jsonOutput {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "  ")
					return enc.Encode(env)
				}
				return printContent(env)
			}
		case <-deadline:
			return fmt.Errorf("timeout waiting for response (30s)")
		}
	}
}

func printContent(env client.Envelope) error {
	var p outputPayload
	if err := json.Unmarshal(env.Payload, &p); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	fmt.Println(p.Content)
	return nil
}
