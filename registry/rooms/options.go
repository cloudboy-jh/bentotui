package rooms

import (
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"
)

type DividerMode string

const (
	DividerNone   DividerMode = "none"
	DividerSubtle DividerMode = "subtle"
	DividerNormal DividerMode = "normal"
)

type Option func(*layoutOptions)

type layoutOptions struct {
	gutter  int
	divider DividerMode
}

func WithGutter(n int) Option {
	return func(o *layoutOptions) {
		o.gutter = clampInt(n, 0, 2)
	}
}

func WithDivider(mode string) Option {
	return func(o *layoutOptions) {
		switch strings.ToLower(strings.TrimSpace(mode)) {
		case "subtle":
			o.divider = DividerSubtle
		case "normal":
			o.divider = DividerNormal
		default:
			o.divider = DividerNone
		}
	}
}

func resolveLayoutOptions(opts []Option) layoutOptions {
	o := layoutOptions{divider: DividerNone}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}

func dividerCell(width, height int, mode DividerMode) Sizable {
	if width <= 0 {
		return Static("")
	}
	ch := " "
	if mode == DividerSubtle {
		ch = "."
	}
	if mode == DividerNormal {
		ch = "|"
	}
	line := strings.Repeat(ch, width)
	block := strings.Repeat(line+"\n", max(0, height-1)) + line
	return Static(engine.Constrain(block, width, height))
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
