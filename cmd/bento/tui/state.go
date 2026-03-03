// Package tui provides the TUI application state and views.
package tui

// View represents the current UI view.
type View int

const (
	ViewMenu View = iota
	ViewInitForm
	ViewComponentList
	ViewDoctor
)

// State holds the application state shared across views.
type State struct {
	CurrentView View
	Width       int
	Height      int

	// Menu state
	MenuSelection int

	// Init form state
	AppName   string
	Module    string
	FormFocus int // 0 = app name, 1 = module, 2 = submit

	// Component list state
	SelectedComponents map[string]bool
	ComponentCursor    int

	// Doctor state
	DoctorRunning bool
	DoctorIndex   int
	DoctorResults []DoctorCheck

	// Log state
	LogLines []string
}

// DoctorCheck represents a single check result for the doctor view.
type DoctorCheck struct {
	Label string
	Pass  bool
	Note  string
	Shown bool
}

// NewState creates a new state with defaults.
func NewState() *State {
	return &State{
		CurrentView:        ViewMenu,
		MenuSelection:      0,
		AppName:            "",
		Module:             "",
		FormFocus:          0,
		SelectedComponents: make(map[string]bool),
		ComponentCursor:    0,
		DoctorRunning:      false,
		DoctorIndex:        0,
		DoctorResults:      make([]DoctorCheck, 0),
		LogLines:           make([]string, 0),
	}
}

// AddLog adds a line to the log.
func (s *State) AddLog(line string) {
	s.LogLines = append(s.LogLines, line)
}

// ClearLog clears the log.
func (s *State) ClearLog() {
	s.LogLines = make([]string, 0)
}
