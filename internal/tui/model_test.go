package tui

import (
	"strings"
	"testing"

	"github.com/CognitiveOS-Project/cli/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	if m.state != StateIdle {
		t.Fatalf("initial state should be idle, got %d", m.state)
	}
	if m.connStatus != ConnDisconnected {
		t.Fatalf("initial conn status should be disconnected, got %d", m.connStatus)
	}
}

func TestIdleState(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	if m.state != StateIdle {
		t.Fatalf("expected idle state, got %d", m.state)
	}
}

func TestKeyPressTransitionsToListening(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateListening {
		t.Fatalf("expected listening state, got %d", m2.state)
	}
	if m2.input.String() != "h" {
		t.Fatalf("expected input 'h', got '%s'", m2.input.String())
	}
}

func TestEnterSubmitsInput(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("hello world")

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateProcessing {
		t.Fatalf("expected processing state, got %d", m2.state)
	}
	if m2.input.String() != "" {
		t.Fatalf("expected empty input after submit, got '%s'", m2.input.String())
	}
}

func TestEmptyEnterDoesNotSubmit(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	// Empty input should not transition to processing
	if m2.state != StateListening {
		t.Fatalf("expected listening state (empty input), got %d", m2.state)
	}
}

func TestEscReturnsToIdle(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("test")

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle state after esc, got %d", m2.state)
	}
}

func TestCtrlCFromProcessing(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateProcessing

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after ctrl+c, got %d", m2.state)
	}
}

func TestCtrlDFromIdle(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after ctrl+d, got %d", m2.state)
	}
}

func TestBackspace(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("hello")

	for i := 0; i < 3; i++ {
		msg := tea.KeyMsg{Type: tea.KeyBackspace}
		result, _ := m.Update(msg)
		m2 := result.(Model)
		m = m2
	}

	if m.input.String() != "he" {
		t.Fatalf("expected input 'he', got '%s'", m.input.String())
	}
}

func TestHistoryNavigation(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateListening

	m.history = []string{"first", "second", "third"}
	m.historyIdx = 3

	msgUp := tea.KeyMsg{Type: tea.KeyUp}
	result, _ := m.Update(msgUp)
	m2 := result.(Model)

	if m2.input.String() != "third" {
		t.Fatalf("expected 'third' from history, got '%s'", m2.input.String())
	}
}

func TestOutputMsgTransition(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := outputMsg("Hello from Wide Model")
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateResponding {
		t.Fatalf("expected responding state, got %d", m2.state)
	}
	if m2.lastOutput != "Hello from Wide Model" {
		t.Fatalf("expected output, got '%s'", m2.lastOutput)
	}
}

func TestConnStatusDisconnectedMessage(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := connStatusMsg(ConnDisconnected)
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.connStatus != ConnDisconnected {
		t.Fatalf("expected disconnected, got %d", m2.connStatus)
	}
}

func TestConnStatusConnectedMessage(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := connStatusMsg(ConnConnected)
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.connStatus != ConnConnected {
		t.Fatalf("expected connected, got %d", m2.connStatus)
	}
}

func TestViewRendersSomething(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.ready = true
	m.connStatus = ConnConnected

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestViewIdleContainsReady(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.ready = true
	m.connStatus = ConnConnected

	view := m.View()
	if !strings.Contains(view, "ready") {
		t.Fatalf("idle view should contain 'ready', got: %s", view)
	}
}

func TestViewListeningContainsInput(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.ready = true
	m.connStatus = ConnConnected
	m.state = StateListening
	m.input.WriteString("test input")

	view := m.View()
	if !strings.Contains(view, "test input") {
		t.Fatalf("listening view should contain input, got: %s", view)
	}
}

func TestViewProcessingContainsWorking(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.ready = true
	m.connStatus = ConnConnected
	m.state = StateProcessing

	view := m.View()
	if !strings.Contains(view, "Working") {
		t.Fatalf("processing view should contain 'Working', got: %s", view)
	}
}

func TestCodeEntryMode(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateCodeEntry

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.codeInput.String() != "A" {
		t.Fatalf("expected code input 'A', got '%s'", m2.codeInput.String())
	}
}

func TestCodeEntryEsc(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateCodeEntry
	m.codeInput.WriteString("ABCD")

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after esc from code entry, got %d", m2.state)
	}
}

func TestRespondingEsc(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)
	m.state = StateResponding

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatalf("expected idle after esc from responding, got %d", m2.state)
	}
}

func TestWindowSize(t *testing.T) {
	conn := client.New("/test.sock")
	m := NewModel(conn)

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	result, _ := m.Update(msg)
	m2 := result.(Model)

	if !m2.ready {
		t.Fatal("expected ready after WindowSizeMsg")
	}
	if m2.width != 80 {
		t.Fatalf("expected width 80, got %d", m2.width)
	}
}
