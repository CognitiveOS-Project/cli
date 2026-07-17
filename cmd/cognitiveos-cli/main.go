package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/CognitiveOS-Project/cli/internal/client"
	"github.com/CognitiveOS-Project/cli/internal/cli"
	"github.com/CognitiveOS-Project/cli/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	fs := flag.NewFlagSet("cognitiveos-cli", flag.ExitOnError)
	socketPath := fs.String("socket", "/cognitiveos/run/daemon.sock", "daemon socket path")
	tuiMode := fs.Bool("tui", false, "launch interactive TUI (default)")
	cmdText := fs.String("cmd", "", "send a command and print the response (non-interactive)")
	jsonOutput := fs.Bool("json", false, "print full JSON envelope (requires --cmd)")
	showVersion := fs.Bool("version", false, "print version and exit")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "cognitiveos-cli — CognitiveOS terminal interface\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli                  Launch interactive TUI (default)\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli --tui            Launch interactive TUI (explicit)\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli --cmd <text>     Send command, print response, exit\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli --cmd <text> --json  Print full JSON response\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli --version        Print version\n")
		fmt.Fprintf(os.Stderr, "  cognitiveos-cli --help           Print this help\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if *showVersion {
		fmt.Printf("cognitiveos-cli %s\n", version)
		os.Exit(0)
	}

	if *jsonOutput && *cmdText == "" {
		fmt.Fprintln(os.Stderr, "error: --json requires --cmd")
		os.Exit(1)
	}

	conn := client.New(*socketPath)

	if *cmdText != "" {
		if err := cli.Run(conn, *cmdText, *jsonOutput); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

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

	_ = *tuiMode

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
