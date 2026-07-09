package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/CognitiveOS-Project/cli/internal/client"
	"github.com/CognitiveOS-Project/cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	socketPath := flag.String("socket", "/cognitiveos/run/daemon.sock", "daemon socket path")
	flag.Parse()

	conn := client.New(*socketPath)

	if _, err := os.Stat(*socketPath); os.IsNotExist(err) {
		cmd := exec.Command("cognitiveosd")
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("spawn cognitiveosd: %v", err)
		}
		defer func() {
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				log.Printf("failed to signal cognitiveosd: %v", err)
			}
			if err := cmd.Wait(); err != nil {
				log.Printf("cognitiveosd exited with error: %v", err)
			}
		}()
		for i := 0; i < 8; i++ {
			if _, err := os.Stat(*socketPath); err == nil {
				break
			}
			time.Sleep(250 * time.Millisecond)
		}
	}

	if err := conn.Connect(); err != nil {
		log.Fatalf("connect: %v", err)
	}

	p := tea.NewProgram(
		tui.NewModel(conn),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(fmt.Errorf("program error: %w", err))
	}

	conn.Close()
	os.Exit(0)
}
