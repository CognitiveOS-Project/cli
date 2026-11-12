package tui

import (
	"testing"

	"github.com/CognitiveOS-Project/cli/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleIdleKeyCtrlCQuits(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if m2.state != m.state {
		t.Fatal("expected same state on ctrl+c in idle")
	}
}

func TestHandleIdleKeyCtrlDStartsConfirming(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if !m2.confirming {
		t.Fatal("expected confirming=true after ctrl+d in idle")
	}
}

func TestHandleIdleKeyYWithoutConfirming(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if m2.confirming {
		t.Fatal("expected confirming=false without prior confirmation")
	}
	if m2.state != StateIdle {
		t.Fatal("expected idle state")
	}
}

func TestHandleIdleKeyYWithConfirming(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.confirming = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if m2.confirming {
		t.Fatal("expected confirming=false after yes")
	}
}

func TestHandleIdleKeyNWithoutConfirming(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if m2.confirming {
		t.Fatal("expected confirming=false without prior confirmation")
	}
}

func TestHandleIdleKeySlashStartsVoice(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	result, _ := m.handleIdleKey(msg)
	m2 := result.(Model)

	if m2.state != StateVoiceRecording {
		t.Fatalf("expected voice recording state, got %d", m2.state)
	}
}

func TestHandleListeningKeyEnterSubmits(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("hello")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.state != StateProcessing {
		t.Fatalf("expected processing state, got %d", m2.state)
	}
	if m2.input.String() != "" {
		t.Fatalf("expected empty input after submit, got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeyEmptyEnterIgnored(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.state != StateListening {
		t.Fatalf("expected listening state for empty input, got %d", m2.state)
	}
}

func TestHandleListeningKeyEsc(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("test")

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle state after esc, got %d", m2.state)
	}
	if m2.input.String() != "" {
		t.Fatalf("expected empty input after esc, got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeyHistoryUp(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.history = []string{"first", "second", "third"}
	m.historyIdx = 3

	msg := tea.KeyMsg{Type: tea.KeyUp}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.input.String() != "third" {
		t.Fatalf("expected 'third' from history, got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeyBackspace(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("hello")

	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.input.String() != "hell" {
		t.Fatalf("expected 'hell' after backspace, got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeySpace(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.state != StateListening {
		t.Fatalf("expected listening state after space, got %d", m2.state)
	}
	if m2.input.String() != " " {
		t.Fatalf("expected space in input, got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeySpaceWhenConfirmingIgnored(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.confirming = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	result, _ := m.handleListeningKey(msg)
	m2 := result.(Model)

	if m2.input.String() != "" {
		t.Fatalf("expected empty input when confirming, got '%s'", m2.input.String())
	}
}

func TestHandleProcessingKeyCtrlC(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateProcessing

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	result, _ := m.handleProcessingKey(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after ctrl+c in processing, got %d", m2.state)
	}
}

func TestHandleRespondingKeyEsc(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateResponding

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.handleRespondingKey(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after esc in responding, got %d", m2.state)
	}
}

func TestHandleCodeEntryKey(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateCodeEntry

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
	result, _ := m.handleCodeEntryKey(msg)
	m2 := result.(Model)

	if m2.codeInput.String() != "A" {
		t.Fatalf("expected code input 'A', got '%s'", m2.codeInput.String())
	}
}

func TestHandleCodeEntryKeyEsc(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateCodeEntry
	m.codeInput.WriteString("ABCD")

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.handleCodeEntryKey(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after esc from code entry, got %d", m2.state)
	}
}
