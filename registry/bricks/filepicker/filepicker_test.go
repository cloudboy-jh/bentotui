package filepicker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestWindowSizeMsgAppliesSize(t *testing.T) {
	m := New(".")
	_, _ = m.Update(tea.WindowSizeMsg{Width: 44, Height: 9})
	w, h := m.GetSize()
	if w != 44 || h != 9 {
		t.Fatalf("expected 44x9, got %dx%d", w, h)
	}
}

func TestDidSelectFileSetsSelectedPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pkg.go")
	if err := os.WriteFile(path, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	m := New(dir)
	initCmd := m.Init()
	if initCmd == nil {
		t.Fatal("expected init command")
	}
	_, _ = m.Update(initCmd())
	_, _ = m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	if got := m.SelectedPath(); cleanPath(got) != cleanPath(path) {
		t.Fatalf("expected selected %q, got %q", path, got)
	}
	if !strings.Contains(m.Status(), "selected") {
		t.Fatalf("expected selected status, got %q", m.Status())
	}
}

func TestDidSelectDisabledFileSetsBlockedStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	m := New(dir)
	m.SetAllowedTypes(".go")
	initCmd := m.Init()
	if initCmd == nil {
		t.Fatal("expected init command")
	}
	_, _ = m.Update(initCmd())
	_, _ = m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))

	if got := m.SelectedPath(); got != "" {
		t.Fatalf("expected empty selected path for disabled file, got %q", got)
	}
	if !strings.Contains(m.Status(), "blocked") {
		t.Fatalf("expected blocked status, got %q", m.Status())
	}
}

func TestCleanPathDefaultsToDot(t *testing.T) {
	if got := cleanPath("   "); got != "." {
		t.Fatalf("expected dot path, got %q", got)
	}
}
