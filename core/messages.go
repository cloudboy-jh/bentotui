package core

import tea "charm.land/bubbletea/v2"

// NavigateMsg requests a route change in the router.
type NavigateMsg struct {
	Page string
}

// Navigate creates a message that switches to the named page.
func Navigate(page string) tea.Msg {
	return NavigateMsg{Page: page}
}
