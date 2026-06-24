package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/CognitiveOS-Project/cli/internal/client"
	"github.com/CognitiveOS-Project/cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	socketPath := flag.String("socket", "/cognitiveos/run/daemon.sock", "daemon socket path")
	flag.Parse()

	conn := client.New(*socketPath)

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
