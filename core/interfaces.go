package core

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Component is the minimum contract for BentoTUI components.
type Component interface {
	tea.Model
}

// Sizeable components can react to terminal size changes.
type Sizeable interface {
	Component
	SetSize(width, height int)
	GetSize() (width, height int)
}

// Focusable components participate in focus management.
type Focusable interface {
	Component
	Focus()
	Blur()
	IsFocused() bool
}

// Positional components can be placed at absolute coordinates.
type Positional interface {
	Component
	SetPosition(x, y int)
}

// Bindable components expose keybindings for help/status rendering.
type Bindable interface {
	Component
	Bindings() []key.Binding
}

// Page is a routable component in the application shell.
type Page interface {
	Component
	Sizeable
	Title() string
}
