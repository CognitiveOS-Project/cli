package tui

import (
	"strings"
	"testing"

	"github.com/CognitiveOS-Project/cli/internal/client"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestConn() *client.Conn {
	return client.New("/test.sock")
}

func teaKeyMsgCtrlC() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyCtrlC}
}

func teaKeyMsgCtrlD() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyCtrlD}
}

func teaKeyMsgBackspace() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyBackspace}
}

func teaKeyMsgCtrlL() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyCtrlL}
}

func teaKeyMsgEnter() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyEnter}
}

func TestStylesRenderNonEmpty(t *testing.T) {
	styles := []string{
		titleStyle.Render("CognitiveOS"),
		promptStyle.Render("> "),
		inputStyle.Render("test"),
		outputStyle.Render("output"),
		spinnerStyle.Render("..."),
		errorBoxStyle.Render("error"),
		suggestionStyle.Render("try this"),
		keyHintStyle.Render("[Enter] send"),
		codeEntryStyle.Render("code entry"),
		codeTitleStyle.Render("Code"),
		codeHintStyle.Render("[Enter] submit"),
		maskedInputStyle.Render("●●●●"),
		shutdownStyle.Render("Shutdown?"),
		mediaBarStyle.Render("[FRAMEBUFFER]"),
		mediaTitleStyle.Render("Media: image"),
		mediaHintStyle.Render("close"),
		codeBlockStyle.Render("```"),
		listStyle.Render("  - item"),
		tableStyle.Render("| a | b |"),
		urlStyle.Render("https://example.com"),
		actionActiveStyle.Render("[retry]"),
		actionInactiveStyle.Render("[close]"),
		scrollIndicatorStyle.Render("▲"),
		voiceWaveformStyle.Render("🎤 ▁"),
	}

	for _, s := range styles {
		if s == "" {
			t.Fatalf("expected non-empty style output")
		}
	}
}

func TestRenderURLs(t *testing.T) {
	result := renderURLs("Visit https://example.com for details")
	if !strings.Contains(result, "https://example.com") {
		t.Fatalf("expected URL to be preserved in output, got: %s", result)
	}
}

func TestRenderURLsNoURL(t *testing.T) {
	result := renderURLs("No URL here")
	if result != "No URL here" {
		t.Fatalf("expected unchanged output, got: %s", result)
	}
}

func TestHandleIdleKeyCtrlC(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateIdle

	result, _ := m.handleKey(teaKeyMsgCtrlC())
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatal("expected idle state after ctrl+c")
	}
}

func TestHandleIdleKeyCtrlD(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateIdle

	result, _ := m.handleKey(teaKeyMsgCtrlD())
	m2 := result.(Model)

	if !m2.confirming {
		t.Fatal("expected confirming=true after ctrl+d")
	}
}

func TestHandleListeningKeyBackspace(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateListening
	m.input.WriteString("hello")

	result, _ := m.handleKey(teaKeyMsgBackspace())
	m2 := result.(Model)

	if m2.input.String() != "hell" {
		t.Fatalf("expected 'hell', got '%s'", m2.input.String())
	}
}

func TestHandleListeningKeyCtrlL(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateListening

	result, _ := m.handleKey(teaKeyMsgCtrlL())
	m2 := result.(Model)

	if m2.state != StateListening {
		t.Fatal("expected listening state after ctrl+l")
	}
}

func TestHandleCodeEntryKeyEnterTooShort(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateCodeEntry
	m.codeInput.WriteString("ABC")

	result, _ := m.handleKey(teaKeyMsgEnter())
	m2 := result.(Model)

	if m2.state != StateCodeEntry {
		t.Fatal("expected to stay in code entry with short code")
	}
}

func TestHandleCodeEntryKeyEnterLongEnough(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateCodeEntry
	m.codeInput.WriteString("ABCD12")

	result, _ := m.handleKey(teaKeyMsgEnter())
	m2 := result.(Model)

	if m2.state != StateIdle {
		t.Fatal("expected idle state after valid code entry")
	}
}

func TestHandleCodeEntryKeyBackspace(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateCodeEntry
	m.codeInput.WriteString("ABCD")

	result, _ := m.handleKey(teaKeyMsgBackspace())
	m2 := result.(Model)

	if m2.codeInput.String() != "ABC" {
		t.Fatalf("expected 'ABC', got '%s'", m2.codeInput.String())
	}
}

func TestSpinnerCycles(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateProcessing

	initialIdx := m.spinnerIdx
	result, _ := m.Update(spinnerTickMsg{})
	m2 := result.(Model)

	expectedIdx := (initialIdx + 1) % len(m.spinnerChars)
	if m2.spinnerIdx != expectedIdx {
		t.Fatalf("expected spinner index %d, got %d", expectedIdx, m2.spinnerIdx)
	}
}

func TestVoiceRecordingFrameAdvances(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateVoiceRecording

	initialIdx := m.waveformIdx
	result, _ := m.Update(spinnerTickMsg{})
	m2 := result.(Model)

	if m2.waveformIdx != initialIdx+1 {
		t.Fatalf("expected waveform index %d, got %d", initialIdx+1, m2.waveformIdx)
	}
}

func TestHistoryNavigationDownToBottom(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateListening
	m.history = []string{"first", "second"}
	m.historyIdx = 0

	msgDown := tea.KeyMsg{Type: tea.KeyDown}
	result, _ := m.Update(msgDown)
	m2 := result.(Model)

	if m2.input.String() != "second" {
		t.Fatalf("expected 'second' from history, got '%s'", m2.input.String())
	}
}

func TestRenderMedia(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateMedia
	m.ready = true
	m.mediaPayload.Type = "image"
	m.mediaPayload.Paths = []string{"/tmp/photo.jpg"}

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty media view")
	}
	if !strings.Contains(view, "FRAMEBUFFER") {
		t.Fatalf("media view should contain FRAMEBUFFER, got: %s", view)
	}
}

func TestRenderError(t *testing.T) {
	conn := newTestConn()
	m := NewModel(conn)
	m.state = StateError
	m.ready = true
	m.errMsg = "connection refused"

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty error view")
	}
	if !strings.Contains(view, "connection refused") {
		t.Fatalf("error view should contain error message, got: %s", view)
	}
}
