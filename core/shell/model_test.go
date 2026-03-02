package shell

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/router"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
)

func TestOpenThemeDialogCmdUsesCompactBounds(t *testing.T) {
	msg := openThemeDialogCmd(280, 90)()
	open, ok := msg.(dialog.OpenMsg)
	if !ok {
		t.Fatalf("expected dialog.OpenMsg, got %T", msg)
	}
	custom, ok := open.Dialog.(dialog.Custom)
	if !ok {
		t.Fatalf("expected dialog.Custom, got %T", open.Dialog)
	}
	if custom.Width > 88 {
		t.Fatalf("expected compact dialog width <= 88, got %d", custom.Width)
	}
	if custom.Width < 52 {
		t.Fatalf("expected compact dialog width >= 52, got %d", custom.Width)
	}
	if custom.Height > 24 {
		t.Fatalf("expected compact dialog height <= 24, got %d", custom.Height)
	}
}

type spyPage struct {
	name      string
	gotTheme  bool
	themeName string
	width     int
	height    int
}

func (p *spyPage) Init() tea.Cmd { return nil }

func (p *spyPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if changed, ok := msg.(theme.ThemeChangedMsg); ok {
		p.gotTheme = true
		p.themeName = changed.Name
	}
	return p, nil
}

func (p *spyPage) View() tea.View { return tea.NewView(p.name) }

func (p *spyPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *spyPage) GetSize() (int, int) { return p.width, p.height }

func (p *spyPage) Title() string { return p.name }

func TestThemeSyncAppliesAfterRouteChange(t *testing.T) {
	_, _ = theme.SetTheme(theme.DefaultName)

	var first *spyPage
	var second *spyPage
	m := New(
		WithPages(
			router.Page("one", func() core.Page {
				first = &spyPage{name: "one"}
				return first
			}),
			router.Page("two", func() core.Page {
				second = &spyPage{name: "two"}
				return second
			}),
		),
	)
	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	tApplied, err := theme.SetTheme(theme.DraculaName)
	if err != nil {
		t.Fatalf("expected dracula theme to be available: %v", err)
	}
	_, _ = m.Update(theme.ThemeChangedMsg{Name: theme.DraculaName, Theme: tApplied})
	if first == nil || !first.gotTheme {
		t.Fatal("expected current page to receive theme changed message")
	}

	_, _ = m.Update(core.Navigate("two"))
	if second == nil {
		t.Fatal("expected second page to be created on navigation")
	}
	if !second.gotTheme {
		t.Fatal("expected navigated page to receive current theme message")
	}
	if second.themeName != theme.DraculaName {
		t.Fatalf("expected navigated page theme %q, got %q", theme.DraculaName, second.themeName)
	}
}
