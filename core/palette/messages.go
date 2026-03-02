package palette

// OpenCommandPaletteMsg is broadcast to open the command palette overlay.
// Handled by the shell; ignored if a dialog is already open.
type OpenCommandPaletteMsg struct{}

// OpenCommandPalette returns an OpenCommandPaletteMsg for use in tea.Cmd returns.
func OpenCommandPalette() OpenCommandPaletteMsg { return OpenCommandPaletteMsg{} }
