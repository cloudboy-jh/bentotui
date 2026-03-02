package theme

import (
	"fmt"
	"sort"
	"sync"
)

var (
	mu          sync.RWMutex
	registry    = make(map[string]Theme)
	currentName = DefaultName
	current     Theme
)

func init() {
	for name, t := range builtinThemes {
		if err := registerThemeLocked(name, t); err != nil {
			panic(err)
		}
	}
	if name, err := loadStoredThemeName(); err == nil {
		if t, ok := registry[name]; ok {
			currentName = name
			current = t
		}
	}
	if current.Name == "" {
		if t, ok := registry[DefaultName]; ok {
			currentName = DefaultName
			current = t
		}
	}
}

// CurrentTheme returns the active theme. Safe for concurrent use — callers
// should call this in View() rather than storing a cached copy.
func CurrentTheme() Theme {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

func CurrentThemeName() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentName
}

// SetTheme sets the active theme by name and persists the choice to disk.
// Safe to call from main() before tea.NewProgram().Run().
func SetTheme(name string) (Theme, error) {
	return applyTheme(name, true)
}

// PreviewTheme sets the active theme without persisting. Used by ThemePicker
// during live preview — ESC reverts by calling PreviewTheme(baseTheme).
func PreviewTheme(name string) (Theme, error) {
	return applyTheme(name, false)
}

// RegisterTheme adds a custom theme to the registry. Call before Run().
func RegisterTheme(name string, t Theme) error {
	mu.Lock()
	defer mu.Unlock()
	return registerThemeLocked(name, t)
}

func applyTheme(name string, persist bool) (Theme, error) {
	mu.Lock()
	defer mu.Unlock()
	t, ok := registry[name]
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q", name)
	}
	currentName = name
	current = t
	if persist {
		_ = saveThemeName(name)
	}
	return t, nil
}

func registerThemeLocked(name string, t Theme) error {
	if name == "" {
		return fmt.Errorf("theme name is required")
	}
	t.Name = name
	if err := validateTheme(t); err != nil {
		return fmt.Errorf("invalid theme %q: %w", name, err)
	}
	registry[name] = t
	if current.Name == "" {
		currentName = name
		current = t
	}
	return nil
}

func availableThemeNames() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		if name != DefaultName {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	// Default is always first; the rest are sorted.
	if _, ok := registry[DefaultName]; ok {
		names = append([]string{DefaultName}, names...)
	}
	return names
}
