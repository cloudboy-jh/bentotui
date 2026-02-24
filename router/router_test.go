package router

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type testPage struct {
	title  string
	width  int
	height int
}

func (p *testPage) Init() tea.Cmd                           { return nil }
func (p *testPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return p, nil }
func (p *testPage) View() tea.View                          { return tea.NewView(p.title) }
func (p *testPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}
func (p *testPage) GetSize() (width, height int) { return p.width, p.height }
func (p *testPage) Title() string                { return p.title }

func TestLazyFactoryOnlyBuildsVisitedPages(t *testing.T) {
	hits := map[string]int{}

	m := New(
		Page("home", func() core.Page {
			hits["home"]++
			return &testPage{title: "home"}
		}),
		Page("settings", func() core.Page {
			hits["settings"]++
			return &testPage{title: "settings"}
		}),
	)

	if hits["home"] != 1 {
		t.Fatalf("expected initial page to be created once, got %d", hits["home"])
	}
	if hits["settings"] != 0 {
		t.Fatalf("expected settings page to not be created yet, got %d", hits["settings"])
	}

	_, _ = m.Update(core.Navigate("settings"))

	if hits["settings"] != 1 {
		t.Fatalf("expected settings page to be created on first visit, got %d", hits["settings"])
	}

	_, _ = m.Update(core.Navigate("settings"))
	if hits["settings"] != 1 {
		t.Fatalf("expected settings page to be cached after first visit, got %d", hits["settings"])
	}
}
