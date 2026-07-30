package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CognitiveOS-Project/cli/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateIdle State = iota
	StateListening
	StateProcessing
	StateResponding
	StateMedia
	StateError
	StateCodeEntry
)

type connectionStatus int

const (
	ConnDisconnected connectionStatus = iota
	ConnConnecting
	ConnConnected
	ConnFailed
)

type Model struct {
	state        State
	connStatus   connectionStatus
	conn         *client.Conn
	input        strings.Builder
	output       strings.Builder
	lastOutput   string
	errMsg       string
	suggestion   string
	codeInput    strings.Builder
	history      []string
	historyIdx   int
	spinnerIdx   int
	spinnerChars []string
	width        int
	height       int
	ready        bool
	confirming   bool

	mediaPayload struct {
		Type  string   `json:"type"`
		Paths []string `json:"paths"`
	}

	actions     []string
	actionIdx   int
	showActions bool
	scrollOffset int
}

type outputMsg struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

type statusMsg string
type connStatusMsg connectionStatus
type reconnectMsg struct{}
type cancelProcessingMsg struct{}

func NewModel(conn *client.Conn) Model {
	return Model{
		state:        StateIdle,
		conn:         conn,
		connStatus:   ConnDisconnected,
		spinnerChars: []string{".", "..", "..."},
		actions:      []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		connectCmd(m.conn),
		spinnerTickCmd(),
	)
}

func connectCmd(conn *client.Conn) tea.Cmd {
	return func() tea.Msg {
		for i := 0; ; i++ {
			if err := conn.Connect(); err == nil {
				_ = conn.RequestStatus()
				return connStatusMsg(ConnConnected)
			}
			if i >= 30 {
				return connStatusMsg(ConnFailed)
			}
			time.Sleep(time.Second)
		}
	}
}

func reconnectCmd(conn *client.Conn) tea.Cmd {
	return func() tea.Msg {
		for {
			if err := conn.Connect(); err == nil {
				_ = conn.RequestStatus()
				return connStatusMsg(ConnConnected)
			}
			time.Sleep(time.Second)
		}
	}
}

func spinnerTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

type spinnerTickMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case connStatusMsg:
		m.connStatus = connectionStatus(msg)
		if m.connStatus == ConnConnected {
			if err := m.conn.RequestStatus(); err != nil {
				m.state = StateError
				m.errMsg = "failed to request initial status: " + err.Error()
			}
			return m, listenCmd(m.conn, &m)
		}
		if m.connStatus == ConnFailed {
			return m, reconnectCmd(m.conn)
		}
		return m, nil

	case reconnectMsg:
		if m.connStatus == ConnFailed || m.connStatus == ConnDisconnected {
			if err := m.conn.Connect(); err == nil {
				_ = m.conn.RequestStatus()
				m.connStatus = ConnConnected
				return m, listenCmd(m.conn, &m)
			}
		}
		return m, reconnectCmd(m.conn)

	case outputMsg:
		m.output.Reset()
		m.output.WriteString(msg.Content)
		m.lastOutput = msg.Content
		m.scrollOffset = 0
		if msg.ContentType == "media" {
			m.state = StateMedia
			m.actions = []string{"close", "save"}
		} else {
			m.state = StateResponding
			m.actions = []string{"retry", "close"}
		}
		m.actionIdx = 0
		m.showActions = false
		return m, nil

	case statusMsg:
		m.output.Reset()
		m.output.WriteString(string(msg))
		return m, nil

	case cancelProcessingMsg:
		m.state = StateIdle
		m.input.Reset()
		return m, nil

	case spinnerTickMsg:
		if m.state == StateProcessing {
			m.spinnerIdx = (m.spinnerIdx + 1) % len(m.spinnerChars)
		}
		return m, spinnerTickCmd()

	default:
		return m, nil
	}
}

