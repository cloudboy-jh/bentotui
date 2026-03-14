package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

type Sizable = engine.Sizable

func Static(s string) Sizable {
	return engine.Static(s)
}

func RenderFunc(fn func(width, height int) string) Sizable {
	return engine.RenderFunc(fn)
}

func min(a, b int) int {
	return engine.Min(a, b)
}

func max(a, b int) int {
	return engine.Max(a, b)
}
