package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateIdle:
		return m.handleIdleKey(msg)
	case StateListening:
		return m.handleListeningKey(msg)
	case StateProcessing:
		return m.handleProcessingKey(msg)
	case StateResponding:
		return m.handleRespondingKey(msg)
	case StateError:
		return m.handleErrorKey(msg)
	case StateCodeEntry:
		return m.handleCodeEntryKey(msg)
	default:
		return m, nil
	}
}

func (m Model) handleIdleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "ctrl+d":
		return m, tea.Quit
	case "ctrl+alt+s":
		m.conn.SendSystemCode("security", "")
		return m, nil
	case "/":
		m.conn.SendVoice()
		return m, nil
	case "esc":
		return m, nil
	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
			m.state = StateListening
		}
		return m, nil
	}
}

func (m Model) handleListeningKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		input := m.input.String()
		if input == "" {
			return m, nil
		}
		m.output.Reset()
		m.output.WriteString(input)
		m.history = append(m.history, input)
		m.historyIdx = len(m.history)
		m.input.Reset()
		m.state = StateProcessing
		m.conn.SendInput(input)
		m.conn.RequestStatus()
		return m, nil

	case "esc":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "ctrl+c":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "up":
		if len(m.history) > 0 && m.historyIdx > 0 {
			m.historyIdx--
			m.input.Reset()
			m.input.WriteString(m.history[m.historyIdx])
		}
		return m, nil

	case "down":
		if m.historyIdx < len(m.history)-1 {
			m.historyIdx++
			m.input.Reset()
			m.input.WriteString(m.history[m.historyIdx])
		} else {
			m.input.Reset()
			m.historyIdx = len(m.history)
		}
		return m, nil

	case "backspace":
		s := m.input.String()
		if len(s) > 0 {
			m.input.Reset()
			m.input.WriteString(s[:len(s)-1])
		}
		return m, nil

	case "ctrl+l":
		return m, tea.ClearScreen

	case "/":
		m.conn.SendVoice()
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

func (m Model) handleProcessingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.state = StateIdle
		m.input.Reset()
		return m, nil
	case "ctrl+l":
		return m, tea.ClearScreen
	default:
		return m, nil
	}
}

func (m Model) handleRespondingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		input := m.input.String()
		if input == "" {
			return m, nil
		}
		if strings.ToLower(input) == "close" || strings.ToLower(input) == "close_media" {
			if m.state == StateMedia {
				m.state = StateResponding
			}
			return m, nil
		}
		m.output.Reset()
		m.output.WriteString(input)
		m.history = append(m.history, input)
		m.historyIdx = len(m.history)
		m.input.Reset()
		m.state = StateProcessing
		m.conn.SendInput(input)
		return m, nil

	case "esc":
		m.state = StateIdle
		m.input.Reset()
		m.output.Reset()
		return m, nil

	case "ctrl+c":
		m.state = StateIdle
		m.input.Reset()
		return m, nil

	case "tab":
		return m, nil

	case "backspace":
		s := m.input.String()
		if len(s) > 0 {
			m.input.Reset()
			m.input.WriteString(s[:len(s)-1])
		}
		return m, nil

	case "ctrl+l":
		return m, tea.ClearScreen

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

func (m Model) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateIdle
		m.errMsg = ""
		m.suggestion = ""
		m.input.Reset()
		return m, nil

	case "enter":
		input := m.input.String()
		if input != "" {
			m.output.Reset()
			m.output.WriteString(input)
			m.history = append(m.history, input)
			m.historyIdx = len(m.history)
			m.input.Reset()
			m.state = StateProcessing
			m.conn.SendInput(input)
		} else {
			m.state = StateIdle
		}
		m.errMsg = ""
		m.suggestion = ""
		return m, nil

	case "backspace":
		s := m.input.String()
		if len(s) > 0 {
			m.input.Reset()
			m.input.WriteString(s[:len(s)-1])
		}
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

func (m Model) handleCodeEntryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		code := m.codeInput.String()
		if len(code) >= 4 {
			m.conn.SendSystemCode("unlock", code)
			m.codeInput.Reset()
			m.state = StateIdle
		}
		return m, nil

	case "esc":
		m.codeInput.Reset()
		m.state = StateIdle
		return m, nil

	case "backspace":
		s := m.codeInput.String()
		if len(s) > 0 {
			m.codeInput.Reset()
			m.codeInput.WriteString(s[:len(s)-1])
		}
		return m, nil

	case "-":
		if m.codeInput.Len() > 0 && !strings.HasSuffix(m.codeInput.String(), "-") {
			m.codeInput.WriteString("-")
		}
		return m, nil

	default:
		if msg.Type == tea.KeyRunes && m.codeInput.Len() < 9 {
			for _, r := range msg.String() {
				if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' {
					m.codeInput.WriteString(strings.ToUpper(string(r)))
					if m.codeInput.Len() == 4 || m.codeInput.Len() == 9 {
						// auto-dash handled by typing it
					}
				}
			}
		}
		return m, nil
	}
}
