# CognitiveOS CLI (TUI)

The human interface — a Bubble Tea TUI that replaces the traditional desktop/app paradigm.

## Features

- Clean text prompt: "Listening..." as default state
- Voice input capture (via ALSA / audio-mcp)
- Text input for keyboard interaction
- Direct framebuffer overlay for images/video (communicates with display-mcp)
- Connects to cognitiveosd via Unix socket

## Build

```bash
go build -o bin/cognitiveos-cli ./cmd/cognitiveos-cli
```

## Architecture

The TUI is thin — it captures input and displays output. All intelligence lives in cognitiveosd and the Wide Model. The TUI can crash and restart without affecting the OS.

## Dependencies

- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`
- CognitiveOS internal: `cognitiveosd` daemon socket
