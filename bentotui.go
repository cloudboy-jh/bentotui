package bentotui

import (
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/router"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/footer"
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

func WithFooterBar(v bool) app.Option {
	return app.WithFooterBar(v)
}

func WithFooter(model *footer.Model) app.Option {
	return app.WithFooter(model)
}

func WithFullScreen(v bool) app.Option {
	return app.WithFullScreen(v)
}

func Page(name string, factory func() core.Page) router.Route {
	return router.Page(name, factory)
}
