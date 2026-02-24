package shell

import (
	"image"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/screen"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/surface"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/dialog"
	"github.com/cloudboy-jh/bentotui/ui/components/footer"
)

type Option func(*Model)

type Model struct {
	router     *router.Model
	dialogs    *dialog.Manager
	footer     *footer.Model
	theme      theme.Theme
	showFooter bool
	fullScreen bool
	width      int
	height     int
}

func New(opts ...Option) *Model {
	m := &Model{
		router:     router.New(),
		dialogs:    dialog.New(),
		footer:     footer.New(),
		theme:      theme.CurrentTheme(),
		showFooter: true,
		fullScreen: true,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.footer.SetTheme(m.theme)
	m.dialogs.SetTheme(m.theme)
	return m
}

func WithTheme(t theme.Theme) Option {
	return func(m *Model) {
		m.theme = t
		m.footer.SetTheme(t)
		m.dialogs.SetTheme(t)
	}
}

func WithPages(routes ...router.Route) Option {
	return func(m *Model) {
		m.router = router.New(routes...)
	}
}

func WithFooterBar(v bool) Option {
	return func(m *Model) { m.showFooter = v }
}

func WithFullScreen(v bool) Option {
	return func(m *Model) { m.fullScreen = v }
}

func WithFooter(model *footer.Model) Option {
	return func(m *Model) {
		if model != nil {
			m.footer = model
			m.footer.SetTheme(m.theme)
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.router.Init(), m.dialogs.Init(), m.footer.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	forwardTheme := false
	openThemeDialog := false
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.syncViewport(v.Width, v.Height)
	case theme.OpenThemePickerMsg:
		if !m.dialogs.IsOpen() {
			openThemeDialog = true
		}
	case theme.ThemeChangedMsg:
		m.theme = v.Theme
		m.footer.SetTheme(v.Theme)
		m.dialogs.SetTheme(v.Theme)
		forwardTheme = true
	}

	if openThemeDialog {
		return m, openThemeDialogCmd(m.width, m.height)
	}

	_, dialogCmd := m.dialogs.Update(msg)

	var pageCmd tea.Cmd
	if forwardTheme || !m.dialogs.IsOpen() {
		_, pageCmd = m.router.Update(msg)
	}

	_, footerCmd := m.footer.Update(msg)

	return m, tea.Batch(dialogCmd, pageCmd, footerCmd)
}

func (m *Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = m.fullScreen
	v.BackgroundColor = lipgloss.Color(m.theme.Background)
	if m.width <= 0 || m.height <= 0 {
		return v
	}
	canvas := uv.NewScreenBuffer(m.width, m.height)
	m.draw(canvas, canvas.Bounds())
	v.SetContent(strings.ReplaceAll(canvas.Render(), "\r\n", "\n"))
	return v
}

func (m *Model) Router() *router.Model {
	return m.router
}

func (m *Model) Dialogs() *dialog.Manager {
	return m.dialogs
}

func (m *Model) syncViewport(width, height int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	m.width = width
	m.height = height
	bodyHeight := height
	if m.showFooter {
		bodyHeight--
	}
	if bodyHeight < 0 {
		bodyHeight = 0
	}
	m.router.SetSize(width, bodyHeight)
	m.dialogs.SetSize(width, height)
	m.footer.SetSize(width, 1)
}

func (m *Model) draw(scr uv.Screen, area image.Rectangle) {
	screen.Clear(scr)

	w := area.Dx()
	h := area.Dy()
	if w <= 0 || h <= 0 {
		return
	}

	bodyHeight := h
	if m.showFooter {
		bodyHeight--
	}
	if bodyHeight < 0 {
		bodyHeight = 0
	}

	shellBG := surface.Fill(w, h, m.theme.Background)
	bodyView := m.router.View()

	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(shellBG).ID("shell-bg").X(area.Min.X).Y(area.Min.Y).Z(0),
		lipgloss.NewLayer(core.ViewLayer(bodyView)).ID("body").X(area.Min.X).Y(area.Min.Y).Z(1),
	}

	if m.showFooter {
		layers = append(layers,
			lipgloss.NewLayer(core.ViewLayer(m.footer.View())).
				ID("footer").
				X(area.Min.X).
				Y(area.Min.Y+bodyHeight).
				Z(2),
		)
	}

	if m.dialogs.IsOpen() {
		scrim := surface.Fill(w, h, m.theme.Scrim)
		dlgView := m.dialogs.View()
		dlg := core.ViewString(dlgView)
		dlgW, dlgH := lipgloss.Size(dlg)
		dlgX := area.Min.X + max(0, (w-dlgW)/2)
		dlgY := area.Min.Y + max(0, (h-dlgH)/2)
		layers = append(layers,
			lipgloss.NewLayer(scrim).ID("scrim").X(area.Min.X).Y(area.Min.Y).Z(3),
			lipgloss.NewLayer(core.ViewLayer(dlgView)).ID("dialog").X(dlgX).Y(dlgY).Z(4),
		)
	}

	lipgloss.NewCanvas(layers...).Draw(scr, area)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func openThemeDialogCmd(width, height int) tea.Cmd {
	picker := dialog.NewThemePicker()
	picker.SetSize(width, height)
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Theme",
			Content:     picker,
			Width:       48,
			Height:      14,
		})
	}
}
