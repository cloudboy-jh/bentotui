package shell

import (
	"image"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/screen"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/palette"
	"github.com/cloudboy-jh/bentotui/core/router"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/bar"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
)

type Option func(*Model)

type Model struct {
	router     *router.Model
	dialogs    *dialog.Manager
	header     *bar.Model
	footer     *bar.Model
	theme      theme.Theme
	showHeader bool
	showFooter bool
	fullScreen bool
	width      int
	height     int
	commands   []dialog.Command
}

func New(opts ...Option) *Model {
	m := &Model{
		router:     router.New(),
		dialogs:    dialog.New(),
		header:     bar.New(),
		footer:     bar.New(),
		theme:      theme.CurrentTheme(),
		showHeader: true,
		showFooter: true,
		fullScreen: true,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.header.SetTheme(m.theme)
	m.footer.SetTheme(m.theme)
	m.dialogs.SetTheme(m.theme)
	return m
}

func WithTheme(t theme.Theme) Option {
	return func(m *Model) {
		m.theme = t
		m.header.SetTheme(t)
		m.footer.SetTheme(t)
		m.dialogs.SetTheme(t)
	}
}

func WithHeaderBar(v bool) Option {
	return func(m *Model) { m.showHeader = v }
}

func WithPages(routes ...router.Route) Option {
	return func(m *Model) {
		m.router = router.New(routes...)
	}
}

func WithFooterBar(v bool) Option {
	return func(m *Model) { m.showFooter = v }
}

func WithHeader(model *bar.Model) Option {
	return func(m *Model) {
		if model != nil {
			m.header = model
			m.header.SetTheme(m.theme)
		}
	}
}

func WithFullScreen(v bool) Option {
	return func(m *Model) { m.fullScreen = v }
}

func WithFooter(model *bar.Model) Option {
	return func(m *Model) {
		if model != nil {
			m.footer = model
			m.footer.SetTheme(m.theme)
		}
	}
}

// WithCommands registers commands that appear in the command palette when
// opened via palette.OpenCommandPaletteMsg. App-level commands are merged
// with any built-in shell commands.
func WithCommands(commands ...dialog.Command) Option {
	return func(m *Model) {
		m.commands = append(m.commands, commands...)
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.router.Init(), m.dialogs.Init(), m.header.Init(), m.footer.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	forwardTheme := false
	syncThemeAfterRoute := false
	openThemeDialog := false
	openPaletteDialog := false
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.syncViewport(v.Width, v.Height)
	case core.NavigateMsg:
		syncThemeAfterRoute = true
	case theme.OpenThemePickerMsg:
		if !m.dialogs.IsOpen() {
			openThemeDialog = true
		}
	case palette.OpenCommandPaletteMsg:
		if !m.dialogs.IsOpen() {
			openPaletteDialog = true
		}
	case theme.ThemeChangedMsg:
		m.theme = v.Theme
		m.header.SetTheme(v.Theme)
		m.footer.SetTheme(v.Theme)
		m.dialogs.SetTheme(v.Theme)
		forwardTheme = true
	}

	if openThemeDialog {
		return m, openThemeDialogCmd(m.width, m.height)
	}
	if openPaletteDialog {
		return m, openCommandPaletteCmd(m.width, m.height, m.commands)
	}

	_, dialogCmd := m.dialogs.Update(msg)

	var pageCmd tea.Cmd
	if forwardTheme || !m.dialogs.IsOpen() {
		_, pageCmd = m.router.Update(msg)
		if syncThemeAfterRoute {
			_, themeCmd := m.router.Update(theme.ThemeChangedMsg{Name: theme.CurrentThemeName(), Theme: m.theme})
			pageCmd = tea.Batch(pageCmd, themeCmd)
		}
	}

	_, headerCmd := m.header.Update(msg)
	_, footerCmd := m.footer.Update(msg)

	return m, tea.Batch(dialogCmd, pageCmd, headerCmd, footerCmd)
}

func (m *Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = m.fullScreen
	v.BackgroundColor = lipgloss.Color(m.theme.Surface.Canvas)
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
	if m.showHeader {
		bodyHeight--
	}
	if m.showFooter {
		bodyHeight--
	}
	if bodyHeight < 0 {
		bodyHeight = 0
	}
	m.router.SetSize(width, bodyHeight)
	m.dialogs.SetSize(width, height)
	m.header.SetSize(width, 1)
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
	if m.showHeader {
		bodyHeight--
	}
	if m.showFooter {
		bodyHeight--
	}
	if bodyHeight < 0 {
		bodyHeight = 0
	}

	shellBG := primitives.Fill(w, h, m.theme.Surface.Canvas)
	bodyView := m.router.View()

	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(shellBG).ID("shell-bg").X(area.Min.X).Y(area.Min.Y).Z(0),
		lipgloss.NewLayer(core.ViewLayer(bodyView)).ID("body").X(area.Min.X).Y(area.Min.Y + boolToInt(m.showHeader)).Z(1),
	}

	if m.showHeader {
		layers = append(layers,
			lipgloss.NewLayer(core.ViewLayer(m.header.View())).
				ID("header").
				X(area.Min.X).
				Y(area.Min.Y).
				Z(2),
		)
	}

	if m.showFooter {
		layers = append(layers,
			lipgloss.NewLayer(core.ViewLayer(m.footer.View())).
				ID("footer").
				X(area.Min.X).
				Y(area.Min.Y+bodyHeight+boolToInt(m.showHeader)).
				Z(3),
		)
	}

	if m.dialogs.IsOpen() {
		scrim := primitives.Fill(w, h, m.theme.Dialog.Scrim)
		dlgView := m.dialogs.View()
		dlg := core.ViewString(dlgView)
		dlgW, dlgH := lipgloss.Size(dlg)
		dlgX := area.Min.X + max(0, (w-dlgW)/2)
		dlgY := area.Min.Y + max(0, (h-dlgH)/2)
		layers = append(layers,
			lipgloss.NewLayer(scrim).ID("scrim").X(area.Min.X).Y(area.Min.Y).Z(4),
			lipgloss.NewLayer(core.ViewLayer(dlgView)).ID("dialog").X(dlgX).Y(dlgY).Z(5),
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func openCommandPaletteCmd(width, height int, commands []dialog.Command) tea.Cmd {
	p := dialog.NewCommandPalette(commands)
	maxWidth := max(44, width-8)
	maxHeight := max(12, height-6)
	modalWidth := clampInt(72, min(52, maxWidth), min(88, maxWidth))
	modalHeight := clampInt(height-10, min(14, maxHeight), min(28, maxHeight))
	p.SetSize(maxInt(1, modalWidth-4), maxInt(1, modalHeight-4))
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Commands",
			Content:     p,
			Width:       modalWidth,
			Height:      modalHeight,
		})
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func openThemeDialogCmd(width, height int) tea.Cmd {
	picker := dialog.NewThemePicker()
	maxWidth := max(44, width-8)
	maxHeight := max(12, height-6)
	modalWidth := clampInt(72, min(52, maxWidth), min(88, maxWidth))
	modalHeight := clampInt(height-10, min(14, maxHeight), min(24, maxHeight))
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Theme",
			Content:     picker,
			Width:       modalWidth,
			Height:      modalHeight,
		})
	}
}

func clampInt(v, minV, maxV int) int {
	if maxV < minV {
		return minV
	}
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
