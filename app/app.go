package app

import (
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/shell"
	"github.com/cloudboy-jh/bentotui/statusbar"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Deprecated: use package shell directly for new code.
type Model = shell.Model

// Deprecated: use package shell directly for new code.
type Option = shell.Option

func New(opts ...Option) *Model { return shell.New(opts...) }

func WithTheme(t theme.Theme) Option { return shell.WithTheme(t) }

func WithPages(routes ...router.Route) Option { return shell.WithPages(routes...) }

func WithStatusBar(v bool) Option { return shell.WithStatusBar(v) }

func WithFullScreen(v bool) Option { return shell.WithFullScreen(v) }

func WithStatus(model *statusbar.Model) Option { return shell.WithStatus(model) }
