package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Width(60)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			Align(lipgloss.Center).
			Width(56)

	indicatorDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	readyText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true).
			Align(lipgloss.Center).
			Width(56)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Align(lipgloss.Center).
			Width(56)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB"))

	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Align(lipgloss.Center).
			Width(56)

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#EF4444")).
			Padding(0, 1).
			Width(50).
			Align(lipgloss.Center)

	suggestionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Italic(true)

	keyHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Align(lipgloss.Center).
			Width(56)

	codeEntryStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F59E0B")).
			Padding(0, 2).
			Width(50).
			Align(lipgloss.Center)

	codeTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B")).
			Align(lipgloss.Center).
			Width(46)

	maskedInputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	codeHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Align(lipgloss.Center).
			Width(46)

	shutdownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true).
			Align(lipgloss.Center).
			Width(56)

	mediaBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true).
			Align(lipgloss.Center).
			Width(56)

	mediaTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB")).
			Align(lipgloss.Center).
			Width(56)

	mediaHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Align(lipgloss.Center).
			Width(56)

	codeBlockStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 2)

	listStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	tableStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Underline(true)

	actionActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Bold(true).
				Padding(0, 1)

	actionInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280")).
				Padding(0, 1)

	scrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280")).
				Italic(true).
				Padding(0, 2)

	voiceWaveformStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Bold(true).
				Align(lipgloss.Center).
				Width(56)
)
