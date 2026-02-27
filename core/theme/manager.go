package theme

import (
	"fmt"
	"sync"
)

var (
	mu          sync.RWMutex
	currentName = DefaultName
	current     = Preset(DefaultName)
)

func init() {
	if name, err := loadStoredThemeName(); err == nil {
		if _, ok := lookupPreset(name); ok {
			currentName = name
			current = Preset(name)
		}
	}
}

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

func SetTheme(name string) (Theme, error) {
	return applyTheme(name, true)
}

func PreviewTheme(name string) (Theme, error) {
	return applyTheme(name, false)
}

func applyTheme(name string, persist bool) (Theme, error) {
	t, ok := lookupPreset(name)
	if !ok {
		return Theme{}, fmt.Errorf("unknown theme %q", name)
	}
	mu.Lock()
	currentName = name
	current = t
	mu.Unlock()
	if persist {
		_ = saveThemeName(name)
	}
	return t, nil
}

func lookupPreset(name string) (Theme, bool) {
	for _, n := range AvailableThemes() {
		if n == name {
			return Preset(name), true
		}
	}
	return Theme{}, false
}
