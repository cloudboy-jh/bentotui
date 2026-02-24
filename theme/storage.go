package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type storedConfig struct {
	Theme string `json:"theme"`
}

func configPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "bentotui", "theme.json"), nil
}

func loadStoredThemeName() (string, error) {
	path, err := configPath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cfg storedConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return "", err
	}
	return cfg.Theme, nil
}

func saveThemeName(name string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(storedConfig{Theme: name}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
