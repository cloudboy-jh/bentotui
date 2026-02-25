package app

import (
	"github.com/cloudboy-jh/bentotui/core/router"
	"github.com/cloudboy-jh/bentotui/core/shell"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/footer"
)

// Deprecated: use package core/shell directly for new code.
type Model = shell.Model

// Deprecated: use package core/shell directly for new code.
type Option = shell.Option

func New(opts ...Option) *Model { return shell.New(opts...) }

func WithTheme(t theme.Theme) Option { return shell.WithTheme(t) }

func WithPages(routes ...router.Route) Option { return shell.WithPages(routes...) }

func WithFooterBar(v bool) Option { return shell.WithFooterBar(v) }

func WithFullScreen(v bool) Option { return shell.WithFullScreen(v) }

func WithFooter(model *footer.Model) Option { return shell.WithFooter(model) }
