# CognitiveOS CLI

The human interface to CognitiveOS. Two modes: interactive TUI for live sessions, non-interactive CLI for scripting.

## Modes

- **TUI** (default): Bubble Tea interactive terminal with 7 display modes, keybindings, voice input
- **CLI** (`--cmd`): Send command, print response, exit. `--json` for full envelope output.

## Features

- Clean text prompt: "Listening..." as default state
- Voice input capture (via ALSA / audio-mcp)
- Text input for keyboard interaction
- Direct framebuffer overlay for images/video (communicates with display-mcp)
- Connects to cognitiveosd via Unix socket
- Non-interactive mode for scripting and automation

## Build

```bash
make build    # compile to build/bin/cognitiveos-cli
make test     # run tests
make lint     # go vet
make clean    # remove build artifacts
```

## Architecture

Both modes use the same `internal/client` package to communicate with `cognitiveosd` via Unix socket. The TUI is thin — it captures input and displays output. All intelligence lives in cognitiveosd and the Wide Model. Either interface can crash and restart without affecting the OS.

```
cmd/cognitiveos-cli/   Entry point, flag parsing
internal/
├── client/            Unix socket client (shared by both modes)
├── cli/               Non-interactive CLI mode
└── tui/               Interactive TUI mode (Bubble Tea)
```

## Dependencies

- `github.com/charmbracelet/bubbletea` — TUI framework (TUI mode only)
- `github.com/charmbracelet/lipgloss` — terminal styling (TUI mode only)
- CognitiveOS internal: `cognitiveosd` daemon socket

## Cloning Convention
- Use SSH (git@github.com:) for development.
- Use HTTPS (https://github.com/) for build scripts that clone public dependencies.
