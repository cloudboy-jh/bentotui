package theme

type ThemeChangedMsg struct {
	Name  string
	Theme Theme
}

type OpenThemePickerMsg struct{}

func OpenThemePicker() OpenThemePickerMsg { return OpenThemePickerMsg{} }
