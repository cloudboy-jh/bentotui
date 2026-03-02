package dialog

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type spyDialog struct {
	enterCount int
	width      int
	height     int
}

func (d *spyDialog) Init() tea.Cmd { return nil }

func (d *spyDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
		d.enterCount++
	}
	return d, nil
}

func (d *spyDialog) View() tea.View { return tea.NewView("spy") }

func (d *spyDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

func (d *spyDialog) GetSize() (int, int) { return d.width, d.height }

func (d *spyDialog) Title() string { return "Spy" }

func TestCustomDialogReceivesEnter(t *testing.T) {
	m := New()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	spy := &spyDialog{}
	_, _ = m.Update(Open(spy))
	if !m.IsOpen() {
		t.Fatal("expected dialog to open")
	}

	_, _ = m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	if spy.enterCount != 1 {
		t.Fatalf("expected custom dialog to receive enter once, got %d", spy.enterCount)
	}
	if !m.IsOpen() {
		t.Fatal("expected custom dialog to remain open after enter")
	}
}

func TestConfirmDialogEnterClosesAndConfirms(t *testing.T) {
	m := New()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	confirmed := false
	_, _ = m.Update(Open(Confirm{
		DialogTitle: "Confirm",
		Message:     "Proceed",
		OnConfirm: func() tea.Msg {
			confirmed = true
			return nil
		},
	}))

	_, cmd := m.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	if cmd == nil {
		t.Fatal("expected confirm command on enter")
	}
	_ = cmd()
	if !confirmed {
		t.Fatal("expected confirm callback to run")
	}
	if m.IsOpen() {
		t.Fatal("expected confirm dialog to close on enter")
	}
}

func TestCustomDialogSizeClampsWithinViewport(t *testing.T) {
	m := New()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})

	picker := NewThemePicker()
	_, _ = m.Update(Open(Custom{
		DialogTitle: "Theme",
		Content:     picker,
		Width:       120,
		Height:      60,
	}))

	if !m.IsOpen() {
		t.Fatal("expected dialog to be open")
	}

	w, h := picker.GetSize()
	if w <= 0 || h <= 0 {
		t.Fatalf("expected positive picker bounds, got %dx%d", w, h)
	}
	if w > 36 || h > 8 {
		t.Fatalf("expected picker bounds clamped to viewport, got %dx%d", w, h)
	}
}