func listenCmd(conn *client.Conn, model *Model) tea.Cmd {
	return func() tea.Msg {
		for env := range conn.Messages {
			switch env.Type {
			case "output_deliver":
				var payload struct {
					Content     string `json:"content"`
					ContentType string `json:"content_type"`
					Media       *struct {
						Type  string   `json:"type"`
						Paths []string `json:"paths"`
					} `json:"media,omitempty"`
				}
				if err := json.Unmarshal(env.Payload, &payload); err == nil {
					if model != nil && payload.Media != nil {
						model.mediaPayload.Type = payload.Media.Type
						model.mediaPayload.Paths = payload.Media.Paths
					}
					return outputMsg{Content: payload.Content, ContentType: payload.ContentType}
				}
			case "status_response":
				return statusMsg("connected")
			case "input_accepted":
				return statusMsg("sent")
			case "audit_report":
				return statusMsg("audit received")
			}
		}
		return reconnectMsg{}
	}
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	if m.connStatus == ConnConnecting {
		return "\n\n  Connecting...\n"
	}
	if m.connStatus == ConnFailed {
		return "\n\n  ⚠ Daemon unavailable. Check system.\n\n  Retrying connection...\n"
	}

	switch m.state {
	case StateIdle:
		return m.renderIdle()
	case StateListening:
		return m.renderListening()
	case StateProcessing:
		return m.renderProcessing()
	case StateResponding:
		return m.renderResponding()
	case StateMedia:
		return m.renderMedia()
	case StateError:
		return m.renderError(m.errMsg)
	case StateCodeEntry:
		return m.renderCodeEntry()
	default:
		return m.renderIdle()
	}
}

func (m Model) renderIdle() string {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(titleStyle.Render("CognitiveOS"))
	b.WriteString("\n\n")
	if m.confirming {
		b.WriteString(shutdownStyle.Render("Shutdown system?  [Y]es  [N]o"))
	} else {
		b.WriteString(indicatorDot.Render("●") + " " + readyText.Render("ready"))
	}
	b.WriteString("\n\n\n")
	b.WriteString(hintStyle.Render("(press / to speak, type anything to begin)"))
	b.WriteString("\n")
	return appStyle.Render(b.String())
}

