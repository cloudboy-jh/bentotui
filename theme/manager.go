package theme

import (
	"fmt"
	"sort"
	"sync"
)

// The global manager is an app-level convenience.
// Bricks do NOT call CurrentTheme() internally — they use whatever Theme
// was passed via WithTheme() or SetTheme(). The global store is only for
// apps that want a single active theme shared across their entire UI.
var (
	mu          sync.RWMutex
	registry    = map[string]Theme{}
	currentName = DefaultName
	current     Theme
)

func init() {
	// Register all built-in presets.
	for name, t := range presets {
		registry[name] = t
	}
	current = presets[DefaultName]
}

// CurrentTheme returns the app-level active theme.
// Safe for concurrent use. Bricks call this only as a fallback when no
// explicit theme has been provided via WithTheme() / SetTheme().
func CurrentTheme() Theme {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

// CurrentThemeName returns the name of the active theme.
func CurrentThemeName() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentName
}

// SetTheme sets and persists the active theme by name.
func SetTheme(name string) (Theme, error) {
	mu.Lock()
	defer mu.Unlock()
	t, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("theme %q not found", name)
	}
	currentName = name
	current = t
	return t, nil
}

// PreviewTheme sets the active theme without persisting (for live preview).
func PreviewTheme(name string) (Theme, error) {
	return SetTheme(name)
}

// RegisterTheme adds a custom theme to the registry.
func RegisterTheme(name string, t Theme) error {
	if name == "" {
		return fmt.Errorf("theme name is required")
	}
	mu.Lock()
	defer mu.Unlock()
	registry[name] = t
	return nil
}

// AvailableThemes returns all registered theme names, default first.
func AvailableThemes() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		if name != DefaultName {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return append([]string{DefaultName}, names...)
}
