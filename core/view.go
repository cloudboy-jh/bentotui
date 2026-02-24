package core

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func ViewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}

func ViewLayer(v tea.View) tea.Layer {
	if v.Content == nil {
		empty := tea.NewView("")
		return empty.Content
	}
	return v.Content
}
