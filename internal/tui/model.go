package tui

import (
	"encoding/json"
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
}

type outputMsg string
type statusMsg string
type connStatusMsg connectionStatus

func NewModel(conn *client.Conn) Model {
	return Model{
		state:      StateIdle,
		conn:       conn,
		connStatus: ConnDisconnected,
		spinnerChars: []string{"|", "/", "-", "\\"},
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
		for i := 0; i < 30; i++ {
			if err := conn.Connect(); err == nil {
				conn.RequestStatus()
				return connStatusMsg(ConnConnected)
			}
			time.Sleep(time.Second)
		}
		return connStatusMsg(ConnFailed)
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
			return m, listenCmd(m.conn)
		}
		return m, nil

	case outputMsg:
		m.output.Reset()
		m.output.WriteString(string(msg))
		m.lastOutput = string(msg)
		m.state = StateResponding
		return m, nil

	case statusMsg:
		m.output.Reset()
		m.output.WriteString(string(msg))
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

func listenCmd(conn *client.Conn) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case env, ok := <-conn.Messages:
				if !ok {
					return connStatusMsg(ConnDisconnected)
				}
				switch env.Type {
				case "output_deliver":
					var payload struct {
						Content     string `json:"content"`
						ContentType string `json:"content_type"`
					}
					if err := 	json.Unmarshal(env.Payload, &payload); err == nil {
						return outputMsg(payload.Content)
					}
				case "status_response":
					return statusMsg("connected")
				case "input_accepted":
					return statusMsg("sent")
				case "audit_report":
					return statusMsg("audit received")
				}
			}
		}
	}
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	if m.connStatus == ConnFailed {
		return m.renderError("\n\n  ⚠ SYSTEM HALTED\n\n  Raw Model integrity check failed.\n  CognitiveOS cannot operate without a verified guardrail.\n\n  Please reflash firmware.\n")
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
	b.WriteString(indicatorDot.Render("●") + " " + readyText.Render("ready"))
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
	for _, line := range lines {
		b.WriteString(promptStyle.Render("> ") + outputStyle.Render(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(m.renderOutput(m.lastOutput))
	b.WriteString("\n")
	b.WriteString(keyHintStyle.Render("─────────────────────"))
	b.WriteString("\n")
	b.WriteString(promptStyle.Render("> ") + inputStyle.Render(m.input.String()))
	b.WriteString("\n\n")
	b.WriteString(keyHintStyle.Render("[Enter] send  [Esc] idle  [Tab] actions"))
	b.WriteString("\n")
	return appStyle.Render(b.String())
}

func (m Model) renderOutput(text string) string {
	var b strings.Builder
	for _, line := range strings.Split(text, "\n") {
		b.WriteString(outputStyle.Render(line))
		b.WriteString("\n")
	}
	return b.String()
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
