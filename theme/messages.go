package theme

// ThemeChangedMsg is broadcast when the active theme changes (via SetTheme
// or PreviewTheme). Components that cache derived state (e.g. textinput
// styles) should listen for this and call syncStyles().
type ThemeChangedMsg struct {
	Name  string
	Theme Theme
}

// OpenThemePickerMsg signals the app to open the theme picker dialog.
type OpenThemePickerMsg struct{}

func OpenThemePicker() OpenThemePickerMsg { return OpenThemePickerMsg{} }
