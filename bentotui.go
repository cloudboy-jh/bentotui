package bentotui

import (
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/router"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/bar"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
)

func New(opts ...app.Option) *app.Model { return app.New(opts...) }

func WithTheme(t theme.Theme) app.Option { return app.WithTheme(t) }
func WithHeaderBar(v bool) app.Option    { return app.WithHeaderBar(v) }
func WithFooterBar(v bool) app.Option    { return app.WithFooterBar(v) }
func WithFullScreen(v bool) app.Option   { return app.WithFullScreen(v) }
func WithHeader(m *bar.Model) app.Option { return app.WithHeader(m) }
func WithFooter(m *bar.Model) app.Option { return app.WithFooter(m) }

func WithPages(routes ...router.Route) app.Option { return app.WithPages(routes...) }

func Page(name string, factory func() core.Page) router.Route {
	return router.Page(name, factory)
}

// WithCommands registers commands available in the command palette.
func WithCommands(commands ...dialog.Command) app.Option {
	return app.WithCommands(commands...)
}

// Command is an alias for dialog.Command.
type Command = dialog.Command
