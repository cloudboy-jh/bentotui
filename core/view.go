package core

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func ViewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	return fmt.Sprint(v.Content)
}
