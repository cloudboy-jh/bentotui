package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/layout"
	"github.com/cloudboy-jh/bentotui/ui/components/panel"
)

func main() {
	app := bentotui.New(
		bentotui.WithPages(
			bentotui.Page("home", func() core.Page { return newHomePage() }),
		),
		bentotui.WithFooterBar(true),
	)

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui demo failed: %v\n", err)
	}
}

type homePage struct {
	root   *layout.Split
	width  int
	height int
}

func newHomePage() *homePage {
	sidebar := panel.New(panel.Title("Sidebar"), panel.Content(staticText("Sessions\nFiles\nSearch")))
	main := panel.New(panel.Title("Main"), panel.Content(staticText("Welcome to BentoTUI v0.1.")))
	root := layout.Horizontal(
		layout.Fixed(30, sidebar),
		layout.Flex(1, main),
	)
	return &homePage{root: root}
}

func (p *homePage) Init() tea.Cmd { return nil }

func (p *homePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := p.root.Update(msg)
	return p, cmd
}

func (p *homePage) View() tea.View { return p.root.View() }

func (p *homePage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.root.SetSize(width, height)
}

func (p *homePage) GetSize() (width, height int) { return p.width, p.height }

func (p *homePage) Title() string { return "Home" }

type staticText string

func (s staticText) Init() tea.Cmd                           { return nil }
func (s staticText) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s staticText) View() tea.View                          { return tea.NewView(string(s)) }