func (m Model) renderListening() string {
	var b strings.Builder
	b.WriteString("\n")
	lines := strings.Split(m.output.String(), "\n")
	for _, line := range lines {
		if line != "" {
			b.WriteString(promptStyle.Render("> ") + outputStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(promptStyle.Render("> ") + inputStyle.Render(m.input.String()))
	b.WriteString("\n\n")
	b.WriteString(keyHintStyle.Render("[Enter] send  [Esc] cancel"))
	b.WriteString("\n")
	return appStyle.Render(b.String())
}

func (m Model) renderProcessing() string {
	var b strings.Builder
	b.WriteString("\n")
	lines := strings.Split(m.output.String(), "\n")
	for _, line := range lines {
		if line != "" {
			b.WriteString(promptStyle.Render("> ") + outputStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	spinner := m.spinnerChars[m.spinnerIdx]
	b.WriteString(spinnerStyle.Render(spinner + " Working..."))
	b.WriteString("\n\n")
	b.WriteString(keyHintStyle.Render("[Ctrl+C] cancel"))
	b.WriteString("\n")
	return appStyle.Render(b.String())
}

func (m Model) renderResponding() string {
	var b strings.Builder
	b.WriteString("\n")
	lines := strings.Split(m.output.String(), "\n")
	visibleHeight := m.height - 8
	if m.lastOutput != "" {
		lastLines := strings.Count(m.lastOutput, "\n") + 1
		visibleHeight -= lastLines
	}
	visibleHeight -= 3 // separator + input + hints
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	start := m.scrollOffset
	if start >= len(lines) {
		start = 0
	}
	end := start + visibleHeight
	if end > len(lines) {
		end = len(lines)
	}
	if m.scrollOffset > 0 {
		b.WriteString(scrollIndicatorStyle.Render("▲"))
		b.WriteString("\n")
	}
	for _, line := range lines[start:end] {
		b.WriteString(promptStyle.Render("> ") + outputStyle.Render(line))
		b.WriteString("\n")
	}
	if end < len(lines) {
		b.WriteString(scrollIndicatorStyle.Render(fmt.Sprintf("▼ %d more lines", len(lines)-end)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(m.renderOutput(m.lastOutput))
	b.WriteString("\n")
	b.WriteString(keyHintStyle.Render("─────────────────────"))
	b.WriteString("\n")
	b.WriteString(promptStyle.Render("> ") + inputStyle.Render(m.input.String()))
	b.WriteString("\n\n")
	if m.showActions && len(m.actions) > 0 {
		for i, a := range m.actions {
			if i == m.actionIdx {
				b.WriteString(actionActiveStyle.Render("[" + a + "]"))
			} else {
				b.WriteString(actionInactiveStyle.Render("[" + a + "]"))
			}
			b.WriteString(" ")
		}
		b.WriteString("\n")
	} else if m.scrollOffset > 0 || end < len(lines) {
		b.WriteString(keyHintStyle.Render("[Enter] send  [Esc] idle  [Shift+↑/↓] scroll"))
		b.WriteString("\n")
	} else {
		b.WriteString(keyHintStyle.Render("[Enter] send  [Esc] idle  [Tab] actions"))
		b.WriteString("\n")
	}
	return appStyle.Render(b.String())
}

func (m Model) renderOutput(text string) string {
	var b strings.Builder
	inCodeBlock := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				b.WriteString(codeBlockStyle.Render(line))
			} else {
				b.WriteString(codeBlockStyle.Render(line))
			}
			b.WriteString("\n")
			continue
		}

		if inCodeBlock {
			b.WriteString(codeBlockStyle.Render(line))
			b.WriteString("\n")
			continue
		}

		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			b.WriteString(listStyle.Render("  " + trimmed))
			b.WriteString("\n")
			continue
		}

		if len(trimmed) > 2 && trimmed[0] >= '1' && trimmed[0] <= '9' && trimmed[1] == '.' {
			b.WriteString(listStyle.Render("  " + trimmed))
			b.WriteString("\n")
			continue
		}

		if strings.Count(trimmed, "|") >= 3 && strings.Contains(trimmed, "-") {
			b.WriteString(tableStyle.Render(trimmed))
			b.WriteString("\n")
			continue
		}

		rendered := renderURLs(line)
		b.WriteString(outputStyle.Render(rendered))
		b.WriteString("\n")
	}
	return b.String()
}

func renderURLs(text string) string {
	var b strings.Builder
	remaining := text
	for {
		idx := strings.Index(remaining, "http")
		if idx < 0 {
			b.WriteString(remaining)
			break
		}
		b.WriteString(remaining[:idx])
		url := remaining[idx:]
		end := strings.IndexAny(url, " \t\n\r,.;:!?\"')")
		if end > 0 {
			url = url[:end]
		}
		b.WriteString(urlStyle.Render(url))
		remaining = remaining[idx+len(url):]
	}
	return b.String()
}

func (m Model) renderMedia() string {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(mediaBarStyle.Render("[FRAMEBUFFER]"))
	b.WriteString("\n\n")
	b.WriteString(mediaTitleStyle.Render("Media: " + m.mediaPayload.Type))
	b.WriteString("\n\n")
	b.WriteString(mediaHintStyle.Render("Type \"close\" to return, \"save\" to store"))
	b.WriteString("\n\n")
	b.WriteString(promptStyle.Render("> ") + inputStyle.Render(m.input.String()))
	return appStyle.Render(b.String())
}

func (m Model) renderError(msg string) string {
	var b strings.Builder
	b.WriteString("\n")
	lines := strings.Split(m.output.String(), "\n")
	for _, line := range lines {
		if line != "" {
			b.WriteString(promptStyle.Render("> ") + outputStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(errorBoxStyle.Render("⚠ " + msg))
	b.WriteString("\n")
	if m.suggestion != "" {
		b.WriteString("\n")
		b.WriteString(suggestionStyle.Render("Try: " + m.suggestion))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(promptStyle.Render("> ") + inputStyle.Render(m.input.String()))
	b.WriteString("\n")
	return appStyle.Render(b.String())
}

func (m Model) renderCodeEntry() string {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(codeTitleStyle.Render("⚠ System Code Required"))
	b.WriteString("\n\n")
	b.WriteString("    Enter unlock code:")
	b.WriteString("\n    ")
	masked := strings.Repeat("●", m.codeInput.Len())
	if len(masked) < 9 {
		masked += "▊"
	}
	b.WriteString(maskedInputStyle.Render(masked))
	b.WriteString("\n\n")
	b.WriteString(codeHintStyle.Render("[Enter] submit  [Esc] cancel"))
	b.WriteString("\n")
	return codeEntryStyle.Render(b.String())
}
