# cli — CognitiveOS TUI

The human interface — a Go Bubble Tea TUI with 7 display modes. Replaces the traditional desktop/app paradigm with a clean terminal-based prompt.

## Display Modes

| Mode | Description |
|------|-------------|
| **idle** | Default — minimal "Listening..." prompt |
| **listening** | Shows "Listening..." while waiting for input |
| **processing** | "Thinking..." with spinning indicator |
| **responding** | Streaming AI response output |
| **error** | Red error state with message |
| **code entry** | Multi-line text input for code blocks |

## Keybindings

| Key | Action |
|-----|--------|
| `Esc` | Cancel / back |
| `Ctrl+C` | Quit |
| `Enter` | Submit text |
| `Tab` | Cycle display modes (debug) |

## Build

```bash
go build -o bin/cognitiveos-cli ./cmd/cognitiveos-cli
```

The TUI connects to cognitiveosd via Unix socket at `/cognitiveos/run/daemon.sock` with 30s retry. It is thin by design — all intelligence lives in the daemon and Wide Model. The TUI can crash and restart without affecting the OS.

## Dependencies

- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`

## Author

**Jean Machuca** — [GitHub](https://github.com/jeanmachuca) · [Sponsor](https://github.com/sponsors/jeanmachuca)

## License

MIT
