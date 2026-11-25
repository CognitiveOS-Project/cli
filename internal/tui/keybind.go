package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// handleListeningKey handles key input while in listening state.
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
		_ = m.conn.SendInput(input)
		_ = m.conn.RequestStatus()
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

	case " ":
		if m.confirming {
			return m, nil
		}
		m.input.WriteString(" ")
		return m, nil

	case "/":
		if err := m.conn.SendVoice(); err != nil {
			m.state = StateError
			m.errMsg = "failed to send voice command: " + err.Error()
			return m, nil
		}
		m.state = StateVoiceRecording
		m.waveformIdx = 0
		return m, spinnerTickCmd()

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

// handleProcessingKey handles key input while in processing state.
func (m Model) handleProcessingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		_ = m.conn.SendInput(input)
		_ = m.conn.RequestStatus()
		return m, nil

	case "esc":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "ctrl+c":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "ctrl+l":
		return m, tea.ClearScreen

	case " ":
		if m.confirming {
			return m, nil
		}
		m.input.WriteString(" ")
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
	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

// handleRespondingKey handles key input while in responding state.
func (m Model) handleRespondingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.showActions && len(m.actions) > 0 {
			action := m.actions[m.actionIdx]
			m.showActions = false
			switch action {
			case "retry":
				input := m.output.String()
				if input == "" {
					return m, nil
				}
				m.state = StateProcessing
				_ = m.conn.SendInput(input)
				return m, nil
			case "close":
				m.state = StateIdle
				m.input.Reset()
				m.output.Reset()
				return m, nil
			case "save":
				return m, nil
			}
		}
		input := m.input.String()
		if input == "" {
			return m, nil
		}
		m.output.Reset()
		m.output.WriteString(input)
		m.scrollOffset = 0
		m.history = append(m.history, input)
		m.historyIdx = len(m.history)
		m.input.Reset()
		m.state = StateProcessing
		_ = m.conn.SendInput(input)
		return m, nil

	case "esc":
		if m.state == StateMedia {
			m.state = StateResponding
			return m, nil
		}
		m.showActions = false
		m.state = StateIdle
		m.input.Reset()
		m.output.Reset()
		return m, nil

	case "ctrl+l":
		return m, tea.ClearScreen

	case " ":
		if m.confirming {
			return m, nil
		}
		m.input.WriteString(" ")
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

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

// handleIdleKey handles key input while in idle state.
func (m Model) handleIdleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		_ = m.conn.SendInput(input)
		_ = m.conn.RequestStatus()
		return m, nil

	case "esc":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "ctrl+c":
		m.input.Reset()
		return m, nil
	case "ctrl+d":
		m.confirming = true
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
	case "h", "H":
		m.state = StateListening
		m.input.WriteString("h")
		return m, nil

	case "y", "Y":
		if m.confirming {
			m.confirming = false
			_ = m.conn.SendSystemCode("idle", "")
			return m, nil
		}
		m.state = StateListening
		m.input.WriteString("y")
		return m, nil

	case "/":
		if err := m.conn.SendVoice(); err != nil {
			m.state = StateError
			m.errMsg = "failed to send voice command: " + err.Error()
			return m, nil
		}
		m.state = StateVoiceRecording
		m.waveformIdx = 0
		return m, spinnerTickCmd()

	case " ":
		if m.confirming {
			return m, nil
		}
		m.input.WriteString(" ")
		m.state = StateListening
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

// handleVoiceRecordingKey handles key input while in voice recording state.
func (m Model) handleVoiceRecordingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.state = StateIdle
		return m, nil
	default:
		return m, nil
	}
}

// handleErrorKey handles key input while in error state.
func (m Model) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateIdle
		m.input.Reset()
		return m, nil

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
		_ = m.conn.SendInput(input)
		_ = m.conn.RequestStatus()
		return m, nil

	case "backspace":
		s := m.input.String()
		if len(s) > 0 {
			m.input.Reset()
			m.input.WriteString(s[:len(s)-1])
		}
		return m, nil

	case " ":
		m.input.WriteString(" ")
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.input.WriteString(msg.String())
		}
		return m, nil
	}
}

// handleCodeEntryKey handles key input while in code entry state.
func (m Model) handleCodeEntryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		code := m.input.String()
		if len(code) >= 4 {
			_ = m.conn.SendSystemCode("unlock", code)
			m.input.Reset()
			m.state = StateIdle
		}
		return m, nil

	case "esc":
		m.input.Reset()
		m.state = StateIdle
		return m, nil

	case "backspace":
		s := m.input.String()
		if len(s) > 0 {
			m.input.Reset()
			m.input.WriteString(s[:len(s)-1])
		}
		return m, nil

	case " ":
		m.input.WriteString(" ")
		return m, nil

	default:
		if msg.Type == tea.KeyRunes && m.codeInput.Len() < 9 {
			for _, r := range msg.String() {
				if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' {
					m.codeInput.WriteString(strings.ToUpper(string(r)))
				}
			}
		}
		return m, nil
	}
}

// handleKey dispatches key input to the appropriate state handler.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Alt && (msg.String() == "s" || msg.String() == "ctrl+s") {
		_ = m.conn.SendSystemCode("security", "")
		return m, nil
	}
	if msg.String() == "ctrl+alt+s" {
		_ = m.conn.SendSystemCode("security", "")
		return m, nil
	}

	switch m.state {
	case StateIdle:
		return m.handleIdleKey(msg)
	case StateListening:
		return m.handleListeningKey(msg)
	case StateVoiceRecording:
		return m.handleVoiceRecordingKey(msg)
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