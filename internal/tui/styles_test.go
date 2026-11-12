package tui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestAppStyleHasPadding(t *testing.T) {
	s := appStyle
	top, right, bottom, left := s.GetPadding()
	if top != 1 || right != 2 || bottom != 1 || left != 2 {
		t.Fatalf("expected padding (1,2,1,2), got (%d,%d,%d,%d)", top, right, bottom, left)
	}
}

func TestAppStyleHasWidth(t *testing.T) {
	s := appStyle
	w := s.GetWidth()
	if w != 60 {
		t.Fatalf("expected width 60, got %d", w)
	}
}

func TestTitleStyleIsBold(t *testing.T) {
	s := titleStyle
	if !s.GetBold() {
		t.Fatal("expected titleStyle to be bold")
	}
}

func TestTitleStyleIsCentered(t *testing.T) {
	s := titleStyle
	align := s.GetAlign()
	if align != lipgloss.Center {
		t.Fatalf("expected titleStyle to be center aligned, got %v", align)
	}
}

func TestTitleStyleHasPurpleColor(t *testing.T) {
	s := titleStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#7C3AED" {
		t.Fatalf("expected titleStyle foreground #7C3AED, got %s", fg)
	}
}

func TestIndicatorDotHasGreenColor(t *testing.T) {
	s := indicatorDot
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#10B981" {
		t.Fatalf("expected indicatorDot foreground #10B981, got %s", fg)
	}
}

func TestReadyTextIsItalic(t *testing.T) {
	s := readyText
	if !s.GetItalic() {
		t.Fatal("expected readyText to be italic")
	}
}

func TestPromptStyleIsBold(t *testing.T) {
	s := promptStyle
	if !s.GetBold() {
		t.Fatal("expected promptStyle to be bold")
	}
}

func TestInputStyleHasLightForeground(t *testing.T) {
	s := inputStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#E5E7EB" {
		t.Fatalf("expected inputStyle foreground #E5E7EB, got %s", fg)
	}
}

func TestOutputStyleHasLightForeground(t *testing.T) {
	s := outputStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#D1D5DB" {
		t.Fatalf("expected outputStyle foreground #D1D5DB, got %s", fg)
	}
}

func TestSpinnerStyleIsOrange(t *testing.T) {
	s := spinnerStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#F59E0B" {
		t.Fatalf("expected spinnerStyle foreground #F59E0B, got %s", fg)
	}
	if s.GetAlign() != lipgloss.Center {
		t.Fatal("expected spinnerStyle to be center aligned")
	}
}

func TestErrorBoxStyleExists(t *testing.T) {
	s := errorBoxStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) == "" {
		t.Fatal("expected errorBoxStyle to have foreground color")
	}
}

func TestMaskedInputStyleHasOrangeForeground(t *testing.T) {
	s := maskedInputStyle
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#F59E0B" {
		t.Fatalf("expected maskedInputStyle foreground #F59E0B, got %s", fg)
	}
}

func TestShutdownStyleIsRedBold(t *testing.T) {
	s := shutdownStyle
	if !s.GetBold() {
		t.Fatal("expected shutdownStyle to be bold")
	}
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#EF4444" {
		t.Fatalf("expected shutdownStyle foreground #EF4444, got %s", fg)
	}
}

func TestMediaBarStyleIsGreenBold(t *testing.T) {
	s := mediaBarStyle
	if !s.GetBold() {
		t.Fatal("expected mediaBarStyle to be bold")
	}
	fg := s.GetForeground()
	if fmt.Sprintf("%s", fg) != "#10B981" {
		t.Fatalf("expected mediaBarStyle foreground #10B981, got %s", fg)
	}
}
