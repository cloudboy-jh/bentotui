package theme

// ThemeChangedMsg is broadcast when the active theme changes.
// The app model handles this to update any bricks that need the new theme.
type ThemeChangedMsg struct {
	Name  string
	Theme Theme
}
