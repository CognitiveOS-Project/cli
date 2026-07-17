# cognitiveos-cli — CognitiveOS Terminal Interface

The human interface to CognitiveOS. Operates in two modes: an interactive TUI (terminal user interface) for live sessions, and a non-interactive CLI for scripting and automation.

## Modes

### TUI Mode (default)

Interactive terminal UI with 7 display modes. Replaces the traditional desktop/app paradigm with a clean terminal-based prompt.

```
cognitiveos-cli              # default
cognitiveos-cli --tui        # explicit
```

| Mode | Description |
|------|-------------|
| **idle** | Default — minimal "Listening..." prompt |
| **listening** | Shows "Listening..." while waiting for input |
| **processing** | "Thinking..." with spinning indicator |
| **responding** | Streaming AI response output |
| **error** | Red error state with message |
| **code entry** | Multi-line text input for code blocks |

| Key | Action |
|-----|--------|
| `Esc` | Cancel / back |
| `Ctrl+C` | Quit |
| `Enter` | Submit text |
| `Tab` | Cycle display modes (debug) |

### CLI Mode (non-interactive)

Send a command, print the response, and exit. Useful for scripting, automation, and pipes.

```
cognitiveos-cli --cmd "what time is it"
cognitiveos-cli --cmd "show me my photos" --json
```

| Flag | Description |
|------|-------------|
| `--cmd <text>` | Send a command and print the response |
| `--json` | Print the full JSON envelope (requires `--cmd`) |

Plain text output (default):
```
$ cognitiveos-cli --cmd "what time is it"
The time is 3:42 PM.
```

JSON output:
```
$ cognitiveos-cli --cmd "what time is it" --json
{
  "type": "output_deliver",
  "id": "a1b2c3d4-...",
  "from": "cognitiveosd",
  "payload": {
    "content": "The time is 3:42 PM."
  }
}
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--socket` | `/cognitiveos/run/daemon.sock` | Daemon socket path |
| `--tui` | `false` | Launch interactive TUI (same as default) |
| `--cmd <text>` | — | Send command, print response, exit |
| `--json` | `false` | Print full JSON envelope (requires `--cmd`) |
| `--version` | — | Print version and exit |
| `--help` | — | Print usage |

## Architecture

Both modes use the same `internal/client` package to communicate with `cognitiveosd` via Unix socket. The TUI is thin — it captures input and displays output. All intelligence lives in the daemon and Wide Model. Either interface can crash and restart without affecting the OS.

```
cognitiveos-cli
├── cmd/cognitiveos-cli/   Entry point, flag parsing
├── internal/
│   ├── client/            Unix socket client (shared by both modes)
│   ├── cli/               Non-interactive CLI mode
│   └── tui/               Interactive TUI mode (Bubble Tea)
```

## Build

```bash
make build    # Compile to build/bin/cognitiveos-cli
make test     # Run tests
make lint     # Run go vet
make clean    # Remove build artifacts
```

## Dependencies

- `github.com/charmbracelet/bubbletea` — TUI framework (TUI mode only)
- `github.com/charmbracelet/lipgloss` — terminal styling (TUI mode only)
- CognitiveOS internal: `cognitiveosd` daemon socket

## Related

- [CognitiveOS](https://github.com/CognitiveOS-Project/cognitiveos) — main project repository
- [cognitive-os.org](https://cognitive-os.org) — project website
- [cognitiveosd](https://github.com/CognitiveOS-Project/cognitiveosd) — daemon that this tool connects to
- [core-mcp-bridges](https://github.com/CognitiveOS-Project/core-mcp-bridges) — display-mcp used for media rendering
- [coginit](https://github.com/CognitiveOS-Project/coginit) — boot manager that supervises this binary
- [Product Specs](https://github.com/CognitiveOS-Project/product-specs) — CLI/TUI specification
- [CognitiveOS Project](https://github.com/CognitiveOS-Project) — GitHub organization

## Contributing

1. Branch from `main`
2. Use topic branches: `feature/<name>`, `fix/<name>`
3. Open a PR to `main` with a clear title and description
4. Merge after review

See the [SDLC repo](https://github.com/CognitiveOS-Project/sdlc) for the full contribution guide, code review standards, and testing strategy.

## Author

**Jean Machuca** — [GitHub](https://github.com/jeanmachuca) · [Sponsor](https://github.com/sponsors/jeanmachuca)

## License

MIT
