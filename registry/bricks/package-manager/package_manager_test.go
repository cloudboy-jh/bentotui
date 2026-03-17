package packagemanager

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestInstallFlowCompletes(t *testing.T) {
	m := New([]string{"a", "b"})
	_, _ = m.Update(installedPkgMsg("a"))
	if m.Done() {
		t.Fatal("expected not done after first package")
	}
	_, _ = m.Update(installedPkgMsg("b"))
	if !m.Done() {
		t.Fatal("expected done after last package")
	}
}

func TestFailureStopsFlow(t *testing.T) {
	m := New([]string{"a"})
	want := errors.New("boom")
	_, _ = m.Update(failedPkgMsg{pkg: "a", err: want})
	if !m.Done() {
		t.Fatal("expected done on failure")
	}
	if m.Error() == nil || m.Error().Error() != want.Error() {
		t.Fatalf("expected failure error %q, got %v", want, m.Error())
	}
}

func TestWindowSizeMsgAppliesSize(t *testing.T) {
	m := New([]string{"a"})
	_, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
	w, h := m.GetSize()
	if w != 80 || h != 10 {
		t.Fatalf("expected size 80x10, got %dx%d", w, h)
	}
}
