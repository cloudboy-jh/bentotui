package bentotui

import (
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/theme"
)

func New(opts ...app.Option) *app.Model {
	return app.New(opts...)
}

func WithTheme(t theme.Theme) app.Option {
	return app.WithTheme(t)
}

func WithPages(routes ...router.Route) app.Option {
	return app.WithPages(routes...)
}

func WithStatusBar(v bool) app.Option {
	return app.WithStatusBar(v)
}

func Page(name string, factory func() core.Page) router.Route {
	return router.Page(name, factory)
}
